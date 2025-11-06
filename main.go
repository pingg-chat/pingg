package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pingg-chat/pingg/config"
	"github.com/pingg-chat/pingg/models"
	"github.com/pingg-chat/pingg/utils"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <user_id>")
		return
	}

	id, err := strconv.ParseInt(os.Args[1], 10, 64)

	if err != nil {
		fmt.Println("Invalid user ID:", os.Args[1])
		return
	}

	config.SetID(id)

	user := models.User{ID: id}
	user.Load()

	utils.Dd(user)
}
