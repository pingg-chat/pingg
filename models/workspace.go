package models

import (
	"github.com/pingg-chat/pingg/api"
)

type Workspace struct {
	ID          int64
	Name        string
	Description string
	Icon        string
}

func (w *Workspace) Load() error {
	workspace, err := api.Get[Workspace]("workspaces")

	if err != nil {
		return err
	}

	*w = *workspace
	return nil
}
