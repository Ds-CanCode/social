package funcs

import (
	"net/http"

	pkg "funcs/pkg"
)

func GetInvitaion(w http.ResponseWriter, r *http.Request) {
	id, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusUnauthorized, pkg.ErrUnauthorized)
		return
	}

	var notif []NOTIF
	notif, err = GetInvitaionDB(id)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
	err = pkg.Encode(w, &notif)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	err = ReadAll(4, id)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
}

func GetRequestGroup(w http.ResponseWriter, r *http.Request) {
	id, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusUnauthorized, pkg.ErrUnauthorized)
		return
	}

	var notif []NOTIF
	notif, err = GetRequestGroupDB(id)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	err = pkg.Encode(w, &notif)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	err = ReadAll(3, id)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
}

func GetRequestUser(w http.ResponseWriter, r *http.Request) {
	id, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusUnauthorized, pkg.ErrUnauthorized)
		return
	}

	var notif []NOTIF
	notif, err = GetRequestUserDB(id)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	err = pkg.Encode(w, &notif)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	err = ReadAll(2, id)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
}

func GetEvent(w http.ResponseWriter, r *http.Request) {
	id, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusUnauthorized, pkg.ErrUnauthorized)
		return
	}

	var notif []NOTIF
	notif, err = GetEventDB(id)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}

	err = pkg.Encode(w, &notif)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
		return
	}
}
