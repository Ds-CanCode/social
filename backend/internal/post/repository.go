package funcs

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	database "funcs/internal/database"
	pkg "funcs/pkg"
)

func InsertPost(post POST) error {
	db := database.GetDb()
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Insert the post first
	query := `INSERT INTO posts (user_id, group_id ,title ,content ,pathimg ,types, createdAt)
	VALUES(?,?,?,?,?,?,?)`
	result, err := tx.Exec(query, post.Creator.Id, post.GroupId, post.Title, post.Content, post.Path, post.Type, time.Now())
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return err
	}
	postId, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	// If the post is semi-private, insert into allowedUsersPost
	if post.Type == "semi-private" {
		for _, user := range post.AllowedUser {
			query := `INSERT INTO allowedUsersPost (user_id, post_id) VALUES(?,?)`
			_, err := tx.Exec(query, user.Id, postId)
			if err != nil {
				fmt.Println("post not inserre", err)
				tx.Rollback()
				return err
			}
		}
	}

	// Commit the transaction
	return tx.Commit()
}

func GroupPost(id, groupID int) []POST {
	query := `SELECT p.* FROM posts p
				JOIN group_members gm ON p.user_id = gm.user_id
				WHERE p.types = "group"
				ORDER BY p.createdAt DESC;
			`

	db := database.GetDb()
	rows, err := db.Query(query, id, groupID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var posts []POST
	for rows.Next() {
		post, err := ScanPostRow(rows, id)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		post.Creator, err = GetUserInfo(post.Creator.Id)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		reactionCount, err := pkg.CountRect(post.Id)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		post.Likes = reactionCount.LikeCount
		post.DisLikes = reactionCount.DislikeCount
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating rows:", err)
	}
	return posts
}

func getFeedPosts(id,offset  int) ([]POST, error) {
	var posts []POST
	db := database.GetDb()
	rows, err := db.Query(feedQuery, id, id, id ,  offset )
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return posts, err
	}

	for rows.Next() {
		post, err := ScanPostRow(rows, id)
		if err != nil {
			return posts, err
		}

		post.Likes, err = GetReactions(post.Id, "LIKE")
		if err != nil {
			return posts, err
		}

		post.DisLikes, err = GetReactions(post.Id, "DISLIKE")
		if err != nil {
			return posts, err
		}

		post.Creator, err = GetUserInfo(post.Creator.Id)
		if err != nil {
			return posts, err
		}

		posts = append(posts, post)
	}
	return posts, nil
}

func GetReactions(postId int, reaction_type string) (int, error) {
	db := database.GetDb()
	query := "SELECT COUNT(*) FROM reactionsPosts WHERE post_id=? AND reaction_type=?"
	var nbr int

	err := db.QueryRow(query, postId, reaction_type).Scan(&nbr)
	if err == sql.ErrNoRows {
		err = nil
	}
	return nbr, err
}

func GetUserPostDB(userId, targetUserId int) ([]POST, error) {
	db := database.GetDb()
	var posts []POST
	query := `SELECT id,user_id,group_id,title,content,pathimg,types,createdAt FROM posts p
				WHERE (
						(types!='group' AND types!='semi-private')
						OR 
						(p.types = 'semi-private' AND EXISTS (
							SELECT 1 FROM allowedUsersPost aup
							WHERE aup.user_id = ? 
							AND aup.post_id = p.id)
						) 
					  ) 
					  AND user_id=? 
					  ORDER BY createdAt DESC;`

	rows, err := db.Query(query, userId, targetUserId)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return posts, err
	}

	for rows.Next() {
		post, err := ScanPostRow(rows, userId)
		if err != nil {
			return posts, err
		}

		if post.Type == "semi-private" {
			post.AllowedUser, err = GetAllowedUsers(post.Id)
			if err != nil {
				return posts, err
			}
		}

		reactionCount, err := pkg.CountRect(post.Id)
		if err != nil {
			return posts, err
		}
		post.Likes = reactionCount.LikeCount
		post.DisLikes = reactionCount.DislikeCount

		posts = append(posts, post)
	}

	return posts, err
}

