package api

import (
	"crypto/tls"
	"crypto/x509"
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

// NewHTTPClientWithHeaders returns a *http.Client that attaches a defined set of Headers to all requests.
// If specified, the client will also trust the CA certificate specified in the client configuration.
func (c *Client) newHTTPClientWithHeaders(headers map[string]string) *http.Client {
	if c.config.Insecure {
		// Return HTTP client where we skip certificate verification
		return &http.Client{
			Transport: &headerTransport{
				base: &http.Transport{
					Proxy: http.ProxyFromEnvironment,
					DialContext: (&net.Dialer{
						Timeout:   30 * time.Second,
						KeepAlive: 30 * time.Second,
						DualStack: true,
					}).DialContext,
					ForceAttemptHTTP2:     true,
					MaxIdleConns:          100,
					IdleConnTimeout:       90 * time.Second,
					TLSHandshakeTimeout:   10 * time.Second,
					ExpectContinueTimeout: 1 * time.Second,

					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: c.config.Insecure,
					},
				},
				headers: headers,
			},
		}
	}

	if len(c.config.CACertificatePEM) > 0 {
		// Create a certificate pool and return a HTTP client with the specified specified CA certificate.
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(c.config.CACertificatePEM))
		return &http.Client{
			Transport: &headerTransport{
				base: &http.Transport{
					Proxy: http.ProxyFromEnvironment,
					DialContext: (&net.Dialer{
						Timeout:   30 * time.Second,
						KeepAlive: 30 * time.Second,
						DualStack: true,
					}).DialContext,
					ForceAttemptHTTP2:     true,
					MaxIdleConns:          100,
					IdleConnTimeout:       90 * time.Second,
					TLSHandshakeTimeout:   10 * time.Second,
					ExpectContinueTimeout: 1 * time.Second,

					TLSClientConfig: &tls.Config{
						RootCAs:            caCertPool,
						InsecureSkipVerify: c.config.Insecure,
					},
				},
				headers: headers,
			},
		}
	}

	// Return a regular default HTTP client
	return &http.Client{
		Transport: &headerTransport{
			base:    http.DefaultTransport,
			headers: headers,
		},
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
