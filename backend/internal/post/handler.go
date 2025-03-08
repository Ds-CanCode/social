package funcs

import (
	"fmt"
	"net/http"
	"strconv"

	is "funcs/internal/Groups"
	pkg "funcs/pkg"
)

func Post(w http.ResponseWriter, r *http.Request) {
	var post POST
	var err error
	var status, id int

	id, err = pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusUnauthorized, err)
		return
	}

	post, status, err = ReadPost(w, r, id)
	if err != nil {
		pkg.SendResponseStatus(w, status, err)
		return
	}

	err = InsertPost(post)
	if err != nil {
		fmt.Println("err", err)
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
	id, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusUnauthorized, err)
		return
	}
	query := r.URL.Query()

	offset, _ := strconv.Atoi(query.Get("offset"))
	if offset < 0 {
		pkg.SendResponseStatus(w, http.StatusBadRequest, pkg.ErrInvalidNamber)
		return
	}
	fmt.Println(offset)
	var posts []POST
	posts, err = getFeedPosts(id, offset )
	if err != nil {
		fmt.Println(err )
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	fmt.Println(posts)
	err = pkg.Encode(w, &posts)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, pkg.ErrInvalidNamber)
		return
	}
}

func GetGroupPost(w http.ResponseWriter, r *http.Request) {
	id, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusUnauthorized, err)
		return
	}

	query := r.URL.Query()
	groupId, _ := strconv.Atoi(query.Get("groupId"))
	if groupId == 0 {
		pkg.SendResponseStatus(w, http.StatusBadRequest, pkg.ErrInvalidNamber)
		return
	}

	isFollow, err := is.IsFollowing(id, groupId)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
	if !isFollow {
		pkg.SendResponseStatus(w, http.StatusBadRequest, pkg.ErrNotFollowing)
		return
	}

	posts := GroupPost(id, groupId)
	err = pkg.Encode(w, &posts)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, pkg.ErrInvalidNamber)
		return
	}
}

func GetuserPost(w http.ResponseWriter, r *http.Request) {
	id, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusUnauthorized, err)
		return
	}

	query := r.URL.Query()
	taget, _ := strconv.Atoi(query.Get("targetId"))
	if taget == 0 {
		taget = id
	}

	posts, status, err := GetPostForProfile(id, taget)
	if err != nil {
		pkg.SendResponseStatus(w, status, err)
		return
	}

	err = pkg.Encode(w, &posts)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, pkg.ErrInvalidNamber)
		return
	}
}
