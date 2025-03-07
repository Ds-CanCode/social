package funcs

import (
	"database/sql"

	database "funcs/internal/database"
)

func GetRequestUserDB(id int) ([]NOTIF, error) {
	var notif []NOTIF
	query := `SELECT user1 ,accepted ,is_read ,createdAt FROM folowers WHERE user2=?`
	db := database.GetDb()

	rows, err := db.Query(query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return notif, nil
		}
		return notif, err
	}
	defer rows.Close()

	for rows.Next() {
		var n NOTIF
		err = rows.Scan(&n.Sender.Id, &n.Accepted, &n.IsRead, &n.CreatedAt)
		if err != nil {
			return notif, err
		}

		n.Sender, err = GetuserInfo(n.Sender.Id)
		if err != nil {
			return notif, err
		}

		notif = append(notif, n)
	}
	return notif, nil
}

func GetRequestGroupDB(id int) ([]NOTIF, error) {
	var notif []NOTIF
	query := `SELECT r.group_id, r.user_id, r.requested_at , r.status FROM join_requests r 
				JOIN groups g ON r.group_id=g.id 
				WHERE g.creator_id=?
				ORDER BY r.requested_at DESC;
			`
	db := database.GetDb()

	rows, err := db.Query(query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return notif, nil
		}
		return notif, err
	}
	defer rows.Close()

	for rows.Next() {
		var n NOTIF
		var status string
		err = rows.Scan(&n.Group.Id, &n.Sender.Id, &n.CreatedAt, &status)
		if err != nil {
			return notif, err
		}

		n.Sender, err = GetuserInfo(n.Sender.Id)
		if err != nil {
			return notif, err
		}

		n.Group, err = GetGroupInfo(n.Group.Id)
		if err != nil {
			return notif, err
		}

		n.Accepted = status == "accepted"
		notif = append(notif, n)
	}

	return notif, nil
}

func GetInvitaionDB(id int) ([]NOTIF, error) {
	var notif []NOTIF
	var status string
	query := `SELECT group_id, invited_by, invited_at ,status FROM invitations WHERE invited_user=?`
	db := database.GetDb()

	rows, err := db.Query(query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return notif, nil
		}
		return notif, err
	}
	defer rows.Close()

	for rows.Next() {
		var n NOTIF
		err = rows.Scan(&n.Group.Id, &n.Sender.Id, &n.CreatedAt, &status)
		if err != nil {
			return notif, err
		}

		n.Sender, err = GetuserInfo(n.Sender.Id)
		if err != nil {
			return notif, err
		}

		n.Group, err = GetGroupInfo(n.Group.Id)
		if err != nil {
			return notif, err
		}

		n.Accepted = status == "accepted"
		notif = append(notif, n)
	}

	return notif, nil
}

func GetuserInfo(id int) (SUser, error) {
	query := "SELECT id, firstName, lastName, nickname, avatar FROM users WHERE id=?"
	db := database.GetDb()
	var user SUser

	err := db.QueryRow(query, id).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Nickname, &user.Path)
	if err == sql.ErrNoRows {
		return user, nil
	}

	return user, err
}

func GetGroupInfo(id int) (SGroup, error) {
	query := "SELECT id, title, imagePath FROM groups WHERE id=?"
	db := database.GetDb()
	var group SGroup

	err := db.QueryRow(query, id).Scan(&group.Id, &group.Title, &group.Path)
	if err == sql.ErrNoRows {
		return group, nil
	}

	return group, err
}

func GetEventDB(id int) ([]NOTIF, error) {
	var notif []NOTIF
	query := `SELECT e.user_id, e.group_id, e.title, e.eventDATE FROM posts_event e
			JOIN group_members g ON g.id=e.group_id
			WHERE g.user_id=?
	`

	db := database.GetDb()
	rows, err := db.Query(query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return notif, nil
		}
		return notif, err
	}
	defer rows.Close()

	for rows.Next() {
		var n NOTIF
		err = rows.Scan(&n.Sender.Id, &n.Group.Id, &n.Title, &n.CreatedAt)
		if err != nil {
			return notif, err
		}

		n.Sender, err = GetuserInfo(n.Sender.Id)
		if err != nil {
			return notif, err
		}

		n.Group, err = GetGroupInfo(n.Group.Id)
		if err != nil {
			return notif, err
		}

		notif = append(notif, n)
	}

	return notif, nil
}

func ReadUserRequest(userId int) error {
	query := `UPDATE folowers SET is_read=1 WHERE user2=?`
	db := database.GetDb()

	_, err := db.Exec(query, userId)

	return err
}

func ReadInvitation(userId int) error {
	db := database.GetDb()
	query := `UPDATE invitations SET is_read=1 WHERE invited_user=?`
	_, err := db.Exec(query, userId)
	return err
}

func ReadJoinRequest(userId int) error {
	db := database.GetDb()
	query := `UPDATE join_requests SET is_read=1 WHERE user_id=?`
	_, err := db.Exec(query, userId)
	return err
}
