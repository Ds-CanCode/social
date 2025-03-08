package funcs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	ws "funcs/internal/chat"
	structure "funcs/internal/notification"
	pkg "funcs/pkg"
)

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	var group Group
	var err error
	var userId int

	userId, err = pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusUnauthorized, err)
		return
	}

	group, status, err := readGroupInfo(r, userId)
	if err != nil {
		pkg.SendResponseStatus(w, status, err)
		return
	}

	if InsertGroupInfo(group) != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	// send succes  json message to user
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	id, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusUnauthorized, err)
		return
	}
	var event Event
	event, err = ReadEventInfo(r)
	if err != nil {
		fmt.Println("err", err)
		pkg.SendResponseStatus(w, http.StatusBadRequest, err)
		return
	}
	event.CreatorId = id

	isMember, err := IsFollowing(event.CreatorId, event.GroupId)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
	if !isMember {
		pkg.SendResponseStatus(w, http.StatusForbidden, pkg.ErrNotMember)
		return
	}
	event.Id, err = InsertEvent(event)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	/********************************************/
	/****Send Notif to all member except creator*/
	/********************************************/
}

func GroupInfo(w http.ResponseWriter, r *http.Request) {
	userId, err := pkg.GetIdBySession(w, r)
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

	exist, err := CheckGroupExists(groupId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !exist {
		http.Error(w, "Group not found", http.StatusBadRequest)
		return
	}

	var group Group
	group, err = GetGroupInfo(groupId)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	group.MemberStatus, err = GetMemberStatus(userId, groupId)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	err = pkg.Encode(w, &group)
	if err != nil {
		fmt.Println(err)
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
}

func GetuserGroup(w http.ResponseWriter, r *http.Request) {
	userId, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusUnauthorized, err)
		return
	}

	groups, err := GetuserGroupDB(userId)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	err = pkg.Encode(w, &groups)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
}

func InviteToGroup(w http.ResponseWriter, r *http.Request) {
	userId, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusUnauthorized, err)
		return
	}

	err = r.ParseMultipartForm(2 << 20) // 2 MB limit for uploaded files
	if err != nil {
		fmt.Println(err)
		return
	}
	users := r.Form["users"]
	group := r.FormValue("group")

	groupId, _ := strconv.Atoi(group)
	if groupId == 0 {
		pkg.SendResponseStatus(w, http.StatusBadRequest, pkg.ErrInvalidNamber)
		return
	}

	for i := range users {
		id, _ := strconv.Atoi(users[i])
		if id == 0 {
			pkg.SendResponseStatus(w, http.StatusBadRequest, pkg.ErrInvalidNamber)
			return
		}

		err = InviteUser(userId, id, groupId)
		if err != nil {
			pkg.SendResponseStatus(w, http.StatusInternalServerError, pkg.ErrInvalidNamber)
			return
		}

		userInfo, _ := structure.GetuserInfo(id)
		var follownotif structure.NOTIF = structure.NOTIF{
			Type:   "groupInvitation",
			Sender: structure.SUser(userInfo),
		}
		ws.SendRealTimeNotification([]int{id}, follownotif)

	}
}

func AcceptInvitation(w http.ResponseWriter, r *http.Request) {
	userId, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusUnauthorized, err)
		return
	}

	groupId, _ := strconv.Atoi(r.URL.Query().Get("groupId"))
	if groupId == 0 {
		pkg.SendResponseStatus(w, http.StatusBadRequest, pkg.ErrInvalidNamber)
		return
	}

	_, err = CheckIfMember(userId, groupId)
	if err != nil {
		if err != pkg.ErrNotMember {
			pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
			return
		}
	}

	_, err = CheckIfInveted(userId, groupId)
	if err != nil {
		if err == pkg.ErrNoInvitation {
			pkg.SendResponseStatus(w, http.StatusBadRequest, pkg.ErrNoInvitation)
		} else {
			pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		}
		return
	}

	err = Accept(userId, groupId)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
}

func AcceptRequest(w http.ResponseWriter, r *http.Request) {
	userId, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusUnauthorized, err)
		return
	}

	groupId, _ := strconv.Atoi(r.URL.Query().Get("groupId"))
	if groupId == 0 {
		pkg.SendResponseStatus(w, http.StatusBadRequest, pkg.ErrInvalidNamber)
		return
	}

	isCreator, err := CheckIfCreator(userId, groupId)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
	if !isCreator {
		fmt.Println("err1", err)
		pkg.SendResponseStatus(w, http.StatusBadRequest, pkg.ErrNotCreator)
		return
	}

	target, _ := strconv.Atoi(r.URL.Query().Get("target"))
	if groupId == 0 {
		pkg.SendResponseStatus(w, http.StatusBadRequest, pkg.ErrInvalidNamber)
		return
	}

	_, err = CheckIfMember(userId, groupId)
	if err != nil {
		if err == pkg.ErrNotMember {
			pkg.SendResponseStatus(w, http.StatusBadRequest, pkg.ErrNotMember)
		} else {
			pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		}
		return
	}

	_, err = CheckIfSendRequest(target, groupId)
	if err != nil {
		if err == pkg.ErrNoRequest {
			pkg.SendResponseStatus(w, http.StatusBadRequest, pkg.ErrNoRequest)
		} else {
			pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		}
		return
	}

	err = Accept(target, groupId)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
}

