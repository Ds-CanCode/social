package funcs

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	dataBase "funcs/internal/database"

	pkg "funcs/pkg"
)

func PrivateMessages(Sender, receiver int, Message string) error {
	db := dataBase.GetDb()
	/*CHECK RECEIVER VALIDITY*/
	var isValide bool
	err := db.QueryRow(`SELECT 1 FROM users WHERE id = ?`, receiver).Scan(&isValide)
	if err != nil {
		return err
	}

	if isValide {
		_, err = db.Exec(`INSERT INTO users_messages ( sender_id, receiver_id, content, createdAt) VALUES (?, ?, ?, ?)`, Sender, receiver, Message, time.Now())
		if err != nil {
			return err
		}
	}
	return nil
}

func GroupMessages(Sender, receiver int, Message string) error {
	db := dataBase.GetDb()
	/*CHECK Group VALIDITY*/
	var isValide bool
	err := db.QueryRow(`SELECT 1 FROM groups WHERE id = ?`, receiver).Scan(&isValide)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if isValide {
		_, err := db.Exec(`INSERT INTO groups_messages (sender_id, group_id, content, createdAt) VALUES (?, ?, ?, ?)`, Sender, receiver, Message, time.Now())
		if err != nil {
			return err
		}
	}

	return nil
}

func GetGroupusers(id int) ([]int, error) {
	db := dataBase.GetDb()

	rows, err := db.Query(`SELECT user_id FROM group_members WHERE group_id = ?`, id)
	if err != nil {
		return []int{}, err
	}
	defer rows.Close()

	var userIDs []int
	for rows.Next() {
		var userID int
		err := rows.Scan(&userID)
		if err != nil {
			return []int{}, err
		}
		userIDs = append(userIDs, userID)
	}

	if err := rows.Err(); err != nil {
		return []int{}, err
	}

	return userIDs, nil
}

func CheckIfMember(userId, groupId int) (STATUS, error) {
	var status STATUS
	query := "SELECT role,joined_at FROM group_members WHERE user_id=? AND group_id=?"
	db := dataBase.GetDb()

	err := db.QueryRow(query, userId, groupId).Scan(&status.Status, &status.Since)
	if errors.Is(err, sql.ErrNoRows) {
		err = pkg.ErrNotMember
	}
	return status, err
}

func Getusernamebyid(id int) (string, error) {
	db := dataBase.GetDb()
	var nickname string
	err := db.QueryRow(`SELECT lastName FROM users WHERE id=?`, id).Scan(&nickname)
	if err != nil {
		return "", err
	}
	return nickname, nil
}

// func GetPrivatemessagescount(user1, user2 string, db *sql.DB) (int, error) {
// 	var count int
// 	err := db.QueryRow(`
// 		SELECT COUNT(*)
// 		FROM users_messages
// 		WHERE (sender_id = ? AND receiver_id = ? OR receiver_id = ? AND sender_id = ?)`,
// 		user1, user2, user1, user2).Scan(&count)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return count, nil
// }

// func GetGroupmessagescount(user1, groupID string, db *sql.DB) (int, error) {
// 	var count int
// 	err := db.QueryRow(`
// 		SELECT COUNT(*)
// 		FROM groups_messages
// 		WHERE (sender_id = ? AND group_id = ? OR group_id = ? AND sender_id = ?)`,
// 		user1, groupID, user1, groupID).Scan(&count)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return count, nil
// }

// // func getRows(rows *sql.Rows) ([]Chat, error) {
// // 	var Chats []Chat
// // 	for rows.Next() {
// // 		var sender, message, date string
// // 		err := rows.Scan(&sender, &message, &date)
// // 		if err != nil {
// // 			return nil, err
// // 		}
// // 		messages := Chat{
// // 			Sender:     sender,
// // 			Content:    message,
// // 			Created_at: date,
// // 		}
// // 		Chats = append(Chats, messages)
// // 	}
// // 	if err := rows.Err(); err != nil {
// // 		return nil, err
// // 	}
// // 	return Chats, nil
// // }
