package funcs

import (
	"net/http"
	"os"
	"path/filepath"
)

func Serveimages(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	filePath := filepath.Join("images", path)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Fichier introuvable", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, filePath)
}
