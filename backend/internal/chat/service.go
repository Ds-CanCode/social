package funcs

import (
	"database/sql"

	dataBase "funcs/internal/database"
)

func GetHistoryPrv(usr1, usr2, offset int) ([]Private_Message, error) {
	db := dataBase.GetDb()
	var message []Private_Message
	query := `SELECT * FROM users_messages 
	WHERE (sender_id=? AND receiver_id=?) OR (sender_id=? AND receiver_id=?)
	ORDER BY createdAt DESC
	LIMIT 10 OFFSET ?;
	`

	rows, err := db.Query(query, usr1, usr2, usr2, usr1, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return message, err
	}
	defer rows.Close()

	for rows.Next() {
		var m Private_Message

		err = rows.Scan(&m.Id, &m.SenderID, &m.ReceiverID, &m.Content, &m.Created_at)
		if err != nil {
			return message, err
		}
		m.Type = "Private"

		message = append(message, m)
	}

	return message, nil
}

func GetGroupHistoryPrv(id, groupId, offset int) ([]Group_Mesaage, error) {
	var messages []Group_Mesaage
	query := `SELECT * FROM groups_messages 
				WHERE group_id=?
				ORDER BY createdAt DESC
				LIMIT 10 OFFSET ?;
				`
	db := dataBase.GetDb()

	rows, err := db.Query(query, groupId, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		return messages, err
	}
	defer rows.Close()

	for rows.Next() {
		var m Group_Mesaage

		err = rows.Scan(&m.Id, &m.SenderID, &m.Group, &m.Content, &m.Created_at)
		if err != nil {
			return messages, err
		}

		m.Sender, err = GetuserInfo(m.SenderID)
		if err != nil {
			return messages, err
		}
		if m.Sender.Id == id {
m.Vv = "You"
		} else {
			m.Vv= "Other"
		}

		
		m.Type = "Group"

		messages = append(messages, m)
	}

	return messages, nil
}

func GetuserInfo(id int) (SUser, error) {
	query := "SELECT id, firstName, lastName, nickname, avatar FROM users WHERE id=?"
	db := dataBase.GetDb()
	var user SUser

	err := db.QueryRow(query, id).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Nickname, &user.Path)
	if err == sql.ErrNoRows {
		return user, nil
	}

	return user, err
}
