package config

import (
	"fmt"
)

type User struct {
	ID int64
}

type Api struct {
	Url    string
	ApiKey string
}

type Config struct {
	User User
	Api  Api
}

var config Config

func main() {
	config = Config{
		User: User{
			ID: 12345,
		},
		Api: Api{
			Url:    "http://127.0.0.1:8000/",
			ApiKey: "base64:ASGQzvXwjj5EPlbpfakJ58k7y8ZkLGSct2MfHuSkUw0=",
		},
	}
}

// ---------------------------------------

var userID int64

func SetID(id int64) {
	userID = id
}

func GetID() string {
	return fmt.Sprintf("%d", userID)
}

