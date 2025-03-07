package funcs

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	dataBase "funcs/internal/database"

	"golang.org/x/crypto/bcrypt"
)

func comparePassword(email, password string) (User, error) {
	db := dataBase.GetDb()

	var user User

	query := "SELECT id, email, password, firstName, lastName, nickname FROM users WHERE  email=?"
	row := db.QueryRow(query, email)

	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Nickname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errors.New("user not found")
		}
		return user, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, errors.New("incorrect password")
	}

	user.Password = ""

	return user, nil
}

func AddUser(user User) error {
	Db := dataBase.GetDb()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = Db.Exec("INSERT INTO users (email, password, firstName, lastName, datebirth, avatar, nickname, aboutme, profileType ,createdAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		user.Email, string(hashedPassword), user.FirstName, user.LastName, user.DateOfBirth, user.Avatar, user.Nickname, user.AboutMe, user.ProfileType, time.Now())
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return errors.New("email is already taken")
		}
		return errors.New("an error occurred, please try again later")
	}

	return nil
}

