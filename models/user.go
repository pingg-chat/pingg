package models

import (
	"github.com/pingg-chat/pingg/api"
)

type User struct {
	ID        int64
	Name      string
	Username  string
	Icon      string
	Email     string
	Workspace *Workspace
}

func (u *User) Load() error {
	workspace := Workspace{}
	u.Workspace = &workspace
	user, err := api.Get[User]("me")

	if err != nil {
		return err
	}

	*u = *user
	return nil
}
