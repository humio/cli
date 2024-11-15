package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Masterminds/semver/v3"
)

type StatusResponse struct {
	Status  string
	Version string
}

func (s StatusResponse) IsDown() bool {
	return s.Status != "OK" && s.Status != "WARN"
}

func (s StatusResponse) AtLeast(ver string) (bool, error) {
	assumeLatest := true
	version := strings.Split(s.Version, "-")
	constraint, err := semver.NewConstraint(fmt.Sprintf(">= %s", ver))
	if err != nil || len(version) == 0 {
		return assumeLatest, fmt.Errorf("could not parse constraint of `%s`: %w", fmt.Sprintf(">= %s", ver), err)
	}
	semverVersion, err := semver.NewVersion(version[0])
	if err != nil {
		return assumeLatest, fmt.Errorf("could not parse version of `%s`: %w", version[0], err)
	}

	return constraint.Check(semverVersion), nil
}

func (c *Client) Status() (*StatusResponse, error) {
	resp, err := c.HTTPRequest(http.MethodGet, "api/v1/status", nil)

	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("failed to get response")
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("error getting server status: %s", resp.Status)
	}

	jsonData, err := io.ReadAll(resp.Body)

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
