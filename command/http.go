package command

import (
	"bytes"
	"net/http"
)

var client = &http.Client{}

func getReq(url string, token string) (*http.Response, error) {
	req, reqErr := http.NewRequest("GET", url, bytes.NewBuffer([]byte("")))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "*/*")

	if reqErr != nil {
		return nil, reqErr
	}
	return client.Do(req)
}

func postJSON(url string, jsonStr string, token string) (*http.Response, error) {
	req, reqErr := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	if reqErr != nil {
		return nil, reqErr
	}
	return client.Do(req)
}

func putJSON(url string, jsonStr string, token string) (*http.Response, error) {
	req, reqErr := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	if reqErr != nil {
		return nil, reqErr
	}
	return client.Do(req)
}
