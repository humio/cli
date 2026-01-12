package api

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"time"
)

// We must our own http.Client which adds the authorization header in all requests sent to Humio.
// We use the approach described here: https://github.com/shurcooL/graphql/issues/28#issuecomment-464713908

type headerTransport struct {
	base    http.RoundTripper
	headers map[string]string
}

// newUnixSocketProxyDialer creates a dialer that connects through a unix socket proxy
func newUnixSocketProxyDialer(socketPath string, fallbackDialer func(ctx context.Context, network, addr string) (net.Conn, error)) func(ctx context.Context, network, addr string) (net.Conn, error) {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		// Connect directly to the unix socket
		// The proxy is expected to handle the connection transparently
		conn, err := net.Dial("unix", socketPath)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to unix socket proxy: %w", err)
		}

		return conn, nil
	}
}

func NewHttpTransport(config Config) *http.Transport {
	dialContext := config.DialContext
	if dialContext == nil {
		dialContext = (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext
	}

	// If a unix socket proxy is specified, wrap the dialer to use it
	if config.UnixSocketProxy != "" {
		dialContext = newUnixSocketProxyDialer(config.UnixSocketProxy, dialContext)
	}

	if config.Insecure {
		// Return HTTP transport where we skip certificate verification
		return &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           dialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,

			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: config.Insecure, // #nosec G402
			},
		}
	}

	if len(config.CACertificatePEM) > 0 {
		// Create a certificate pool and return a HTTP transport with the specified specified CA certificate.
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(config.CACertificatePEM))
		return &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           dialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,

			TLSClientConfig: &tls.Config{
				RootCAs:            caCertPool,
				InsecureSkipVerify: config.Insecure, // #nosec G402
			},
		}
	}

	// Return a regular default HTTP client
	return &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

// NewHTTPClientWithHeaders returns a *http.Client that attaches a defined set of Headers to all requests.
func (c *Client) newHTTPClientWithHeaders(headers map[string]string) *http.Client {
	return &http.Client{
		Transport: &headerTransport{
			base:    c.httpTransport,
			headers: headers,
		},
		Timeout: 30 * time.Second,
	}
}

func (h *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := CloneRequest(req)
	for key, val := range h.headers {
		req2.Header.Set(key, val)
	}
	return h.base.RoundTrip(req2)
}

// CloneRequest and CloneHeader copied from https://github.com/kubernetes/apimachinery/blob/a76b7114b20a2e56fd698bba815b1e2c82ec4bff/pkg/util/net/http.go#L469-L491

// CloneRequest creates a shallow copy of the request along with a deep copy of the Headers.
func CloneRequest(req *http.Request) *http.Request {
	r := new(http.Request)

	// shallow clone
	*r = *req

	// deep copy headers
	r.Header = CloneHeader(req.Header)

	return r
}

// CloneHeader creates a deep copy of an http.Header.
func CloneHeader(in http.Header) http.Header {
	out := make(http.Header, len(in))
	for key, values := range in {
		newValues := make([]string, len(values))
		copy(newValues, values)
		out[key] = newValues
	}
	return out
}
