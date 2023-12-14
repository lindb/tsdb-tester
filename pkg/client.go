package pkg

import (
	"github.com/go-resty/resty/v2"
)

type Client struct {
	cli *resty.Client
}

func NewClient(endpoint string) *Client {
	cli := resty.New()
	cli.SetBaseURL(endpoint)
	return &Client{
		cli: cli,
	}
}

func (cli *Client) Put(params map[string]any) ([]byte, error) {
	resp, err := cli.cli.R().
		SetBody(params).
		SetHeader("Accept", "application/json").
		Put("")
	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}
