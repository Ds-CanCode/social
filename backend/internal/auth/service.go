package funcs

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	pkg "funcs/pkg"
)

func ValidateCredentials(user User) bool {
	if ParseAddress(user.Email) {
		return false
	}
	if ValidateLength(user.FirstName) {
		return false
	}
	if ValidateLength(user.LastName) {
		return false
	}
	if ValidateLength(user.Nickname) {
		return false
	}
	if ValidateLength(user.Password) {
		return false
	}
	// if user.Age > 100 || user.Age <= 18 {
	// 	return false
	// }
	return true
}

func ParseAddress(s string) bool {
	return false
}

func ValidateLength(data string) bool {
	if len(data) <= 3 || len(data) >= 32 {
		return true
	}
	return false
}

func validateUser(user User) error {
	if user.Email == "" {
		return fmt.Errorf("l'email est obligatoire")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(user.Email) {
		return fmt.Errorf("format d'email invalide")
	}
	if len(user.Password) < 6 {
		return fmt.Errorf("le mot de passe doit contenir au moins 6 caractères")
	}

	if user.FirstName == "" || strings.TrimSpace(user.FirstName) == "" {
		return fmt.Errorf("le prénom est obligatoire")
	}
	if user.LastName == "" || strings.TrimSpace(user.LastName) == "" {
		return fmt.Errorf("le nom est obligatoire")
	}
	if user.DateOfBirth != "" {
		_, err := time.Parse("2006-01-02", user.DateOfBirth)
		if err != nil {
			return fmt.Errorf("format de date incorrect (attendu: YYYY-MM-DD)")
		}
	}
	return nil
}

func ReadUserInfo(w http.ResponseWriter, r *http.Request) (User, int, error) {
	// Extract form values
	user := User{
		Email:       r.FormValue("email"),
		Password:    r.FormValue("password"),
		FirstName:   r.FormValue("firstName"),
		LastName:    r.FormValue("lastName"),
		DateOfBirth: r.FormValue("dateOfBirth"),
		Nickname:    r.FormValue("nickName"),
		AboutMe:     r.FormValue("aboutMe"),
	}
	if err := validateUser(user); err != nil {
		return user, http.StatusBadRequest, err
	}

	file, header, err := r.FormFile("avatar")

	if err == nil || err.Error() != "http: no such file" {
		if err == nil {

			imgName, err := pkg.SaveImage(file, header)
			user.Avatar = imgName
			if err != nil {
				return user, http.StatusInternalServerError, pkg.ErrProccessingFile
			}
		} else {
			return user, http.StatusInternalServerError, pkg.ErrInvalidFile
		}
	}

	return user, 0, nil
}
