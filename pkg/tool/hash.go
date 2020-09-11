package tool

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func Generate(data string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(data), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
		return "", fmt.Errorf("%v", err)
	}
	return string(hash), nil
}

func Verify(data string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(data)) == nil
}
