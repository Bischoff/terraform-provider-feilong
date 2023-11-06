package feilong_api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const DefaultHost string = "localhost:35000"

type Client struct {
	Host		string
	HTTPClient	*http.Client
}

func NewClient(connector *string) (*Client, error) {
	c := Client{
		HTTPClient:	&http.Client{Timeout: 10 * time.Second},
		Host:		DefaultHost,
	}

	if connector != nil {
		c.Host = *connector
	}

	return &c, nil
}

// For internal use
func (c *Client) doRequest(method string, path string) ([]byte, error) {
	url := "http://" + c.Host + path
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
