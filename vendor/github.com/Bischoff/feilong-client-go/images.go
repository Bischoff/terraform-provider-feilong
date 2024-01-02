/**
  Copyright Contributors to the Feilong Project.

  SPDX-License-Identifier: Apache-2.0
**/

package feilong

import (
	"encoding/json"
)


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#list-images

type ListImagesImage struct {
	ImageName		string		`json:"imagename"`
	ImageOSDistro		string		`json:"imageosdistro"`
	MD5Sum			string		`json:"md5sum"`
	DiskSizeUnits		string		`json:"disk_size_units"`
	ImageSizeInBytes	string		`json:"image_size_in_bytes"`
	Type			string		`json:"type"`
	Comments		string		`json:"comments"`
	LastAccessTime		float64		`json:"last_access_time"`
}

type ListImagesResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
	Output		[]ListImagesImage `json:"output"`
}

func (c *Client) ListImages(imageName *string) (*ListImagesResult, error) {
	var result ListImagesResult

	path := "/images"
	if imageName != nil {
		path = path + "/" + *imageName
	}
	body, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#create-image

type CreateImageMeta struct {
	OSVersion	string		`json:"os_version"`
	MD5Sum		string		`json:"md5sum,omitempty"`
	DiskType	string		`json:"disk_type,omitempty"`
}

type CreateImageParams struct {
	ImageName	string		`json:"image_name"`
	URL		string		`json:"url"`
	ImageMeta	CreateImageMeta		`json:"image_meta"`
	RemoteHost	string		`json:"remote_host,omitempty"`
}

func (c *Client) CreateImage(params *CreateImageParams) error {
	wrapper := createImageWrapper { Image: *params }

	body, err := json.Marshal(&wrapper)
	if err != nil {
		return err
	}

	_, err = c.doRequest("POST", "/images", body)

	return err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#export-image

type ExportImageParams struct {
	DestURL		string		`json:"dest_url"`
	RemoteHost	string		`json:"remote_host,omitempty"`
}

type ExportImageOutput struct {
	ImageName	string		`json:"image_name"`
	ImagePath	string		`json:"image_path"`
	OSVersion	string		`json:"os_version"`
	MD5Sum		string		`json:"md5sum"`
}

type ExportImageResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
	Output		ExportImageOutput `json:"output"`
}

func (c *Client) ExportImage(name string, params *ExportImageParams) (*ExportImageResult, error) {
	wrapper := exportImageWrapper { Location: *params }
	var result ExportImageResult

	body, err := json.Marshal(&wrapper)
	if err != nil {
		return nil, err
	}

	body, err = c.doRequest("PUT", "/images/" + name, body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#get-root-disk-size-of-image

type GetImageRootDiskSizeResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
	Output		string		`json:"output"`
}

func (c *Client) GetImageRootDiskSize(name string) (*GetImageRootDiskSizeResult, error) {
	var result GetImageRootDiskSizeResult

	body, err := c.doRequest("GET", "/images/" + name + "/root_disk_size", nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#delete-image

func (c *Client) DeleteImage(name string) error {
	_, err := c.doRequest("DELETE", "/images/" + name, nil)

	return err
}


// For internal use

type createImageWrapper struct {
	Image		CreateImageParams `json:"image"`
}

type exportImageWrapper struct {
	Location	ExportImageParams `json:"location"`
}
