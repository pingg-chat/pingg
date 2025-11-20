package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pingg-chat/pingg/config"
	"github.com/pingg-chat/pingg/models"
	"github.com/pingg-chat/pingg/tui"
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

	// Run TUI
	if err := tui.Run(&user); err != nil {
		utils.Dd("Error running TUI", err)
	}
}
