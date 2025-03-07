package funcs

import (
	"net/http"
)

func Rootes(mux *http.ServeMux) {
	mux.HandleFunc("/api/login", Login)
	mux.HandleFunc("/api/register", Register)
	mux.HandleFunc("/api/check-auth", CheckAuth)
	mux.HandleFunc("/api/logout", Logout)
}
