package funcs

import (
	"net/http"

	pkg "funcs/pkg"
)

func IsLoggin(funcNext http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := pkg.GetIdBySession(w, r)
		if err != nil {
			pkg.DeleteSession(w, r)
			pkg.SendResponseStatus(w, http.StatusUnauthorized, err)
			return
		} else {
			funcNext(w, r)
		}
	}
}
