package funcs

import (
	"net/http"

	pkg "funcs/pkg"
)

func Search(w http.ResponseWriter, r *http.Request) {
	_, err := pkg.GetIdBySession(w, r)
	if err != nil {
		pkg.SendResponseStatus(w, http.StatusUnauthorized, pkg.ErrUnauthorized)
		return
	}
	state := r.FormValue("type")
	input := r.FormValue("query")

	if state == "people" && input != "" {
		users, err := GetUsers(input)
		if err != nil {
			pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
			return
		}

		err = pkg.Encode(w, &users)
		if err != nil {
			pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
			return
		}

	} else if state == "groups" && input != "" {
		groups, err := GetGroups(input)
		if err != nil {
			pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
			return
		}

		err = pkg.Encode(w, &groups)
		if err != nil {
			pkg.SendResponseStatus(w, http.StatusInternalServerError, err)
			return
		}

	} else {
		pkg.SendResponseStatus(w, http.StatusBadRequest, pkg.ErrBadRequest)
		return
	}
}
