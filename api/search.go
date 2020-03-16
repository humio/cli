package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type QueryJobs struct {
	client *Client
}

func (c *Client) QueryJobs() *QueryJobs { return &QueryJobs{client: c} }

type Query struct {
	QueryString    string            `json:"queryString"`
	Start          string            `json:"start,omitempty"`
	End            string            `json:"end,omitempty"`
	Live           bool              `json:"isLive,omitempty"`
	TimezoneOffset *int              `json:"timeZoneOffsetMinutes,omitempty"`
	Arguments      map[string]string `json:"arguments,omitempty"`
}

type QueryResultMetadata struct {
	EventCount  uint64 `json:"eventCount"`
	IsAggregate bool   `json:"isAggregate"`
	PollAfter   int    `json:"pollAfter"`
	QueryStart  uint64 `json:"queryStart"`
	QueryEnd    uint64 `json:"queryEnd"`
}

type QueryResult struct {
	Cancelled bool                     `json:"cancelled"`
	Done      bool                     `json:"done"`
	Events    []map[string]interface{} `json:"events"`
	Metadata  QueryResultMetadata      `json:"metaData"`
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

	if resp.StatusCode != http.StatusOK {
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

func (q *QueryJobs) Poll(repository string, id string) (QueryResult, error) {
	return q.PollContext(context.Background(), repository, id)
}

func (q *QueryJobs) PollContext(ctx context.Context, repository string, id string) (QueryResult, error) {
	resp, err := q.client.HTTPRequestContext(ctx, http.MethodGet, "api/v1/repositories/"+url.QueryEscape(repository)+"/queryjobs/"+id, bytes.NewBuffer(nil))

	if err != nil {
		return QueryResult{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return QueryResult{}, fmt.Errorf("error polling query job, got status code %d", resp.StatusCode)
	}

	var result QueryResult

	err = json.NewDecoder(resp.Body).Decode(&result)

	return result, err
}

func (q *QueryJobs) Delete(repository string, id string) error {
	_, err := q.client.HTTPRequest(http.MethodDelete, "api/v1/repositories/"+url.QueryEscape(repository)+"/queryjobs/"+id, bytes.NewBuffer(nil))
	return err
}
