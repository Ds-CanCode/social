package funcs

import "net/http"

func Rootes(mux *http.ServeMux) {
	mux.HandleFunc("/api/notif/invitations", GetInvitaion)
	mux.HandleFunc("/api/notif/requestgroup", GetRequestGroup)
	mux.HandleFunc("/api/notif/event", GetEvent)
	mux.HandleFunc("/api/notif/requestuser", GetRequestUser)
}
