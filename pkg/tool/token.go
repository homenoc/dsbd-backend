package tool

import (
	"github.com/google/uuid"
	"log"
)

func GenerateToken(count int) (string, error) {
	var token string
	for i := 0; i < count; i++ {
		tmpToken, err := uuid.NewRandom()
		if err != nil {
			log.Printf("error: token generate |%s", err)
			return "", err
		}
		token += tmpToken.String()
	}
	return token, nil
}
