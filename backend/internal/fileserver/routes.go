package funcs

import "net/http"

func Rootes(mux *http.ServeMux) {
	mux.HandleFunc("/api/images", Serveimages)
}
