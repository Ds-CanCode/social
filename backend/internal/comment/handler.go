package funcs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	pkg "funcs/pkg"
)

func AddComments(w http.ResponseWriter, r *http.Request) {
	id, err := pkg.GetIdBySession(w, r)
	if err != nil {
		fmt.Println("AddComments111", err)
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	var comment Comment
	var status int
	comment, status, err = ReadComment(w, r)
	if err != nil {
		fmt.Println("AddComments")
		pkg.SendResponseStatus(w, status, err)
		return
	}
	comment.UserID = id

	insertedComment, err := AddComment(comment)
	if err != nil {
		fmt.Println("AddComment failed")
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(insertedComment) // Send the inserted comment as response
}

func GetComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	query := r.URL.Query()
	postID := query.Get("PostId")
	postId, err := strconv.Atoi(postID)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusBadRequest, err)
		return
	}
	comments, err := GetComments(postId)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	err = pkg.Encode(w, comments)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusBadRequest, err)
		return
	}
}
