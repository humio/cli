package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type QueryJobs struct {
	client *Client
}

func (c *Client) QueryJobs() *QueryJobs { return &QueryJobs{client: c} }

type Query struct {
	QueryString                string            `json:"queryString"`
	Start                      string            `json:"start,omitempty"`
	End                        string            `json:"end,omitempty"`
	Live                       bool              `json:"isLive,omitempty"`
	TimezoneOffset             *int              `json:"timeZoneOffsetMinutes,omitempty"`
	Arguments                  map[string]string `json:"arguments,omitempty"`
	ShowQueryEventDistribution bool              `json:"showQueryEventDistribution,omitempty"`
}

type QueryResultMetadata struct {
	EventCount       uint64                 `json:"eventCount"`
	ExtraData        map[string]interface{} `json:"extraData"`
	FieldOrder       []string               `json:"fieldOrder"`
	IsAggregate      bool                   `json:"isAggregate"`
	PollAfter        int                    `json:"pollAfter"`
	ProcessedBytes   uint64                 `json:"processedBytes"`
	ProcessedEvents  uint64                 `json:"processedEvents"`
	QueryStart       uint64                 `json:"queryStart"`
	QueryEnd         uint64                 `json:"queryEnd"`
	ResultBufferSize uint64                 `json:"resultBufferSize"`
	TimeMillis       uint64                 `json:"timeMillis"`
	TotalWork        uint64                 `json:"totalWork"`
	WorkDone         uint64                 `json:"workDone"`
}

type QueryResult struct {
	Cancelled bool                     `json:"cancelled"`
	Done      bool                     `json:"done"`
	Events    []map[string]interface{} `json:"events"`
	Metadata  QueryResultMetadata      `json:"metaData"`
}

type QueryError struct {
	error string
}

func (e QueryError) Error() string {
	return e.error
}

func (q QueryJobs) Create(repository string, query Query) (string, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(query)

	if err != nil {
		return "", err
	}

	resp, err := q.client.HTTPRequest(http.MethodPost, "api/v1/repositories/"+url.QueryEscape(repository)+"/queryjobs", &buf)

	if err != nil {
		return "", err
	}

	if resp == nil {
		return "", fmt.Errorf("failed to get response")
	}

	switch resp.StatusCode {
	case http.StatusBadRequest:
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return "", QueryError{string(body)}
	case http.StatusOK:
	default:
		return "", fmt.Errorf("could not create query job, got status code %d", resp.StatusCode)
	}

	var jsonResponse struct {
		ID string `json:"id"`
	}

	err = json.NewDecoder(resp.Body).Decode(&jsonResponse)

	if err != nil {
		return "", err
	}

	return jsonResponse.ID, nil
}

func (q *QueryJobs) PollContext(ctx context.Context, repository string, id string) (QueryResult, error) {
	resp, err := q.client.HTTPRequestContext(ctx, http.MethodGet, "api/v1/repositories/"+url.QueryEscape(repository)+"/queryjobs/"+id, nil, JSONContentType)

	if err != nil {
		return QueryResult{}, err
	}

	if resp == nil {
		return QueryResult{}, fmt.Errorf("failed to get response")
	}

	if resp.StatusCode != http.StatusOK {
		return QueryResult{}, fmt.Errorf("error polling query job, got status code %d", resp.StatusCode)
	}

	var result QueryResult

	err = json.NewDecoder(resp.Body).Decode(&result)

	return result, err
}

func (q *QueryJobs) Delete(repository string, id string) error {
	_, err := q.client.HTTPRequest(http.MethodDelete, "api/v1/repositories/"+url.QueryEscape(repository)+"/queryjobs/"+id, nil)
	return err
}
