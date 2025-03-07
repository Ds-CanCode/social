package funcs

import "net/http"

func Rootes(mux *http.ServeMux) {
	mux.HandleFunc("/api/reactPost/add", PostReactionAdd)
	mux.HandleFunc("/api/reactComment/add", CommentReactionAdd)
	mux.HandleFunc("/api/reactComment/getCurrent", CommentReactionGet)
	//mux.HandleFunc("/api/reactPost/get", PostReactionGet)
}
