package funcs

import (
	"errors"
)

func ReadAll(notifType, id int) error {
	switch notifType {
	case 2:
		// userRequest
		return ReadUserRequest(id)
	case 3:
		// grouprequest
		return ReadJoinRequest(id)
	case 4:
		// invitation
		return ReadInvitation(id)
	}
	return errors.New("cccccccccc")
}
