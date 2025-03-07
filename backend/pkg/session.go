package funcs

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore([]byte("super-secret-key"))

func SetSession(w http.ResponseWriter, r *http.Request, id int) error {
	session, err := Store.Get(r, "session")
	if err != nil {
		fmt.Println("Error retrieving session:", err)
		return ErrSessionNotFound
	}
	session.Values["user_id"] = id
	session.Options = &sessions.Options{
		MaxAge: 2147483647,  // Max possible value (~68 years)
		Domain: "localhost", // Set to your domain (in production, use your actual domain)
		Path:   "/",         // Make cookie available on all paths

	}
	return session.Save(r, w)
}

func GetIdBySession(w http.ResponseWriter, r *http.Request) (int, error) {
	session, err := Store.Get(r, "session")
	if err != nil {
		session.Options.MaxAge = -1
		session.Save(r, w)

		return 0, ErrSessionNotFound
	}

	id, ok := session.Values["user_id"].(int)
	if !ok {
		return 0, ErrInvalidUserID
	}
	return id, nil
}

func DeleteSession(w http.ResponseWriter, r *http.Request) error {
	// Try to retrieve the session
	session, err := Store.Get(r, "session")
	if err != nil {
		// If session retrieval fails, manually delete the cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   true, // Set to false if not using HTTPS
			MaxAge:   -1,   // Immediately expire
		})


		if err := session.Save(r, w); err != nil {
			return err
		}

		return nil

	}

	// Expire the session properly
	session.Options.MaxAge = -1

	if err := session.Save(r, w); err != nil {
		return err
	}

	return nil
}