func GetAllowedUsers(postId int) ([]User, error) {
	query := "SELECT user_id FROM allowedUsersPost WHERE post_id=?"
	var users []User
	db := database.GetDb()

	rows, err := db.Query(query, postId)
	if err != nil {
		return users, err
	}

	for rows.Next() {
		var id int
		var user User
		err = rows.Scan(&id)
		if err != nil {
			return users, err
		}

		user, err = GetUserInfo(id)
		if err != nil {
			return users, err
		}

		users = append(users, user)
	}

	return users, nil
}

func GetUserInfo(id int) (User, error) {
	var user User
	user.Id = id
	query := "SELECT firstName,lastName,nickname,avatar FROM users WHERE id=?"
	db := database.GetDb()

	err := db.QueryRow(query, id).Scan(&user.FirstName, &user.LastName, &user.Nickname, &user.Avatar)
	return user, err
}

func ScanPostRow(rows *sql.Rows, userId int) (POST, error) {
	var post POST
	err := rows.Scan(&post.Id, &post.Creator.Id, &post.GroupId, &post.Title, &post.Content, &post.Path, &post.Type, &post.CreateAt)
	if err != nil {
		return post, err
	}
	post.IsLiked, err = IsLiked(post.Id, userId)
	return post, err
}

func IsLiked(postId, userId int) (int, error) {
	query := "SELECT reaction_type FROM reactionsPosts WHERE post_id=? AND user_id=?"
	db := database.GetDb()

	var reaction string
	err := db.QueryRow(query, postId, userId).Scan(&reaction)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	if reaction == "LIKE" {
		return 1, nil
	}
	return 0, nil
}

// const query string = `SELECT DISTINCT p.*, u.firstName, u.lastName, u.avatar, u.nickname
// 						FROM (
// 							SELECT DISTINCT p.* FROM posts p
// 							WHERE p.types = 'public'

// 							UNION

// 							SELECT DISTINCT p.* FROM posts p
// 							WHERE p.types = 'private'
// 							AND EXISTS (
// 								SELECT 1 FROM folowers f
// 								WHERE (f.user1 = ? AND f.user2 = p.user_id
// 									OR f.user2 = ? AND f.user1 = p.user_id)
// 								AND f.accepted = 1
// 							)

// 							UNION

// 							SELECT DISTINCT p.* FROM posts p
// 							WHERE p.types = 'semi-public'
// 							AND EXISTS (
// 								SELECT 1 FROM allowedUsersPost aup
// 								WHERE aup.user_id = ?
// 								AND aup.post_id = p.id
// 							)

// 							UNION

// 							SELECT DISTINCT p.* FROM posts p
// 							WHERE p.types = 'group'
// 							AND EXISTS (
// 								SELECT 1 FROM groups_members gm
// 								WHERE gm.user_id = ?
// 								AND gm.group_id = p.group_id
// 								AND gm.accepted = 1
// 							)
// 						) AS p
// 						JOIN users u ON u.id = p.user_id
// 						ORDER BY p.createdAt DESC;
// `

const feedQuery string = `
SELECT * FROM (
    SELECT p.*
    FROM posts p
    INNER JOIN users u ON p.user_id = u.id
    WHERE 
        ((p.types = 'public' AND u.profileType = 0 ) OR p.user_id = ?)
        OR 
        (p.types = 'private' AND EXISTS (
            SELECT 1 FROM folowers f 
            WHERE f.user1 = ? AND f.user2 = p.user_id
            AND f.accepted = 1
        ))
        OR 
        (p.types = 'semi-public' AND EXISTS (
            SELECT 1 FROM allowedUsersPost aup
            WHERE aup.user_id = ?
            AND aup.post_id = p.id
        ))
    ORDER BY p.createdAt DESC
    LIMIT 100
) subquery
LIMIT 1 OFFSET ?;
`

// OR
//     (p.types = 'group' AND EXISTS (
//         SELECT 1 FROM groups_members gm
//         WHERE gm.user_id = ?
//         AND gm.group_id = p.group_id
//         AND gm.accepted = 1
//     ))
