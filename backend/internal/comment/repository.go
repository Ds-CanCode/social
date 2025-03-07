package funcs

import (
	"fmt"
	"time"

	dataBase "funcs/internal/database"
	pkg "funcs/pkg"
)

func AddComment(InsertedComment Comment) (Comment, error) {
	db := dataBase.GetDb()
	query := `INSERT INTO comments (user_id, post_id , content, pathimg, createdAt) VALUES (?, ?, ?, ?, ?)`
	result, err := db.Exec(query, InsertedComment.UserID, InsertedComment.PostID, InsertedComment.Content, InsertedComment.PathImg, time.Now())
	if err != nil {
		return InsertedComment, err
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		return InsertedComment, err
	}

	var comment Comment
	query = `SELECT 
    c.id, 
    c.user_id, 
    u.firstName, 
    u.lastName, 
    u.avatar, 
    c.post_id, 
    c.content, 
    c.pathimg, 
    c.createdAt 
	FROM comments c 
	JOIN users u ON c.user_id = u.id
	WHERE c.id = ?;
	`
	err = db.QueryRow(query, lastID).Scan(&comment.ID, &comment.User.Id, &comment.User.FirstName, &comment.User.LastName, &comment.User.Avatar, &comment.PostID, &comment.Content, &comment.PathImg, &comment.CreatedAt)
	if err != nil {
		return comment, err
	}

	return comment, nil
}

func GetComments(postId int) ([]Comment, error) {
	db := dataBase.GetDb()
	query := `
	SELECT 
    c.id, 
    c.user_id, 
    u.firstName, 
    u.lastName, 
    u.avatar, 
    c.post_id, 
    c.content, 
    c.pathimg, 
    c.createdAt 
	FROM comments c
	JOIN users u ON c.user_id = u.id
	WHERE c.post_id = ?
	ORDER BY c.createdAt DESC;
	`
	rows, err := db.Query(query, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.UserID, &comment.User.FirstName, &comment.User.LastName, &comment.User.Avatar, &comment.PostID, &comment.Content, &comment.PathImg, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		reactionCount, err := pkg.CountCommentRect(comment.ID)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		comment.Likes = reactionCount.LikeCount
		comment.DisLikes = reactionCount.DislikeCount
		comments = append(comments, comment)
	}
	return comments, nil
}
