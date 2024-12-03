package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type StatusValue string

const (
	// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
	StatusOK StatusValue = "OK"
	// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
	StatusWarn StatusValue = "WARN"
	// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
	StatusDown StatusValue = "DOWN"
)

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type HealthCheck struct {
	Name          string                 `json:"name"`
	Status        StatusValue            `json:"status"`
	StatusMessage string                 `json:"statusMessage"`
	Fields        map[string]interface{} `json:"fields"`
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
type Health struct {
	Status        StatusValue   `json:"status"`
	StatusMessage string        `json:"statusMessage"`
	Uptime        string        `json:"uptime"`
	Version       string        `json:"version"`
	OK            []HealthCheck `json:"oks"`
	Warn          []HealthCheck `json:"warnings"`
	Down          []HealthCheck `json:"down"`
	rawJson       []byte
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (c *Client) HealthString() (string, error) {
	resp, err := c.HTTPRequest(http.MethodGet, "api/v1/health", nil)
	if err != nil {
		return "", err
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (c *Client) Health() (Health, error) {
	resp, err := c.HTTPRequest(http.MethodGet, "api/v1/health-json", nil)
	if err != nil {
		return Health{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return Health{}, fmt.Errorf("server responded with status code %d", resp.StatusCode)
	}

	var rawJson bytes.Buffer

	var health Health
	err = json.NewDecoder(io.TeeReader(resp.Body, &rawJson)).Decode(&health)

	if health.Down == nil {
		health.Down = []HealthCheck{}
	}

	if health.Warn == nil {
		health.Warn = []HealthCheck{}
	}

	if health.OK == nil {
		health.OK = []HealthCheck{}
	}

	health.rawJson = rawJson.Bytes()

	return health, err
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (h *Health) ChecksMap() map[string]HealthCheck {
	m := map[string]HealthCheck{}

	for _, l := range [][]HealthCheck{h.OK, h.Warn, h.Down} {
		for _, c := range l {
			m[c.Name] = c
		}
	}

	return m
}

// Deprecated: Should no longer be used. https://github.com/CrowdStrike/logscale-go-api-client-example
func (h *Health) Json() []byte {
	return h.rawJson
}
