package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/pingg-chat/pingg/config"
	"github.com/pingg-chat/pingg/utils"
)

func Get(path string) {
	makeRequest(path, "GET")
}

func makeRequest(path string, method string) {
	client := &http.Client{}

	apiUrl := "http://127.0.0.1:8000/" + path

	req, err := http.NewRequest(method, apiUrl, nil)

	if err != nil {
		// tratar erro
		utils.Dd("Error on New Request", err)
		return
	}

	id := config.GetID()

	req.Header.Set("X-Auth-User-Id", id)
	req.Header.Set("X-API-KEY", "base64:ASGQzvXwjj5EPlbpfakJ58k7y8ZkLGSct2MfHuSkUw0=")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		// tratar erro
		fmt.Println("Error on Client Do", err)
		utils.Dd("Error on Client Do", err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.Dd("Error on ReadAll", err)
		return
	}

	utils.Dd(string(body))
}
