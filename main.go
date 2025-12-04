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
	// Usage: pingg <user_id> [width] [height]
	if len(os.Args) < 2 {
		fmt.Println("Usage: pingg <user_id> [width] [height]")
		return
	}

	id, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		fmt.Println("Invalid user ID:", os.Args[1])
		return
	}

	// Parse optional width and height arguments
	var width, height int
	if len(os.Args) >= 4 {
		if w, err := strconv.Atoi(os.Args[2]); err == nil && w > 0 {
			width = w
		}
		if h, err := strconv.Atoi(os.Args[3]); err == nil && h > 0 {
			height = h
		}
	}

	// Set terminal size if provided
	if width > 0 && height > 0 {
		tui.SetTerminalSize(width, height)
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
		utils.Dd("Error running TUI", err.Error())
	}
}
