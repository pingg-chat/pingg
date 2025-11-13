package models

type Channel struct {
	ID          int64
	ShortID     string
	Name        string
	Icon        string
	Description string
	IsPrivate   bool
	IsDM        bool
}
