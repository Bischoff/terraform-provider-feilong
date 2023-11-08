package feilong_api

import (
	"encoding/json"
)

// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#export-file

type ExportFileParams struct {
	SourceFile	string	`json:"source_file"`
}

type ExportFileResult struct {
	Contents	[]byte
}

func (c *Client) ExportFile(params *ExportFileParams) (*ExportFileResult, error) {
	var result ExportFileResult

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	result.Contents, err = c.doRequest("POST", "/files", body)

	return &result, err
}
