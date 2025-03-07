package funcs

import (
	"database/sql"
	"net/http"

	dataBase "funcs/internal/database"
	pkg "funcs/pkg"
)

func PostReactionAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	id, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
	var rectPost reactionPost
	err = pkg.Decode(r, &rectPost)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusBadRequest, err)
		return
	}
	rectPost.UserID = id
	if ReactPost(rectPost) != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	reactionCount, err := pkg.CountRect(rectPost.PostID)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = pkg.Encode(w, reactionCount)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
}

func CommentReactionAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	id, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
	var rectComment reactionComment
	err = pkg.Decode(r, &rectComment)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusBadRequest, err)
		return
	}
	rectComment.UserID = id
	if ReactComment(rectComment) != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
	reactionCommentCount, err := pkg.CountCommentRect(rectComment.CommentID)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = pkg.Encode(w, reactionCommentCount)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
}

func CommentReactionGet(w http.ResponseWriter, r *http.Request) {
	id, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	quer := r.URL.Query()
	commentID := quer.Get("commentId")
	if commentID == "" {
		http.Error(w, "Missing commentId parameter", http.StatusBadRequest)
		return
	}

	db := dataBase.GetDb()

	type Current struct {
		CurrentReaction string `json:"currentReact"`
	}

	var reactionCount Current
	query := `SELECT reaction_type FROM reactionsComments WHERE user_id = ? AND comment_id = ?`

	if id != 0 {
		err = db.QueryRow(query, id, commentID).Scan(&reactionCount.CurrentReaction)
		if err == sql.ErrNoRows {
			reactionCount.CurrentReaction = ""
			err = nil
		}
	}

	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = pkg.Encode(w, reactionCount)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
}

// func PostReactionGet(http.ResponseWriter, *http.Request) {
// }
