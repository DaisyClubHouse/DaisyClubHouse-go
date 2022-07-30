package utils

import "github.com/google/uuid"

func GenerateRandomID() string {
	u := uuid.New()
	return u.String()
}
