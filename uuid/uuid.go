package uuid

import (
	"fmt"

	"github.com/google/uuid"
)

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func NewV4String() string {
	ID, err := uuid.NewRandom()
	if err != nil {
		panic(fmt.Errorf("could not generate new random id"))
	}

	return ID.String()
}
