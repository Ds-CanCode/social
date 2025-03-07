package funcs

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	ws "funcs/internal/chat"
	database "funcs/internal/database"
	structure "funcs/internal/notification"
	pkg "funcs/pkg"
)

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	id, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusUnauthorized, err)
		return
	}

	profile, err := GetUserDB(id)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = pkg.Encode(w, profile)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
}

func getProfile(w http.ResponseWriter, r *http.Request) {
	id, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusUnauthorized, err)
		return
	}

	query := r.URL.Query()
	targetId, _ := strconv.Atoi(query.Get("profileId"))
	if targetId == id {
		http.Error(w, "see other", http.StatusSeeOther)
		return
	}

	profile, err := GetProfileDB(targetId, id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "see other", http.StatusSeeOther)
		}else {
			pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		}
		return
	}

	err = pkg.Encode(w, profile)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
}

func request(w http.ResponseWriter, r *http.Request) {
	id, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusUnauthorized, err)
		return
	}

	query := r.URL.Query()
	targetId, _ := strconv.Atoi(query.Get("profileId"))
	if targetId == 0 {
		pkg.SendResponseStatus(w, http.StatusBadRequest, pkg.ErrInvalidNamber)
		return
	}
	state, err := CheckFollowstatus(id, targetId)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, pkg.ErrInvalidNamber)
		return
	}

	if !state {
		err = Follow(id, targetId)
		if err != nil {
			fmt.Println("can't follow", err)
			pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
			return
		}
		userInfo, _ := structure.GetuserInfo(targetId)
		var follownotif structure.NOTIF = structure.NOTIF{
			Type:   "follow",
			Sender: structure.SUser(userInfo),
		}
		ws.SendRealTimeNotification([]int{targetId}, follownotif)
	} else {
		err = UnFollow(id, targetId)
		if err != nil {
			pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func handlerequest(w http.ResponseWriter, r *http.Request) {
	id, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusUnauthorized, err)
		return
	}
	query := r.URL.Query()
	targetId, _ := strconv.Atoi(query.Get("profileId"))
	if targetId == 0 {
		pkg.SendResponseStatus(w, http.StatusBadRequest, pkg.ErrInvalidNamber)
		return
	}
	state := query.Get("state")
	if state == "accepted" {
		AcceptRequest(id, targetId)
	} else if state == "rejected" {
		RejectRequest(id, targetId)
	} else {
		pkg.SendResponseStatus(w, http.StatusBadRequest, pkg.ErrInvalidNamber)
		return
	}
}

func ChangeProfileType(w http.ResponseWriter, r *http.Request) {
	userId, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusUnauthorized, err)
		return
	}
	state := r.FormValue("type")
	if state == "public" {
		ChangeProfileTypeDB(userId, 0)
	} else if state == "private" {
		ChangeProfileTypeDB(userId, 1)
	} else {
		pkg.SendResponseStatus(w, http.StatusBadRequest, pkg.ErrInvalidNamber)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetFollowers(w http.ResponseWriter, r *http.Request) {
	id, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusUnauthorized, err)
		return
	}
	var followers []Followers
	followers, err = getFollowersDB(id)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	err = pkg.Encode(w, followers)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
}

func GetUserFollowers(w http.ResponseWriter, r *http.Request) {
	id, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusUnauthorized, err)
		return
	}

	targetId, _ := strconv.Atoi(r.URL.Query().Get("profileId"))
	if targetId == 0 {
		targetId = id
	}

	if targetId != id {

		bl, err := GetProfileType(targetId)
		if err != nil {
			pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
			return
		}

		if bl {
			is, err := IsFriends(id, targetId)
			if err != nil {
				pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
				return
			}

			if is == 0 {
				pkg.SendResponseStatus(w, http.StatusForbidden, err)
				return
			}
		}
	}

	followers, err := GetUserFollowersDB(targetId)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	err = pkg.Encode(w, &followers)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
}

func GetUserFollowing(w http.ResponseWriter, r *http.Request) {
	id, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusUnauthorized, err)
		return
	}

	targetId, _ := strconv.Atoi(r.URL.Query().Get("profileId"))
	if targetId == 0 {
		targetId = id
	}

	if targetId != id {

		bl, err := GetProfileType(targetId)
		if err != nil {
			pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
			return
		}

		if bl {
			is, err := IsFriends(id, targetId)
			if err != nil {
				pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
				return
			}

			if is == 0 {
				pkg.SendResponseStatus(w, http.StatusForbidden, err)
				return
			}
		}
	}

	followers, err := GetUserFollowingDB(targetId)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	err = pkg.Encode(w, &followers)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
}

func NowRequest(w http.ResponseWriter, r *http.Request) {
	id, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusUnauthorized, err)
		return
	}
	targetId, _ := strconv.Atoi(r.URL.Query().Get("profileId"))
	if targetId == 0 {
		targetId = id
	}
	var nowF struct {
		Is bool `json:"is"`
	}
	result, err := GetNowRequest(targetId, id)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
	nowF.Is = result
	err = pkg.Encode(w, &nowF)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
}

func GetUserFollowersDB(id int) ([]Followers, error) {
	query1 := `
		SELECT users.id, users.firstName, users.lastName , users.avatar
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

func GetUserFollowingDB(id int) ([]Followers, error) {
	query1 := `
		SELECT users.id, users.firstName, users.lastName, users.avatar
		FROM users
		JOIN folowers ON users.id = folowers.user2
		WHERE folowers.user1 = ? AND folowers.accepted = 1;
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
