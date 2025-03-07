package funcs

import (
	database "funcs/internal/database"
)

func GetUsers(input string) ([]SUser, error) {
	input = "%" + input + "%"
	var users []SUser

	query := "SELECT id, firstName, lastName, nickname, avatar FROM users WHERE (firstName LIKE ?) OR (lastName LIKE ?) OR (nickname LIKE ?)"

	db := database.GetDb()
	rows, err := db.Query(query, input, input, input)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user SUser
		err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Nickname, &user.Path)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return users, err
	}

	return users, nil
}

func GetGroups(input string) ([]SGroup, error) {
	input = "%" + input + "%"
	var groups []SGroup

	query := "SELECT id, title, imagePath FROM groups WHERE (title LIKE ?)"

	db := database.GetDb()
	rows, err := db.Query(query, input)
	if err != nil {
		return groups, err
	}
	defer rows.Close()

	for rows.Next() {
		var group SGroup
		err := rows.Scan(&group.Id, &group.Title, &group.Path)
		if err != nil {
			return groups, err
		}
		groups = append(groups, group)
	}

	if err = rows.Err(); err != nil {
		return groups, err
	}

	return groups, nil
}
