package gen

import (
	uuidLib "github.com/google/uuid"
)

func GenerateUUIDString() (string, error) {
	uuid, err := uuidLib.NewRandom()
	if err != nil {
		return "", err
	}
	return uuid.String(), nil
}
