/**
  Copyright Contributors to the Feilong Project.

  SPDX-License-Identifier: Apache-2.0
**/

package feilong

import (
	"encoding/json"
)


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#get-guests-list

type GetHostGuestListResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
	Output		[]string	`json:"output"`
}

// English grammar: "guest list", not "guests list"
func (c *Client) GetHostGuestList() (*GetHostGuestListResult, error) {
	var result GetHostGuestListResult

	body, err := c.doRequest("GET", "/host/guests", nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#get-host-info

type GetHostInfoCPUInfo struct {
	CECModel	string		`json:"cec_model"`
	Architecture	string		`json:"architecture"`
}

type GetHostInfoOutput struct {
	ZVMHost		string		`json:"zvm_host"`
	ZCCUserID	string		`json:"zcc_userid"`
	IPLTime		string		`json:"ipl_time"`
	HypervisorHostname string	`json:"hypervisor_hostname"`
	HypervisorType	string		`json:"hypervisor_type"`
	HypervisorVersion int		`json:"hypervisor_version"`
	DiskTotal	int		`json:"disk_total"`
	DiskUsed	int		`json:"disk_used"`
	DiskAvailable	int		`json:"disk_available"`
	VCPUs		int		`json:"vcpus"`
	VCPUsUsed	int		`json:"vcpus_used"`
	MemoryMB	float64		`json:"memory_mb"`
	MemoryMBUsed	float64		`json:"memory_mb_used"`
	CPUInfo		GetHostInfoCPUInfo `json:"cpu_info"`
}

type GetHostInfoResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
	Output		GetHostInfoOutput `json:"output"`
}

func (c *Client) GetHostInfo() (*GetHostInfoResult, error) {
	var result GetHostInfoResult

	body, err := c.doRequest("GET", "/host", nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#get-host-disk-pool-info

type GetHostDiskPoolInfoOutput struct {
	DiskAvailable	int		`json:"disk_available"`
	DiskTotal	int		`json:"disk_total"`
	DiskUsed	int		`json:"disk_used"`
}

type GetHostDiskPoolInfoResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
	Output		GetHostDiskPoolInfoOutput `json:"output"`
}

func (c *Client) GetHostDiskPoolInfo(poolName *string) (*GetHostDiskPoolInfoResult, error) {
	var result GetHostDiskPoolInfoResult

	req := "/host/diskpool"
        if poolName != nil {
		req += "?poolname=" + *poolName
	}

	body, err := c.doRequest("GET", req, nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Same, with details:

type Volume struct {
	VolumeName	string		`json:"volume_name"`
	DeviceType	string		`json:"device_type"`
	StartCylinder	string		`json:"start_cylinder"`
	FreeSize	int		`json:"free_size"`
	DASDGroup	string		`json:"dasd_group"`
	RegionName	string		`json:"region_name"`
}

type GetHostDiskPoolDetailsResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
	Output		map[string][]Volume `json:"output"`
}

func (c *Client) GetHostDiskPoolDetails(poolName *string) (*GetHostDiskPoolDetailsResult, error) {
	var result GetHostDiskPoolDetailsResult

	req := "/host/diskpool"
        if poolName != nil {
		req += "?poolname=" + *poolName + "&details=true"
	} else {
		req += "?details=true"
	}

	body, err := c.doRequest("GET", req, nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#get-host-disk-pool-volume-names

type GetHostDiskPoolVolumeNamesOutput struct {
	DiskPoolVolumes	string		`json:"diskpool_volumes"`
}

type GetHostDiskPoolVolumeNamesResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
	Output		GetHostDiskPoolVolumeNamesOutput `json:"output"`
}

func (c *Client) GetHostDiskPoolVolumeNames(poolName *string) (*GetHostDiskPoolVolumeNamesResult, error) {
	var result GetHostDiskPoolVolumeNamesResult

	req := "/host/diskpool_volumes"
        if poolName != nil {
		req += "?poolname=" + *poolName
	}

	body, err := c.doRequest("GET", req, nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#get-host-volume-info

type GetHostVolumeInfoOutput struct {
	VolumeType	string		`json:"volume_type"`
	VolumeSize	string		`json:"volume_size"`
}

type GetHostVolumeInfoResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
	Output		GetHostVolumeInfoOutput `json:"output"`
}

func (c *Client) GetHostVolumeInfo(volumeName string) (*GetHostVolumeInfoResult, error) {
	var result GetHostVolumeInfoResult

	body, err := c.doRequest("GET", "/host/volume?volumename=" + volumeName, nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#get-host-ssi-cluster-info

type GetHostSSIClusterInfoResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
	Output		[]string	`json:"output"`
}

func (c *Client) GetHostSSIClusterInfo() (*GetHostSSIClusterInfoResult, error) {
	var result GetHostSSIClusterInfoResult

	body, err := c.doRequest("GET", "/host/ssi", nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
