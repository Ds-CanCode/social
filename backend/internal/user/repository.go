package funcs

import (
	"database/sql"

	database "funcs/internal/database"
)

// GetRows scans multiple rows and appends them to the followers slice
func GetRows(rows *sql.Rows, followers *[]Followers) error {
	for rows.Next() {
		var follower Followers
		err := rows.Scan(&follower.Id, &follower.FirstName, &follower.LastName, &follower.Avatar)
		if err != nil {
			return err
		}
		*followers = append(*followers, follower)
	}
	return rows.Err()
}

func GetProfileType(id int) (bool, error) {
	query := "SELECT profileType FROM users WHERE id=?"
	db := database.GetDb()
	var profileType bool
	err := db.QueryRow(query, id).Scan(&profileType)
	return profileType, err
}

func IsFriends(id1, id2 int) (int, error) {
	query := `SELECT EXISTS (
				SELECT 1 FROM folowers 
				WHERE (user1 = ? AND user2 = ?)
				AND accepted = 1
			)`
	db := database.GetDb()
	var exists int
	err := db.QueryRow(query, id1, id2, id1, id2).Scan(&exists)
	if err != nil {
		return 0, err
	}
	return exists, nil
}

func CheckFollowstatus(id1, id2 int) (bool, error) {
	query := `SELECT EXISTS (
				SELECT 1 FROM folowers
				WHERE ((user1 = ? AND user2 = ?) 
				AND accepted = 1)
			)`
	db := database.GetDb()
	var exists bool
	err := db.QueryRow(query, id1, id2).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// -1 not friend
// 0 waiting
// 1 friend
func IsFriends2(id1, id2 int) (int, error) {
	query := `SELECT accepted FROM folowers WHERE (user1 = ? AND user2 = ?) OR (user2 = ? AND user1 = ?)`

	db := database.GetDb()
	var exists int
	err := db.QueryRow(query, id1, id2, id1, id2).Scan(&exists)
	if err == sql.ErrNoRows {
		return -1, nil
	}
	if err != nil {
		return 0, err
	}
	return exists, nil
}

func GetUserDB(userId int) (Profile, error) {
	var profile Profile
	query := `SELECT firstName,lastName,datebirth,avatar,nickname,aboutme,profileType,createdAt FROM users
	 WHERE id=?`
	db := database.GetDb()

	err := db.QueryRow(query, userId).Scan(&profile.FirstName, &profile.LastName, &profile.Datebirth, &profile.Avatar, &profile.NickName, &profile.Aboutme, &profile.Type, &profile.CreatedAt)
	if err != nil {
		return profile, err
	}
	profile.Id = userId

	profile.Followers, err = CountFollowers(userId)
	if err != nil {
		return profile, err
	}

	profile.NbrPosts, err = CountPosts(userId)
	if err != nil {
		return profile, err
	}
	return profile, nil
}

func CountPosts(id int) (int, error) {
	query := `SELECT COUNT(*) FROM posts WHERE user_id = ?`
	db := database.GetDb()
	var count int

	err := db.QueryRow(query, id).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func CountFollowers(id int) (int, error) {
	query := `SELECT COUNT(*) FROM folowers 
			  WHERE (user1 = ? OR user2 = ?) 
			  AND accepted = 1`

	db := database.GetDb()
	var count int

	err := db.QueryRow(query, id, id).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func DeleteFollower(id, targetId int) error {
	query := "DELETE  FROM folowers WHERE (user1 = ? AND user2 = ?) ;"
	db := database.GetDb()
	_, err := db.Exec(query, id, targetId)
	return err
}

func AddFollower(id, targetId, state int) error {
	query := "INSERT INTO folowers (user1, user2,accepted) VALUES (?,?,?);"
	db := database.GetDb()
	_, err := db.Exec(query, id, targetId, state)
	return err
}

func AcceptRequest(id, targetId int) error {
	query := "UPDATE folowers SET accepted = ? WHERE user1 = ? AND user2 = ?"
	db := database.GetDb()
	_, err := db.Exec(query, 1, targetId, id)
	return err
}

func RejectRequest(id, targetId int) error {
	return DeleteFollower(targetId, id)
}

func ChangeProfileTypeDB(userId, newType int) error {
	query := "UPDATE users SET profileType=? WHERE id=?"
	db := database.GetDb()
	_, err := db.Exec(query, newType, userId)
	return err
}

func GetProfileInfo(TargetId int, userId int) (Profile, error) {
	var profile Profile
	query := `SELECT firstName,lastName,datebirth,nickname,avatar,aboutme,profileType,createdAt FROM users WHERE id=?`
	queryFollow := `SELECT accepted FROM folowers WHERE user1=? AND user2=? `
	db := database.GetDb()

	err := db.QueryRow(query, TargetId).Scan(&profile.FirstName, &profile.LastName, &profile.Datebirth, &profile.NickName, &profile.Avatar, &profile.Aboutme, &profile.Type, &profile.CreatedAt)
	if err != nil {
		return profile, err
	}
	err = db.QueryRow(queryFollow, userId, TargetId).Scan(&profile.IsFollow)
	if err == sql.ErrNoRows {
		profile.IsFollow = nil
	}
	return profile, nil
}

func GetNowRequest(TargetId int, userId int) (bool, error) {
	queryFollow := `SELECT accepted FROM folowers WHERE user1=? AND user2=? AND accepted = 0`
	db := database.GetDb()
	var is bool
	err := db.QueryRow(queryFollow, TargetId, userId).Scan(&is)
	if err == sql.ErrNoRows {
		return false, err
	}
	return true, nil
}

func getFollowersDB(id int) ([]Followers, error) {
	query1 := `
		SELECT users.id, users.firstName, users.lastName, users.avatar
		FROM users
		JOIN folowers ON users.id = folowers.user1
		WHERE folowers.user2 = ? AND folowers.accepted = 1;
	`
	db := database.GetDb()
	var followers []Followers

	// First query
	rows1, err := db.Query(query1, id)
	if err != nil {
		return nil, err
	}
	defer rows1.Close()

	err = GetRows(rows1, &followers)
	if err != nil {
		return nil, err
	}
	return followers, nil
}

// query2 := `
// 	SELECT users.id, users.firstName, users.lastName
// 	FROM users
// 	JOIN folowers ON users.id = folowers.user2
// 	WHERE folowers.user1 = ? AND folowers.accepted = 1;
// `
// // Second query
// rows2, err := db.Query(query2, id)
// if err != nil {
// 	return nil, err
// }
// defer rows2.Close()

// err = GetRows(rows2, &followers)
// if err != nil {
// 	return nil, err
// }
