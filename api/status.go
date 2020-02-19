package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type StatusResponse struct {
	Status  string
	Version string
}

func (c *Client) Status() (*StatusResponse, error) {
	resp, err := c.HTTPRequest(http.MethodGet, "api/v1/status", nil)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("error getting server status: %s", resp.Status)
	}

	jsonData, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var status StatusResponse
	err = json.Unmarshal(jsonData, &status)

	if err != nil {
		return nil, err
	}

	return &status, nil
}
