package utils

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

func GenerateRandomID() string {
	u := uuid.New()
	return u.String()
}

func GenerateSixFigure() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06v", rand.Int31n(1000000))
}
