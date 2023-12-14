/**
  Copyright Contributors to the Feilong Project.

  SPDX-License-Identifier: Apache-2.0
**/

package feilong

import (
	"encoding/json"
	"strings"
)


// Common structures

type GuestDisk struct {
	Size		string	`json:"size"`
	Format		string	`json:"format,omitempty"`
	IsBootDisk	bool	`json:"is_boot_disk,omitempty"`
	VDev 		string	`json:"vdev,omitempty"`
	DiskPool	string	`json:"disk_pool,omitempty"`
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#list-guests

type ListGuestsResult struct {
	OverallRC	int	`json:"overallRC"`
	ReturnCode	int	`json:"rc"`
	Reason		int	`json:"rs"`
	ErrorMsg	string	`json:"errmsg"`
	ModuleId	int	`json:"modID"`
	Output		[]string `json:"output"`
}

func (c *Client) ListGuests() (*ListGuestsResult, error) {
	var result ListGuestsResult

	body, err := c.doRequest("GET", "/guests", nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#create-guest

type CreateDiskLoadDev struct {
	PortName	string	`json:"portname,omitempty"`
	LUN		string	`json:"lun,omitempty"`
}

type CreateGuestParams struct {
	UserId		string	`json:"userid"`
	VCPUs		int	`json:"vcpus"`
	Memory		int	`json:"memory"`
	UserProfile	string	`json:"user_profile,omitempty"`
	DiskList	[]GuestDisk `json:"disk_list,omitempty"`
	MaxCPU		int	`json:"max_cpu,omitempty"`
	MaxMem		string	`json:"max_mem,omitempty"`
	IPLFrom		string	`json:"ipl_from,omitempty"`
	IPLParam	string	`json:"ipl_param,omitempty"`
	IPLLoadParam	string	`json:"ipl_loadparam,omitempty"`
	DedicateVDevs	[]string `json:"dedicate_vdevs,omitempty"`
	LoadDev		CreateDiskLoadDev `json:"loaddev,omitempty"`
	Account		string	`json:"account,omitempty"`
	CommentList	[]string `json:"comment_list,omitempty"`
}

type CreateGuestResult struct {
	OverallRC	int	`json:"overallRC"`
	ReturnCode	int	`json:"rc"`
	Reason		int	`json:"rs"`
	ErrorMsg	string	`json:"errmsg"`
	ModuleId	int	`json:"modID"`
	Output		[]GuestDisk `json:"output"`
}

func (c *Client) CreateGuest(params *CreateGuestParams) (*CreateGuestResult, error) {
	var result CreateGuestResult

	wrapper := createGuestWrapper { Guest: *params }
	body, err := json.Marshal(&wrapper)
	if err != nil {
		return nil, err
	}

	body, err = c.doRequest("POST", "/guests", body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#guest-add-disks

type GuestAddDisksParams struct {
	DiskList	[]GuestDisk `json:"disk_list,omitempty"`
}

type GuestAddDisksResult struct {
	OverallRC	int	`json:"overallRC"`
	ReturnCode	int	`json:"rc"`
	Reason		int	`json:"rs"`
	ErrorMsg	string	`json:"errmsg"`
	ModuleId	int	`json:"modID"`
	Output		[]GuestDisk `json:"output"`
}

func (c *Client) GuestAddDisks(userid string, params *GuestAddDisksParams) (*GuestAddDisksResult, error) {
	wrapper := guestAddDisksWrapper { DiskInfo: *params }
	var result GuestAddDisksResult

	body, err := json.Marshal(&wrapper)
	if err != nil {
		return nil, err
	}

	body, err = c.doRequest("POST", "/guests/" + userid + "/disks", body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#guest-delete-disks

type GuestDeleteDisksParams struct {
	VDevList	[]string `json:"vdev_list,omitempty"`
}

func (c *Client) GuestDeleteDisks(userid string, params *GuestDeleteDisksParams) (error) {
	wrapper := guestDeleteDisksWrapper { VDevInfo: *params }

	body, err := json.Marshal(&wrapper)
	if err != nil {
		return err
	}

	_, err = c.doRequest("DELETE", "/guests/" + userid + "/disks", body)

	return err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#show-guest-definition

type ShowGuestDefinitionOutput struct {
	UserDirect	[]string `json:"user_direct,omitempty"`
	CheckInfo	[]string `json:"check_info,omitempty"`
}

type ShowGuestDefinitionResult struct {
	OverallRC	int	`json:"overallRC"`
	ReturnCode	int	`json:"rc"`
	Reason		int	`json:"rs"`
	ErrorMsg	string	`json:"errmsg"`
	ModuleId	int	`json:"modID"`
	Output		ShowGuestDefinitionOutput `json:"output"`
}

func (c *Client) ShowGuestDefinition(userid string) (*ShowGuestDefinitionResult, error) {
	var result ShowGuestDefinitionResult

	body, err := c.doRequest("GET", "/guests/" + userid, nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#delete-guest

func (c *Client) DeleteGuest(userid string) (error) {
	_, err := c.doRequest("DELETE", "/guests/" + userid, nil)

	return err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#get-guest-info

type GetGuestInfoOutput struct {
	MaxMemKB	int	`json:"max_mem_kb"`
	NumCPUs		int	`json:"num_cpu"`
	CPUTimeMuSec	int	`json:"cpu_time_us"`
	PowerState	string	`json:"power_state"`
	MemKB		int	`json:"mem_kb"`
	OnlineCPUNum	int	`json:"online_cpu_num"`
	OSDistro	string	`json:"os_distro"`
	KernelInfo	string	`json:"kernel_info"`
}

type GetGuestInfoResult struct {
	OverallRC	int	`json:"overallRC"`
	ReturnCode	int	`json:"rc"`
	Reason		int	`json:"rs"`
	ErrorMsg	string	`json:"errmsg"`
	ModuleId	int	`json:"modID"`
	Output		GetGuestInfoOutput `json:"output"`
}

func (c *Client) GetGuestInfo(userid string) (*GetGuestInfoResult, error) {
	var result GetGuestInfoResult

	body, err := c.doRequest("GET", "/guests/" + userid + "/info", nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#get-guest-adapters-info

type GetGuestAdaptersInfoAdapter struct {
	LANOwner	string	`json:"lan_owner"`
	LANName		string	`json:"lan_name"`
	AdapterAddress	string	`json:"adapter_address,omitempty"`
	AdapterStatus	string	`json:"adapter_status"`
	MACAddress	string	`json:"mac_address,omitempty"`
	IPAddress	string	`json:"mac_ip_address,omitempty"`
	IPVersion	string	`json:"mac_ip_version"`
}

type GetGuestAdaptersInfoOutput struct {
	Adapters	[]GetGuestAdaptersInfoAdapter `json:"adapters,omitempty"`
}

type GetGuestAdaptersInfoResult struct {
	OverallRC	int	`json:"overallRC"`
	ReturnCode	int	`json:"rc"`
	Reason		int	`json:"rs"`
	ErrorMsg	string	`json:"errmsg"`
	ModuleId	int	`json:"modID"`
	Output		GetGuestAdaptersInfoOutput `json:"output"`
}

func (c *Client) GetGuestAdaptersInfo(userid string) (*GetGuestAdaptersInfoResult, error) {
	var result GetGuestAdaptersInfoResult

	body, err := c.doRequest("GET", "/guests/" + userid + "/adapters", nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#create-guest-nic

type CreateGuestNICParams struct {
	VDev		string		`json:"vdev,omitempty"`
	NICId		string		`json:"nic_id,omitempty"`
	MACAddress	string		`json:"mac_addr,omitempty"`
	Active		bool		`json:"active,omitempty"`
}

func (c *Client) CreateGuestNIC(userid string, params *CreateGuestNICParams) (error) {
	wrapper := createGuestNICWrapper { NIC: *params }
	body, err := json.Marshal(&wrapper)
	if err != nil {
		return err
	}

	_, err = c.doRequest("POST", "/guests/" + userid + "/nic", body)

	return err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#start-guest

func (c *Client) StartGuest(userid string) (error) {
	params := simpleAction { Action: "start" }

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	_, err = c.doRequest("POST", "/guests/" + userid + "/action", body)

	return err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#stop-guest

func (c *Client) StopGuest(userid string) (error) {
	params := simpleAction { Action: "stop" }

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	_, err = c.doRequest("POST", "/guests/" + userid + "/action", body)

	return err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#deploy-guest

type DeployGuestParams struct {
	Action		string	`json:"action"`
	Image		string	`json:"image"`
	TransportFiles	string	`json:"transportfiles,omitempty"` // list of comma-separated files
	RemoteHost	string	`json:"remotehost,omitempty"`
	VDev		string	`json:"vdev,omitempty"`
	Hostname	string	`json:"hostname,omitempty"`
	SkipDiskCopy	bool	`json:"skipdiskcopy,omitempty"`
}

func (c *Client) DeployGuest(userid string, params *DeployGuestParams) (error) {
	params.Action = "deploy"

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	// HACK
	// golang JSON encoder does not allow duplicate keys
	// but Feilong accept multiple "transportfiles" keys
	if strings.Contains(params.TransportFiles, ",") {
		from := `"transportfiles":"` + params.TransportFiles + `"`
		to := ""
		transports := strings.Split(params.TransportFiles, ",")
		last := len(transports) - 1
		for i, s := range(transports) {
			to += `"transportfiles":"` + s + `"`
			if i != last { to += "," }
		}
		body = []byte(strings.Replace(string(body), from, to, 1))
	}

	_, err = c.doRequest("POST", "/guests/" + userid + "/action", body)

	return err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#update-guest-nic

type UpdateGuestNICParams struct {
	Couple		bool	`json:"couple"`
	Active		bool	`json:"active,omitempty"`
	VSwitch		string	`json:"vswitch,omitempty"`
}

func (c *Client) UpdateGuestNIC(userid string, vdev string, params *UpdateGuestNICParams) (error) {
	wrapper := updateGuestNICWrapper { Info: *params }

	body, err := json.Marshal(wrapper)
	if err != nil {
		return err
	}

	_, err = c.doRequest("PUT", "/guests/" + userid + "/nic/" + vdev, body)

	return err
}


// For internal use

type simpleAction struct {
	Action		string	`json:"action"`
}

type createGuestWrapper struct {
	Guest		CreateGuestParams `json:"guest"`
}

type guestAddDisksWrapper struct {
	DiskInfo	GuestAddDisksParams `json:"disk_info"`
}

type guestDeleteDisksWrapper struct {
	VDevInfo	GuestDeleteDisksParams `json:"vdev_info"`
}

type createGuestNICWrapper struct {
	NIC		CreateGuestNICParams `json:"nic"`
}

type updateGuestNICWrapper struct {
	Info		UpdateGuestNICParams `json:"info"`
}
