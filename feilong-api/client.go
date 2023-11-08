package feilong_api

import (
	"fmt"
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

const defaultHost string = "localhost:35000"
const timeout time.Duration = 10 * time.Second

type Client struct {
	Host		string
	HTTPClient	*http.Client
}

func NewClient(connector *string) (*Client, error) {
	c := Client{
		HTTPClient:	&http.Client{Timeout: timeout},
		Host:		defaultHost,
	}

	if connector != nil {
		c.Host = *connector
	}

	return &c, nil
}

// For internal use
func (c *Client) doRequest(method string, path string, params []byte) ([]byte, error) {
	url := "http://" + c.Host + path
	reader := bytes.NewReader(params)
	req, err := http.NewRequest(method, url, reader)
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
		return nil, fmt.Errorf("HTTP status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
