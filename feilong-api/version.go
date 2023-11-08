package feilong_api

import (
	"encoding/json"
)

// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#get-feilong-version

type GetFeilongVersionResult struct {
	OverallRC	int	`json:"overallRC"`
	ReturnCode	int	`json:"rc"`
	Reason		int	`json:"rs"`
	ErrorMsg	string	`json:"errmsg"`
	ModuleId	int	`json:"modID"`
	Output struct {
		Version		string	`json:"version"`
		APIVersion	string	`json:"api_version"`
		MaxVersion	string	`json:"max_version"`
		MinVersion	string	`json:"min_version"`
	}
}

func (c *Client) GetFeilongVersion() (*GetFeilongVersionResult, error) {
	var result GetFeilongVersionResult

	body, err := c.doRequest("GET", "/", nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
