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
	// Checking if user id is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <user_id>")
		return
	}

	id, err := strconv.ParseInt(os.Args[1], 10, 64)

	if err != nil {
		fmt.Println("Invalid user ID:", os.Args[1])
		return
	}

	// Setting user id on config
	config.SetID(id)

	// Load user from api
	user := models.User{ID: id}
	err = user.Load()

	if err != nil {
		utils.Dd("Error on Load User", err)
	}
	// travar o processo ... só sair com CTRL+C
	fmt.Println("Press CTRL+C to exit")

	utils.WaitForCtrlC()

	utils.Dump("User:", user)
}
