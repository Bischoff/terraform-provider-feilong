/**
  Copyright Contributors to the Feilong Project.

  SPDX-License-Identifier: Apache-2.0
**/

package feilong

import (
	"encoding/json"
)

// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#report-health-of-smapi

type SMAPIHealthOutput struct {
	TotalSuccess	int		`json:"totalSuccess"`
	TotalFail	int		`json:"totalFail"`
	LastSuccess	string		`json:"lastSuccess"`
	LastFail	string		`json:"lastFail"`
	ContinuousFail	int		`json:"continuousFail"`
	Healthy		bool		`json:"healthy"`
}

// Deprecated call

type SMAPIHealthyResult struct {
	SMAPI		SMAPIHealthOutput `json:"SMAPI"`
}

func (c *Client) SMAPIHealthy() (*SMAPIHealthyResult, error) {
	var result SMAPIHealthyResult

	body, err := c.doRequest("GET", "/smapi-healthy", nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// New call

type SMAPIHealthResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
	Output		SMAPIHealthOutput `json:"output"`
}

func (c *Client) SMAPIHealth() (*SMAPIHealthResult, error) {
	var result SMAPIHealthResult

	body, err := c.doRequest("GET", "/smapi_health", nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
