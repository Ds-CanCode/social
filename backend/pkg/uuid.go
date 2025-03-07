package funcs

import (
	"log"

	"github.com/gofrs/uuid"
)

func GenereteTocken() (string, error) {
	// Create a Version 4 UUID.
	u2, err := uuid.NewV4()
	if err != nil {
		log.Printf("failed to generate UUID: %v", err)
		return "", err
	}
	return u2.String(), nil
}
