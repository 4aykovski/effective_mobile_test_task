package carinfo

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

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
	RegNumber string `json:"regNum"`
	Mark      string `json:"mark"`
	Model     string `json:"model"`
	Year      int    `json:"year,omitempty"`
	Owner     Owner  `json:"owner"`
}

type Owner struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
}

func (c *Client) GetCarInfoByRegNumber(ctx context.Context, regNumber string) ([]byte, error) {
	u := c.httpClient.GetUlrWithMethods(getCarInfoByRegNumberMethod)
	q := url.Values{}
	q.Add("regNum", regNumber)

	req, err := c.httpClient.CreateRequest(ctx, http.MethodGet, u.String(), nil, nil, q)
	if err != nil {
		return nil, fmt.Errorf("can't get car info: %w", err)
	}

	fmt.Println(req.URL.String())

	res, err := c.httpClient.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("can't get car info: %w", err)
	}

	return res, nil
}