func RequestGroup(w http.ResponseWriter, r *http.Request) {
	userId, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusUnauthorized, err)
		return
	}

	groupId, _ := strconv.Atoi(r.URL.Query().Get("groupId"))
	if groupId == 0 {
		pkg.SendResponseStatus(w, http.StatusBadRequest, err)
		return
	}

	status, err := CheckIfMember(userId, groupId)
	if status.Status == "member" || status.Status == "creator" {
		pkg.SendResponseStatus(w, http.StatusBadRequest, pkg.ErrAlreadyMember)
		return
	}
	if err != nil && err != pkg.ErrNotMember {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, pkg.ErrAlreadyMember)
		return
	}

	err = SendRequest(userId, groupId)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	creator, err := GetGroupCreator(groupId)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	var follownotif structure.NOTIF = structure.NOTIF{
		Type:   "groupRequest",
		Sender: structure.SUser(creator),
	}
	ws.SendRealTimeNotification([]int{creator.Id}, follownotif)
}

func GetAllGroups(w http.ResponseWriter, r *http.Request) {
	userId, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusUnauthorized, err)
		return
	}

	groups, err := GetJoinedGroups(userId)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	err = pkg.Encode(w, &groups)
	if err != nil {
		fmt.Println(err)
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
}

func RejectRequest(w http.ResponseWriter, r *http.Request) {
	_, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusUnauthorized, err)
		return
	}

	groupId, _ := strconv.Atoi(r.URL.Query().Get("groupId"))
	if groupId == 0 {
		pkg.SendResponseStatus(w, http.StatusBadRequest, err)
		return
	}

	targetId, _ := strconv.Atoi(r.URL.Query().Get("target"))
	if targetId == 0 {
		pkg.SendResponseStatus(w, http.StatusBadRequest, err)
		return
	}
	

	err = Reject(targetId, groupId)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusBadRequest, err)
		return
	}
}

func GetEvents(w http.ResponseWriter, r *http.Request) {
	_, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusUnauthorized, err)
		return
	}

	group_id, err := strconv.Atoi(r.URL.Query().Get("groupId"))
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusBadRequest, err)
		return
	}

	events, err := GetGroupEvents(group_id)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	err = pkg.Encode(w, &events)
	if err != nil {
		fmt.Println(err)
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
}

func HandleJoinEvent(w http.ResponseWriter, r *http.Request) {
	member_id, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusUnauthorized, err)
		return
	}

	// Parse group and event IDs from query parameters
	groupId, err := strconv.Atoi(r.URL.Query().Get("groupId"))
	if err != nil || groupId == 0 {
		pkg.SendResponseStatus(w, http.StatusBadRequest, fmt.Errorf("invalid group ID"))
		return
	}

	eventId, err := strconv.Atoi(r.URL.Query().Get("eventId"))
	if err != nil || eventId == 0 {
		pkg.SendResponseStatus(w, http.StatusBadRequest, fmt.Errorf("invalid event ID"))
		return
	}

	// Check if user is already registered for this event
	exists, err := CheckEventRegistration(member_id, groupId, eventId)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	if exists {
		pkg.SendResponseStatus(w, http.StatusConflict, fmt.Errorf("already registered for this event"))
		return
	}

	// Join the event
	err = JoinEvent(member_id, groupId, eventId)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	// // Get updated attendees count
	attendeesCount, err := GetEventAttendeesCount(eventId)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	// Send success response with attendees count
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":         true,
		"attendees_count": attendeesCount,
	})
}

// delete the attending of the event
func LeaveGroup(w http.ResponseWriter, r *http.Request) {
	member_id, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusUnauthorized, err)
		return
	}

	groupId, err := strconv.Atoi(r.URL.Query().Get("groupId"))
	if err != nil || groupId == 0 {
		pkg.SendResponseStatus(w, http.StatusBadRequest, fmt.Errorf("invalid group ID"))
		return
	}

	eventId, err := strconv.Atoi(r.URL.Query().Get("eventId"))
	if err != nil || eventId == 0 {
		pkg.SendResponseStatus(w, http.StatusBadRequest, fmt.Errorf("invalid event ID"))
		return
	}

	exists, err := CheckEventRegistration(member_id, groupId, eventId)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	if !exists {
		pkg.SendResponseStatus(w, http.StatusConflict, fmt.Errorf("you are not registered for this event"))
		return
	}

	err = DeleteAttending(member_id, groupId, eventId)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
	attendeesCount, err := GetEventAttendeesCount(eventId)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":         true,
		"attendees_count": attendeesCount,
	})
}

// func CountEventAttends(w http.ResponseWriter, r *http.Request) {
// 	_, err := pkg.GetIdBySession(w, r)
// 	if err != nil {
// 		pkg.SendResponseStatus(w, http.StatusUnauthorized, err)
// 		return
// 	}

// 	eventId, err := strconv.Atoi(r.URL.Query().Get("eventId"))
// 	if err != nil || eventId == 0 {
// 		pkg.SendResponseStatus(w, http.StatusBadRequest, fmt.Errorf("invalid event ID"))
// 		return
// 	}

// 	attendeesCount, err := GetEventAttendeesCount(eventId)
// 	if err != nil {
// 		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"success": true,
// 		"attendees_count": attendeesCount,
// 	})

// }
