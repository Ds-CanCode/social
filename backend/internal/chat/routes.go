package funcs

import "net/http"

func Rootes(mux *http.ServeMux) {
	mux.HandleFunc("/api/ws", ChatService)
	mux.HandleFunc("/api/chathistory", ChatHistory)
	mux.HandleFunc("/api/chatgrouphistory", ChatGroupHistory)
}
