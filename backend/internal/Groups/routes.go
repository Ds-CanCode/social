package funcs

import "net/http"

func Rootes(mux *http.ServeMux) {
	// mux.HandleFunc("/api/groups/countmembers", CountMembers)

	mux.HandleFunc("/api/groups/add", CreateGroup)
	mux.HandleFunc("/api/groups/addevent", CreateEvent)

	mux.HandleFunc("/api/groups/get", GroupInfo)
	mux.HandleFunc("/api/groups/getall", GetAllGroups)
	mux.HandleFunc("/api/groups/getusergroups", GetuserGroup)
	mux.HandleFunc("/api/groups/getevents", GetEvents)
	mux.HandleFunc("/api/groups/handleJoinEvent", HandleJoinEvent)
	mux.HandleFunc("/api/groups/handleDeleteEvent", LeaveGroup)
	// mux.HandleFunc("/api/groups/counteventattends", CountEventAttends)


	mux.HandleFunc("/api/groups/invitetogroup", InviteToGroup)
	mux.HandleFunc("/api/groups/request", RequestGroup)

	mux.HandleFunc("/api/groups/acceptinvitation", AcceptInvitation)
	mux.HandleFunc("/api/groups/acceptrequest", AcceptRequest)

	mux.HandleFunc("/api/groups/rejectrequest", RejectRequest)
}
