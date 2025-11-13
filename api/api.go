package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pingg-chat/pingg/config"
	"github.com/pingg-chat/pingg/utils"
)

func Get[T any](path string) (*T, error) {
	return makeRequest[T](path, "GET")
}

func makeRequest[T any](path string, method string) (*T, error) {
	client := &http.Client{}

	apiUrl := "http://127.0.0.1:8000/" + path

	req, err := http.NewRequest(method, apiUrl, nil)

	if err != nil {
		// tratar erro
		utils.Dd("Error on New Request", err)
	}

	req.Header.Set("X-Auth-User-Id", config.GetID())
	req.Header.Set("X-API-KEY", "base64:ASGQzvXwjj5EPlbpfakJ58k7y8ZkLGSct2MfHuSkUw0=")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		utils.Dd("Error on Client Do", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.Dd("Error on ReadAll", err)
	}

	if resp.StatusCode == 401 {
		utils.Dd("Error 401 Unauthorized")
	}

	var result T

	err = json.Unmarshal(body, &result)

	if err != nil {
		utils.Dd("Error on Unmarshal", err)
	}

	return &result, nil
}
