package gen

import (
	"fmt"
	"github.com/google/uuid"
)

func GenerateUUID() string {
	u, err := uuid.NewRandom()
	if err != nil {
		fmt.Println(err)
		//return ""
	}
	uu := u.String()

	return uu
}
