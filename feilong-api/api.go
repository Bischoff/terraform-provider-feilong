package feilong_api

import (
	"encoding/json"
)

type ZvmCloudConnectorVersion struct {
	RS		int	`json:"rs"`
	OverallRC	int	`json:"overallRC"`
	ModID		int	`json:"modID"`
	RC		int	`json:"rc"`
	ErrMsg		string	`json:"errmsg"`
	Output struct {
		Version		string	`json:"version"`
		ApiVersion	string	`json:"api_version"`
		MaxVersion	string	`json:"max_version"`
		MinVersion	string	`json:"min_version"`
	}
}

func (c *Client) GetZvmCloudConnectorVersion() (*ZvmCloudConnectorVersion, error) {
	var answer ZvmCloudConnectorVersion

	body, err := c.doRequest("GET", "/")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &answer)
	if err != nil {
		return nil, err
	}

	return &answer, nil
}
