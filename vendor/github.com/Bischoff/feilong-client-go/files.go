/**
  Copyright Contributors to the Feilong Project.

  SPDX-License-Identifier: Apache-2.0
**/

package feilong

import (
	"encoding/json"
)

// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#import-file

type ImportFileOutput struct {
	DestURL		string		`json:"dest_url"`
	filesizeInBytes	int		`json:"filesize_in_bytes"`
	MD5Sum		string		`json:"md5sum"`
}

type ImportFileResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
	Output		ImportFileOutput `json:"output"`
}

func (c *Client) ImportFile(file []byte) (*ImportFileResult, error) {
	var result ImportFileResult

	body, err := c.doRequest("PUT", "/files", file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#export-file

type ExportFileParams struct {
	SourceFile	string		`json:"source_file"`
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
