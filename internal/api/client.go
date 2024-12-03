package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Khan/genqlient/graphql"
	"github.com/humio/cli/internal/api/humiographql"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

const defaultUserAgent = "Humio-go-client/unknown"

type Client struct {
	config        Config
	httpTransport *http.Transport
}

type Response struct {
	Data       interface{}            `json:"data"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
	Errors     ErrorList              `json:"errors,omitempty"`
}

type ErrorList []*GraphqlError

type GraphqlError struct {
	Err        error                  `json:"-"`
	Message    string                 `json:"message"`
	Path       ast.Path               `json:"path,omitempty"`
	Locations  []gqlerror.Location    `json:"locations,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
	Rule       string                 `json:"-"`
	State      map[string]string      `json:"state,omitempty"`
}

func (err *GraphqlError) Error() string {
	var res bytes.Buffer
	if err == nil {
		return ""
	}
	filename, _ := err.Extensions["file"].(string)
	if filename == "" {
		filename = "input"
	}

	res.WriteString(filename)

	if len(err.Locations) > 0 {
		res.WriteByte(':')
		res.WriteString(strconv.Itoa(err.Locations[0].Line))
	}

	res.WriteString(": ")
	if ps := err.pathString(); ps != "" {
		res.WriteString(ps)
		res.WriteByte(' ')
	}

	for key, value := range err.State {
		res.WriteString(fmt.Sprintf("(%s: %s) ", key, value))
	}

	res.WriteString(err.Message)

	return res.String()
}
func (err *GraphqlError) pathString() string {
	return err.Path.String()
}

func (errs ErrorList) Error() string {
	var buf bytes.Buffer
	for _, err := range errs {
		buf.WriteString(err.Error())
		buf.WriteByte('\n')
	}
	return buf.String()
}

func (c *Client) MakeRequest(ctx context.Context, req *graphql.Request, resp *graphql.Response) error {
	var httpReq *http.Request
	var err error

	body, err := json.Marshal(req)
	if err != nil {
		return err
	}
	graphqlURL, err := c.Address().Parse("graphql")
	if err != nil {
		return nil
	}
	httpReq, err = http.NewRequest(
		http.MethodPost,
		graphqlURL.String(),
		bytes.NewReader(body))
	if err != nil {
		return err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	if ctx != nil {
		httpReq = httpReq.WithContext(ctx)
	}
	httpClient := c.newHTTPClientWithHeaders(c.headers())
	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		return err
	}
	if httpResp == nil {
		return fmt.Errorf("could not execute http request")
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		var respBody []byte
		respBody, err = io.ReadAll(httpResp.Body)
		if err != nil {
			respBody = []byte(fmt.Sprintf("<unreadable: %v>", err))
		}
		return fmt.Errorf("returned error %v: %s", httpResp.Status, respBody)
	}

	var actualResponse Response
	actualResponse.Data = resp.Data

	err = json.NewDecoder(httpResp.Body).Decode(&actualResponse)
	resp.Extensions = actualResponse.Extensions
	for _, actualError := range actualResponse.Errors {
		gqlError := gqlerror.Error{
			Err:        actualError.Err,
			Message:    actualError.Message,
			Path:       actualError.Path,
			Locations:  actualError.Locations,
			Extensions: actualError.Extensions,
			Rule:       actualError.Rule,
		}
		resp.Errors = append(resp.Errors, &gqlError)
	}
	if err != nil {
		return err
	}

	// This prints all extentions. To use this properly, use a logger
	//if len(actualResponse.Extensions) > 0 {
	//	for _, extension := range resp.Extensions {
	//		fmt.Printf("%v\n", extension)
	//	}
	//}
	if len(actualResponse.Errors) > 0 {
		return actualResponse.Errors
	}
	return nil
}

type Config struct {
	Address           *url.URL
	UserAgent         string
	Token             string
	CACertificatePEM  string
	Insecure          bool
	ProxyOrganization string
	DialContext       func(ctx context.Context, network, addr string) (net.Conn, error)
}

func DefaultConfig() Config {
	config := Config{}

	return config
}

func (c *Client) Address() *url.URL {
	return c.config.Address
}

func (c *Client) Token() string {
	return c.config.Token
}

func (c *Client) CACertificate() string {
	return c.config.CACertificatePEM
}

func (c *Client) Insecure() bool {
	return c.config.Insecure
}

func (c *Client) Config() Config {
	return c.config
}

func NewClient(config Config) *Client {
	httpTransport := NewHttpTransport(config)
	return NewClientWithTransport(config, httpTransport)
}

func NewClientWithTransport(config Config, httpTransport *http.Transport) *Client {
	if config.Address != nil && !strings.HasSuffix(config.Address.Path, "/") {
		config.Address.Path = config.Address.Path + "/"
	}

	if config.UserAgent == "" {
		config.UserAgent = defaultUserAgent
	}

	return &Client{
		config:        config,
		httpTransport: httpTransport,
	}
}

func (c *Client) headers() map[string]string {
	headers := map[string]string{}

	if c.Token() != "" {
		headers["Authorization"] = fmt.Sprintf("Bearer %s", c.Token())
	}

	if c.config.ProxyOrganization != "" {
		headers["ProxyOrganization"] = c.config.ProxyOrganization
	}

	if c.config.UserAgent != "" {
		headers["User-Agent"] = c.config.UserAgent
	}

	return headers
}

// JSONContentType is "application/json"
const JSONContentType string = "application/json"
const ZIPContentType string = "application/zip"

func (c *Client) HTTPRequest(httpMethod string, path string, body io.Reader) (*http.Response, error) {
	return c.HTTPRequestContext(context.Background(), httpMethod, path, body, JSONContentType)
}

func (c *Client) HTTPRequestContext(ctx context.Context, httpMethod string, path string, body io.Reader, contentType string) (*http.Response, error) {
	if body == nil {
		body = bytes.NewReader(nil)
	}

	url, err := c.Address().Parse(path)
	if err != nil {
		return nil, err
	}

	req, reqErr := http.NewRequestWithContext(ctx, httpMethod, url.String(), body)
	if reqErr != nil {
		return nil, reqErr
	}

	headers := c.headers()
	headers["Content-Type"] = contentType

	var client = c.newHTTPClientWithHeaders(headers)
	return client.Do(req)
}

func queryOwnershipToQueryOwnershipType(o humiographql.SharedQueryOwnershipType) humiographql.QueryOwnershipType {
	switch (o).(type) {
	case *humiographql.SharedQueryOwnershipTypeUserOwnership:
		return humiographql.QueryOwnershipTypeUser
	case *humiographql.SharedQueryOwnershipTypeOrganizationOwnership:
		return humiographql.QueryOwnershipTypeOrganization
	default:
		panic("unknown ownership type")
	}
}
