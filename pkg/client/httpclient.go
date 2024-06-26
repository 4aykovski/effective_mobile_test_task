package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

type HTTPClient struct {
	Scheme   string
	Host     string
	BasePath string
	client   http.Client
}

func NewHTTPClient(host, basePath, scheme string, client http.Client) *HTTPClient {
	return &HTTPClient{
		Scheme:   scheme,
		Host:     host,
		BasePath: basePath,
		client:   client,
	}
}

func (hc *HTTPClient) Do(ctx context.Context, r *http.Request) ([]byte, error) {
	res, err := hc.client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("can't do request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusUnauthorized && res.StatusCode != http.StatusBadRequest {
		return nil, fmt.Errorf("can't do request: %w", fmt.Errorf("%w to %s: %d", ErrWrongStatusCode, res.Request.URL.String(), res.StatusCode))
	} else if res.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("can't do request: %w", fmt.Errorf("%w to %s", Err401StatusCode, res.Request.URL.String()))
	} else if res.StatusCode == http.StatusBadRequest {
		return nil, fmt.Errorf("can't do request: %w", fmt.Errorf("%w to %s", Err400StatusCode, res.Request.URL.String()))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("can't do request: %w", fmt.Errorf("can't read response body: %w", err))
	}

	return body, nil
}

// CreateRequest return http.Request with given parameters. If you don't need some of the parameters then give nil.
func (hc *HTTPClient) CreateRequest(ctx context.Context, httpMethod string, url string, header http.Header, body io.Reader, query url.Values) (*http.Request, error) {
	req, err := http.NewRequest(httpMethod, url, body)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCantCreateRequest, err)
	}

	if header != nil {
		req.Header = header
	}

	if query != nil {
		req.URL.RawQuery = query.Encode()
	}

	return req, nil
}

func (hc *HTTPClient) GetFullUrl() *url.URL {
	return &url.URL{
		Scheme: hc.Scheme,
		Host:   hc.Host,
		Path:   path.Join(hc.BasePath, "/"),
	}
}

func (hc *HTTPClient) GetUlrWithMethods(method string) *url.URL {
	u := hc.GetFullUrl()
	u.Path = path.Join(u.Path, method)
	return u
}
