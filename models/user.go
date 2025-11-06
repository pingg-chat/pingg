package models

import "github.com/pingg-chat/pingg/api"

type User struct {
	ID       int64
	Name     string
	Username string
	Icon     string
	Email    string
}

func (u *User) Load() {

	api.Get("me")

	u.Name = "John Doe"
	u.Username = "johndoe"
	u.Icon = "https://example.com/icon.png"
	u.Email = "joe@doe.com"
}
