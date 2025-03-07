package funcs

import "net/http"

func Rootes(mux *http.ServeMux) {
	mux.HandleFunc("/api/post/add", Post)
	mux.HandleFunc("/api/post/getAll", GetPosts)
	mux.HandleFunc("/api/post/get", GetGroupPost)
	mux.HandleFunc("/api/post/getuserpost", GetuserPost)
}
