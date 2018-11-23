package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type StatusResponse struct {
	Status  string
	Version string
}

func (c *Client) Status() (*StatusResponse, error) {
	resp, err := c.httpGET("api/v1/status")

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, errors.New(fmt.Sprintf("error getting server status: %s", resp.Status))
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
