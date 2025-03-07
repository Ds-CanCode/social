package funcs

import (
	"net/http"
	"strconv"

	pkg "funcs/pkg"
)

func HandleFollow(userId, groupId int) error {
	var is bool
	var err error

	is, err = IsFollowing(userId, groupId)
	if err != nil {
		return pkg.ErrServer
	}

	if is {
		return UnFollow(userId, groupId)
	}

	// return Follow(userId, groupId)
	return nil
}

func readGroupInfo(r *http.Request, userId int) (Group, int, error) {
	group := Group{
		CreatorId:   userId,
		Title:       r.FormValue("title"),
		Descreption: r.FormValue("description"),
	}

	err := CheckGroupInfo(group)
	if err != nil {
		return group, http.StatusBadRequest, pkg.ErrContentLenght
	}

	file, header, err := r.FormFile("image")
	if err == nil || err.Error() != "http: no such file" {
		if err == nil {
			group.Path, err = pkg.SaveImage(file, header)
			if err != nil {
				return group, http.StatusInternalServerError, pkg.ErrProccessingFile
			}
		} else {
			return group, http.StatusInternalServerError, pkg.ErrInvalidFile
		}
	}

	return group, 0, nil
}

func CheckGroupInfo(group Group) error {
	if group.Descreption == "" || len(group.Descreption) > 500 {
		return pkg.ErrContentLenght
	}
	if group.Title == "" || len(group.Title) > 500 {
		return pkg.ErrContentLenght
	}
	return nil
}

func GetMemberStatus(userId, groupId int) (STATUS, error) {
	var status STATUS
	var err error
	status, err = CheckIfInveted(userId, groupId)
	if err == nil {
		status.Status = "pending"
		status.Sender, err = GetuserInfo(status.Sender.Id)
		return status, err
	} else if err == pkg.ErrNoInvitation {
		status, err = CheckIfMember(userId, groupId)
		if err == nil {
			return status, nil
		} else if err == pkg.ErrNotMember {
			status, err = CheckIfSendRequest(userId, groupId)
			if err == nil {
				status.Status = "joinRequest"
				return status, nil
			} else if err == pkg.ErrNoRequest {
				status.Status = pkg.ErrNoRequest.Error()
				return status, nil
			}
		}
	}
	return status, err
}

func ReadEventInfo(r *http.Request) (Event, error) {
	var event Event
	event.Title = r.FormValue("title")
	event.Descreption = r.FormValue("description")
	event.Time = r.FormValue("date")

	if event.Title == "" || len(event.Title) > 200 {
		return event, pkg.ErrContentLenght
	}

	if event.Descreption == "" || len(event.Descreption) > 200 {
		return event, pkg.ErrContentLenght
	}

	if event.Time == "" || len(event.Time) > 200 {
		return event, pkg.ErrContentLenght
	}

	event.GroupId, _ = strconv.Atoi(r.FormValue("groupId"))
	if event.GroupId == 0 {
		return event, pkg.ErrInvalidNamber
	}

	return event, nil
}

