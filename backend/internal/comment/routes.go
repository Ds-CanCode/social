package funcs

import "net/http"

func Rootes(mux *http.ServeMux) {
	mux.HandleFunc("/api/comment/add", AddComments)
	mux.HandleFunc("/api/comment/get", GetComment)
}
