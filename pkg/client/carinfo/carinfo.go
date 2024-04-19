package carinfo

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/4aykovski/effective_mobile_test_task/internal/model"
	"github.com/4aykovski/effective_mobile_test_task/pkg/client"
)

const (
	getCarInfoByRegNumberMethod = "/info"
)

type Client struct {
	httpClient *client.HTTPClient
}

func NewClient(httpClient *client.HTTPClient) *Client {
	return &Client{
		httpClient: httpClient,
	}
}

type CarInfo struct {
	RegNumber string      `json:"regNumber"`
	Mark      string      `json:"mark"`
	Model     string      `json:"model"`
	Year      int         `json:"year,omitempty"`
	Owner     model.Owner `json:"owner"`
}

func (c *Client) GetCarInfoByRegNumber(ctx context.Context, regNumber string) ([]byte, error) {
	u := c.httpClient.GetUlrWithMethods(getCarInfoByRegNumberMethod)
	q := url.Values{}
	q.Add("regNum", regNumber)

	req, err := c.httpClient.CreateRequest(ctx, http.MethodGet, u.String(), nil, nil, q)
	if err != nil {
		return nil, fmt.Errorf("can't get car info: %w", err)
	}

	res, err := c.httpClient.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("can't get car info: %w", err)
	}

	return res, nil
}
