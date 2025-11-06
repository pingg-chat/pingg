package main

import (
	"fmt"
	"os"
	"strconv"
)

type User struct {
	ID       int64
	Name     string
	Username string
	Icon     string
	Email    string
}

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

	user := User{ID: id}

	dd(user)
}
