package funcs

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	// session errors
	ErrSessionNotFound = errors.New("session not found")
	ErrInvalidUserID   = errors.New("invalid or missing user_id in session")
	ErrBadRequest      = errors.New("please check your credentials")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrNotFound        = errors.New("not found")
	ErrServer          = errors.New("error in server")
	ErrNotMember       = errors.New("you are not a member")
	ErrNotCreator      = errors.New("you are not the creator of group")
	ErrInvalidNamber   = errors.New("invalid Number")
	ErrNotFriends      = errors.New("not friends")
	ErrAlreadyFriend   = errors.New("already friends")
	ErrTitlePost       = errors.New("title length must be between 1 and 100")
	ErrContentPost     = errors.New("content length must be between 1 and 1000")
	ErrTyoePost        = errors.New("type must be 1 2 3 4")
	ErrEventTime       = errors.New("event must have a deadline")
	ErrNotFollowing    = errors.New("you re not following this group")
	ErrMaxSizeImage    = errors.New("file size exceeds 1MB")
	ErrInvalidFile     = errors.New("invalid file type, only JPEG, PNG, GIF, and WEBP allowed")
	ErrProccessingFile = errors.New("error processing image please try again")
	ErrImagextension   = errors.New("image extension must be jpeg, png, gif")
	ErrNoPost          = errors.New("there is no Post")
	ErrContentLenght   = errors.New("error in lenght text")
	ErrNoInvitation    = errors.New("can't find invitation")
	ErrNoRequest       = errors.New("can't find request")
	ErrAlreadyMember   = errors.New("already member")
)

func SendResponseStatus(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	// If there's an error, send it as JSON
	if err != nil {
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// If there's no error, send an empty response or success message
	json.NewEncoder(w).Encode(map[string]string{"message": "Success"})
}
