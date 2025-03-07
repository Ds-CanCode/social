package funcs

import (
	"encoding/json"
	"net/http"
)

func Decode(r *http.Request, receiver any) error {
	reader := json.NewDecoder(r.Body)
	err := reader.Decode(&receiver)
	if err != nil {
		return err
	}
	return nil
}

func Encode(w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(&data)
	if err != nil {
		return err
	}
	return nil
}

func GetOptionalFormValue(r *http.Request, key string) *string {
	value := r.FormValue(key)
	if value == "" {
		return nil
	}
	return &value
}