package funcs

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	pkg "funcs/pkg"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var user User
	err := pkg.Decode(r, &user)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusBadRequest, err)
		return
	}
	if ValidateLength(user.Email) || ValidateLength(user.Password) {
		fmt.Println("error validate length")
		pkg.SendResponseStatus(w, http.StatusBadRequest, pkg.ErrBadRequest)
		return
	}

	user, err = comparePassword(user.Email, user.Password)
	if err != nil {
		fmt.Println("error compare password")
		pkg.SendResponseStatus(w, http.StatusBadRequest, err)
		return
	}
	err = pkg.SetSession(w, r, user.ID)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	err := pkg.DeleteSession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func Register(w http.ResponseWriter, r *http.Request) {
	// Limit upload size to 2MB
	err := r.ParseMultipartForm(2 << 20) // 2MB limit
	if err != nil {
		log.Println("Error parsing form:", err)
		pkg.SendResponseStatus(w, http.StatusBadRequest, err)
		return
	}
	var user User
	var status int
	user, status, err = ReadUserInfo(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, status, err)
		return
	}

	// Save user to DB
	if err = AddUser(user); err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Inscription réussie"))
}


func CheckAuth(w http.ResponseWriter, r *http.Request) {
	_, err := pkg.GetIdBySession(w, r)
	if err != nil {

		pkg.DeleteSession(w, r)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]bool{"isAuthenticated": false})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"isAuthenticated": true})
}
