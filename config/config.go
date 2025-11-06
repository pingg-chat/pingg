package config

import (
	"fmt"
)

var userID int64

func SetID(id int64) {
	userID = id
}

func GetID() string {
	return fmt.Sprintf("%d", userID)
}
