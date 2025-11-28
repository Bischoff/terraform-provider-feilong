/**
  Copyright Contributors to the Feilong Project.

  SPDX-License-Identifier: Apache-2.0
**/

package feilong

import (
	"encoding/json"
)


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#list-guests

type ListGuestsResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
	Output		[]string	`json:"output"`
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
	PortName	string		`json:"portname,omitempty"`
	LUN		string		`json:"lun,omitempty"`
}

type CreateGuestParams struct {
	UserId		string		`json:"userid"`
	VCPUs		int		`json:"vcpus"`
	Memory		int		`json:"memory"`
	UserProfile	string		`json:"user_profile,omitempty"`
	DiskList	[]GuestDisk	`json:"disk_list,omitempty"`
	MaxCPU		int		`json:"max_cpu,omitempty"`
	MaxMem		string		`json:"max_mem,omitempty"`
	IPLFrom		string		`json:"ipl_from,omitempty"`
	IPLParam	string		`json:"ipl_param,omitempty"`
	IPLLoadParam	string		`json:"ipl_loadparam,omitempty"`
	DedicateVDevs	[]string	`json:"dedicate_vdevs,omitempty"`
	LoadDev		CreateDiskLoadDev `json:"loaddev,omitempty"`
	Account		string		`json:"account,omitempty"`
	CommentList	[]string	`json:"comment_list,omitempty"`
}

type CreateGuestResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
	Output		[]GuestDisk	`json:"output"`
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


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#get-guest-minidisks-info

type GetGuestMinidisksInfoMinidisk struct {
	VDev		string		`json:"vdev"`
	RDev		string		`json:"rdev"`
	AccessType	string		`json:"access_type"`
	DeviceType	string		`json:"device_type"`
	DeviceSize	int		`json:"device_size"`
	DeviceUnits	string		`json:"device_units"`
	VolumeLabel	string		`json:"volume_label"`
}

type GetGuestMinidisksInfoOutput struct {
	Minidisks	[]GetGuestMinidisksInfoMinidisk `json:"minidisks,omitempty"`
}

type GetGuestMinidisksInfoResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
	Output		GetGuestMinidisksInfoOutput `json:"output"`
}

func (c *Client) GetGuestMinidisksInfo(userid string) (*GetGuestMinidisksInfoResult, error) {
	var result GetGuestMinidisksInfoResult

	body, err := c.doRequest("GET", "/guests/" + userid + "/disks", nil)
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

type AddGuestDisksParams struct {
	DiskList	[]GuestDisk	`json:"disk_list,omitempty"`
}

type AddGuestDisksResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
	Output		[]GuestDisk	`json:"output"`
}

func (c *Client) AddGuestDisks(userid string, params *AddGuestDisksParams) (*AddGuestDisksResult, error) {
	var result AddGuestDisksResult

	wrapper := addGuestDisksWrapper { DiskInfo: *params }
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


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#guest-configure-disks

type ConfigureGuestDisk struct {
	VDev		string		`json:"vdev,omitempty"`
	Format		string		`json:"format"`
	MountDirectory	string		`json:"mntdir,omitempty"`
}

type ConfigureGuestDisksParams struct {
	DiskList	[]ConfigureGuestDisk `json:"disk_list,omitempty"`
}

func (c *Client) ConfigureGuestDisks(userid string, params *ConfigureGuestDisksParams) error {
	wrapper := configureGuestDisksWrapper { DiskInfo: *params }
	body, err := json.Marshal(wrapper)
	if err != nil {
		return err
	}

	_, err = c.doRequest("PUT", "/guests/" + userid + "/disks", body)

	return err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#guest-delete-disks

type DeleteGuestDisksParams struct {
	VDevList	[]string	`json:"vdev_list,omitempty"`
}

func (c *Client) DeleteGuestDisks(userid string, params *DeleteGuestDisksParams) error {
	wrapper := deleteGuestDisksWrapper { VDevInfo: *params }
	body, err := json.Marshal(&wrapper)
	if err != nil {
		return err
	}

	_, err = c.doRequest("DELETE", "/guests/" + userid + "/disks", body)

	return err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#attach-volume

type AttachGuestVolumeParams struct {
	AssignerId	string		`json:"assigner_id"`
	FCPList		[]string	`json:"zvm_fcp"`
	FCPTemplateId	string		`json:"fcp_template_id"`
	TargetWWPN	[]string	`json:"target_wwpn"`
	TargetLUN	string		`json:"target_lun"`
	OSVersion	string		`json:"os_version"`
	Multipath	bool		`json:"multipath"`
	MountPoint	string		`json:"mount_point,omitempty"`
	IsRootVolume	bool		`json:"is_root_volume,omitempty"`
	DoRollback	bool		`json:"do_rollback,omitempty"`
}

func (c *Client) AttachGuestVolume(params *AttachGuestVolumeParams) error {
	wrapper2 := attachGuestVolumeWrapper2 { Connection: *params }
	wrapper := attachGuestVolumeWrapper { Info: wrapper2 }
	body, err := json.Marshal(&wrapper)
	if err != nil {
		return err
	}

	_, err = c.doRequest("POST", "/guests/volumes", body)

	return err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#detach-volume

type DetachGuestVolumeParams struct {
	AssignerId	string		`json:"assigner_id"`
	FCPList		[]string	`json:"zvm_fcp"`
	FCPTemplateId	string		`json:"fcp_template_id"`
	TargetWWPN	[]string	`json:"target_wwpn"`
	TargetLUN	string		`json:"target_lun"`
	OSVersion	string		`json:"os_version"`
	Multipath	bool		`json:"multipath"`
	MountPoint	string		`json:"mount_point,omitempty"`
	IsRootVolume	bool		`json:"is_root_volume,omitempty"`
	UpdateConnectionsOnly bool	`json:"update_connections_only,omitempty"`
	DoRollback	bool		`json:"do_rollback,omitempty"`
}

func (c *Client) DetachGuestVolume(params *DetachGuestVolumeParams) error {
	wrapper2 := detachGuestVolumeWrapper2 { Connection: *params }
	wrapper := detachGuestVolumeWrapper { Info: wrapper2 }
	body, err := json.Marshal(&wrapper)
	if err != nil {
		return err
	}

	_, err = c.doRequest("DELETE", "/guests/volumes", body)

	return err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#get-guests-stats-including-cpu-and-memory

type GetGuestsStatsDetails struct {
	GuestCPUs	int		`json:"guest_cpus"`
	UsedCPUTime	int		`json:"used_cpu_time_us"`
	ElapsedCPUTime	int		`json:"elapsed_cpu_time_us"`
	MinCPUCount	int		`json:"min_cpu_count"`
	MaxCPULimit	int		`json:"max_cpu_limit"`
	SamplesCPUInUse	int		`json:"samples_cpu_in_use"`
	SamplesCPUDelay	int		`json:"samples_cpu_delay"`
	UsedMemoryKB	int		`json:"used_mem_kb"`
	MaxMemoryKB	int		`json:"max_mem_kb"`
	MinMemoryKB	int		`json:"min_mem_kb"`
	SharedMemoryKB	int		`json:"shared_mem_kb"`
	TotalMemory	int		`json:"total_memory"`
	AvailableMemory	int		`json:"available_memory"`
	FreeMemory	int		`json:"free_memory"`
}

type GetGuestsStatsResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
	Output		map[string]GetGuestsStatsDetails `json:"output"`
}

func (c *Client) GetGuestsStats(userid string) (*GetGuestsStatsResult, error) {
	var result GetGuestsStatsResult

	body, err := c.doRequest("GET", "/guests/stats?userid=" + userid, nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#get-guests-interface-stats

type GetGuestsInterfaceStats struct {
	VSwitch		string		`json:"vswitch_name,omitempty"`
	VDev		string		`json:"nic_vdev"`
	FramesRec	int		`json:"nic_fr_rx"`
	FramesSent	int		`json:"nic_fr_tx"`
	FramesRecDisc	int		`json:"nic_fr_rx_dsc"`
	FramesSentDisc	int		`json:"nic_fr_tx_dsc"`
	FramesRecErr	int		`json:"nic_fr_rx_err"`
	FramesSentErr	int		`json:"nic_fr_tx_err"`
	BytesRec	int		`json:"nic_rx"`
	BytesSent	int		`json:"nic_tx"`
}

type GetGuestsInterfaceStatsResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
	Output		map[string][]GetGuestsInterfaceStats `json:"output"`
}

func (c *Client) GetGuestsInterfaceStats(userid string) (*GetGuestsInterfaceStatsResult, error) {
	var result GetGuestsInterfaceStatsResult

	body, err := c.doRequest("GET", "/guests/interfacestats?userid=" + userid, nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#get-guests-nic-info

type GetGuestsNICInfo struct {
	UserId		string		`json:"userid"`
	Interface	string		`json:"interface"`
	VSwitch		string		`json:"switch"`
	Port		string		`json:"port"`
	Comments	string		`json:"comments"`
}

type GetGuestsNICInfoResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
	Output		[]GetGuestsNICInfo `json:"output"`
}

func (c *Client) GetGuestsNICInfo(userid *string, nicid *string, vswitch *string) (*GetGuestsNICInfoResult, error) {
	var result GetGuestsNICInfoResult

	req := "/guests/nics"
	sep := "?"
	if userid != nil {
		req = req + sep + "userid=" + *userid
		sep = "&"
	}
	if nicid != nil {
		req = req + sep + "nic_id=" + *nicid
		sep = "&"
	}
	if vswitch != nil {
		req = req + sep + "vswitch=" + *vswitch
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


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#show-guest-definition

type ShowGuestDefinitionOutput struct {
	UserDirect	[]string	`json:"user_direct,omitempty"`
}

type ShowGuestDefinitionResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
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

func (c *Client) DeleteGuest(userid string) error {
	_, err := c.doRequest("DELETE", "/guests/" + userid, nil)

	return err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#get-guest-power-state-from-hypervisor

type GetGuestPowerStateFromHypervisorResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
	Output		string		`json:"output"`
}

func (c *Client) GetGuestPowerStateFromHypervisor(userid string) (*GetGuestPowerStateFromHypervisorResult, error) {
	var result GetGuestPowerStateFromHypervisorResult

	body, err := c.doRequest("GET", "/guests/" + userid + "/power_state_real", nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#get-guest-info

type GetGuestInfoOutput struct {
	MaxMemKB	int		`json:"max_mem_kb"`
	NumCPUs		int		`json:"num_cpu"`
	CPUTimeMuSec	int		`json:"cpu_time_us"`
	PowerState	string		`json:"power_state"`
	MemKB		int		`json:"mem_kb"`
	OnlineCPUNum	int		`json:"online_cpu_num"`
	OSDistro	string		`json:"os_distro"`
	KernelInfo	string		`json:"kernel_info"`
}

type GetGuestInfoResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
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


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#get-guest-os-info

type GetGuestOSInfoOutput struct {
	OSDistro	string		`json:"os_distro"`
	KernelInfo	string		`json:"kernel_info"`
}

type GetGuestOSInfoResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
	Output		GetGuestOSInfoOutput `json:"output"`
}

func (c *Client)  GetGuestOSInfo(userid string) (*GetGuestOSInfoResult, error) {
	var result GetGuestOSInfoResult

	body, err := c.doRequest("GET", "/guests/" + userid + "/os_info", nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#get-guest-online-cpu-num

type GetGuestOnlineCPUNumResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
	Output		int		`json:"output"`
}

func (c *Client) GetGuestOnlineCPUNum(userid string) (*GetGuestOnlineCPUNumResult, error) {
	var result GetGuestOnlineCPUNumResult

	body, err := c.doRequest("GET", "/guests/" + userid + "/online_cpu_num", nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#get-guest-user-direct

type GetGuestUserDirectoryOutput struct {
	UserDirect	[]string	`json:"user_direct,omitempty"`
}

type GetGuestUserDirectoryResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
	Output		GetGuestUserDirectoryOutput `json:"output"`
}

func (c *Client) GetGuestUserDirectory(userid string) (*GetGuestUserDirectoryResult, error) {
	var result GetGuestUserDirectoryResult

	body, err := c.doRequest("GET", "/guests/" + userid + "/user_direct", nil)
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
	LANOwner	string		`json:"lan_owner"`
	LANName		string		`json:"lan_name"`
	AdapterAddress	string		`json:"adapter_address,omitempty"`
	AdapterStatus	string		`json:"adapter_status"`
	MACAddress	string		`json:"mac_address,omitempty"`
	IPAddress	string		`json:"mac_ip_address,omitempty"`
	IPVersion	string		`json:"mac_ip_version"`
}

type GetGuestAdaptersInfoOutput struct {
	Adapters	[]GetGuestAdaptersInfoAdapter `json:"adapters,omitempty"`
}

type GetGuestAdaptersInfoResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
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

func (c *Client) CreateGuestNIC(userid string, params *CreateGuestNICParams) error {
	wrapper := createGuestNICWrapper { NIC: *params }
	body, err := json.Marshal(&wrapper)
	if err != nil {
		return err
	}

	_, err = c.doRequest("POST", "/guests/" + userid + "/nic", body)

	return err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#create-network-interface

type CreateGuestNetworkInterfaceParams struct {
	OSVersion	string		`json:"os_version"`
	GuestNetworks	[]GuestNetwork	`json:"guest_networks"`
	Active		bool		`json:"active,omitempty"`
}

func (c *Client) CreateGuestNetworkInterface(userid string, params *CreateGuestNetworkInterfaceParams) error {
	wrapper := createGuestNetworkInterfaceWrapper { Interface: *params }
	body, err := json.Marshal(&wrapper)
	if err != nil {
		return err
	}

	_, err = c.doRequest("POST", "/guests/" + userid + "/interface", body)

	return err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#delete-network-interface

type DeleteGuestNetworkInterfaceParams struct {
	OSVersion	string		`json:"os_version"`
	VDev		string		`json:"vdev"`
	Active		bool		`json:"active,omitempty"`
}

func (c *Client) DeleteGuestNetworkInterface(userid string, params *DeleteGuestNetworkInterfaceParams) error {
	wrapper := deleteGuestNetworkInterfaceWrapper { Interface: *params }
	body, err := json.Marshal(&wrapper)
	if err != nil {
		return err
	}

	_, err = c.doRequest("DELETE", "/guests/" + userid + "/interface", body)

	return err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#start-guest

func (c *Client) StartGuest(userid string) error {
	params := simpleAction { Action: "start" }

	return c.doAction(userid, &params)
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#stop-guest

func (c *Client) StopGuest(userid string) error {
	params := simpleAction { Action: "stop" }

	return c.doAction(userid, &params)
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#softstop-guest

func (c *Client) SoftStopGuest(userid string) error {
	params := simpleAction { Action: "softstop" }

	return c.doAction(userid, &params)
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#pause-guest

func (c *Client) PauseGuest(userid string) error {
	params := simpleAction { Action: "pause" }

	return c.doAction(userid, &params)
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#unpause-guest

func (c *Client) UnpauseGuest(userid string) error {
	params := simpleAction { Action: "unpause" }

	return c.doAction(userid, &params)
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#reboot-guest

func (c *Client) RebootGuest(userid string) error {
	params := simpleAction { Action: "reboot" }

	return c.doAction(userid, &params)
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#reset-guest

func (c *Client) ResetGuest(userid string) error {
	params := simpleAction { Action: "reset" }

	return c.doAction(userid, &params)
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#get-guest-console-output

type GetGuestConsoleOutputResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
	Output		[]string	`json:"output"`
}

func (c *Client) GetGuestConsoleOutput(userid string) (*GetGuestConsoleOutputResult, error) {
	var result GetGuestConsoleOutputResult
	params := simpleAction { Action: "get_console_output" }

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	body, err = c.doRequest("POST", "/guests/" + userid + "/action", body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#live-migration-of-guest

type LiveMigrateGuestOptions struct {
	MaxTotal	int		`json:"maxtotal,omitempty"`
	MaxQuiesce	int		`json:"maxquiesce,omitempty"`
	Immediate	string		`json:"immediate,omitempty"`
}

type LiveMigrateGuestParams struct {
	Action		string		`json:"action"`
	DestUserId	string		`json:"dest_zcc_userid,omitempty"`
	Destination	string		`json:"destination"`
	Options		LiveMigrateGuestOptions `json:"parms,omitempty"`
	MigrationAction	string		`json:"lgr_action"`
}

func (c *Client) LiveMigrateGuest(userid string, params *LiveMigrateGuestParams) error {
	params.Action = "live_migrate_vm"

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	_, err = c.doRequest("POST", "/guests/" + userid + "/action", body)

	return err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#guest-register

type RegisterGuestParams struct {
	Action		string		`json:"action"`
	Meta		string		`json:"meta"`
	NetSet		string		`json:"net_set"`
	Port		map[string]string `json:"port"`
}

func (c *Client) RegisterGuest(userid string, params *RegisterGuestParams) error {
	params.Action = "register_vm"

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	_, err = c.doRequest("POST", "/guests/" + userid + "/action", body)

	return err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#guest-deregister

func (c *Client) DeregisterGuest(userid string) error {
	params := simpleAction { Action: "deregister_vm" }

	return c.doAction(userid, &params)
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#live-resize-cpus-of-guest

type LiveResizeGuestCPUsParams struct {
	Action		string		`json:"action"`
	CPUCount	int		`json:"cpu_cnt"`
}

func (c *Client) LiveResizeGuestCPUs(userid string, params *LiveResizeGuestCPUsParams) error {
	params.Action = "live_resize_cpus"

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	_, err = c.doRequest("POST", "/guests/" + userid + "/action", body)

	return err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#resize-cpus-of-guest

type ResizeGuestCPUsParams struct {
	Action		string		`json:"action"`
	CPUCount	int		`json:"cpu_cnt"`
}

func (c *Client) ResizeGuestCPUs(userid string, params *ResizeGuestCPUsParams) error {
	params.Action = "resize_cpus"

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	_, err = c.doRequest("POST", "/guests/" + userid + "/action", body)

	return err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#live-resize-memory-of-guest

type LiveResizeGuestMemoryParams struct {
	Action		string		`json:"action"`
	Size		string		`json:"size"`
}

func (c *Client) LiveResizeGuestMemory(userid string, params *LiveResizeGuestMemoryParams) error {
	params.Action = "live_resize_mem"

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	_, err = c.doRequest("POST", "/guests/" + userid + "/action", body)

	return err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#resize-memory-of-guest

type ResizeGuestMemoryParams struct {
	Action		string		`json:"action"`
	Size		string		`json:"size"`
}

func (c *Client) ResizeGuestMemory(userid string, params *ResizeGuestMemoryParams) error {
	params.Action = "resize_mem"

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	_, err = c.doRequest("POST", "/guests/" + userid + "/action", body)

	return err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#deploy-guest

type DeployGuestParams struct {
	Action		string		`json:"action"`
	Image		string		`json:"image"`
	TransportFiles	string		`json:"transportfiles,omitempty"`
	RemoteHost	string		`json:"remotehost,omitempty"`
	VDev		string		`json:"vdev,omitempty"`
	Hostname	string		`json:"hostname,omitempty"`
	SkipDiskCopy	bool		`json:"skipdiskcopy,omitempty"`
}

func (c *Client) DeployGuest(userid string, params *DeployGuestParams) error {
	params.Action = "deploy"

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	_, err = c.doRequest("POST", "/guests/" + userid + "/action", body)

	return err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#capture-guest

type CaptureGuestParams struct {
	Action		string		`json:"action"`
	Image		string		`json:"image"`
	CaptureType	string		`json:"capture_type,omitempty"`
	CompressLevel	int		`json:"compress_level,omitempty"`
}

func (c *Client) CaptureGuest(userid string, params *CaptureGuestParams) error {
	params.Action = "capture"

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	_, err = c.doRequest("POST", "/guests/" + userid + "/action", body)

	return err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#grow-root-volume-of-guest

type GrowGuestRootVolumeParams struct {
	Action		string		`json:"action"`
	OSVersion	string		`json:"os_version"`
}

func (c *Client) GrowGuestRootVolume(userid string, params *GrowGuestRootVolumeParams) error {
	params.Action = "grow_root_volume"

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	_, err = c.doRequest("POST", "/guests/" + userid + "/action", body)

	return err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#get-guest-power-state

type GetGuestPowerStateResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
	Output		string		`json:"output"`
}

func (c *Client) GetGuestPowerState(userid string) (*GetGuestPowerStateResult, error) {
	var result GetGuestPowerStateResult

	body, err := c.doRequest("GET", "/guests/" + userid + "/power_state", nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#update-guest-nic

type UpdateGuestNICParams struct {
	Couple		bool		`json:"couple"`
	Active		bool		`json:"active,omitempty"`
	VSwitch		string		`json:"vswitch,omitempty"`
}

func (c *Client) UpdateGuestNIC(userid string, vdev string, params *UpdateGuestNICParams) error {
	wrapper := updateGuestNICWrapper { Info: *params }
	body, err := json.Marshal(wrapper)
	if err != nil {
		return err
	}

	_, err = c.doRequest("PUT", "/guests/" + userid + "/nic/" + vdev, body)

	return err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#delete-guest-nic

type DeleteGuestNICParams struct {
	Active		bool		`json:"active,omitempty"`
}

func (c *Client) DeleteGuestNIC(userid string, vdev string, params *DeleteGuestNICParams) error {
	body, err := json.Marshal(&params)
	if err != nil {
		return err
	}

	_, err = c.doRequest("DELETE", "/guests/" + userid + "/nic/" + vdev, body)

	return err
}


// For internal use

type simpleAction struct {
	Action		string		`json:"action"`
}

func (c *Client) doAction(userid string, params *simpleAction) error {
	body, err := json.Marshal(*params)
	if err != nil {
		return err
	}

	_, err = c.doRequest("POST", "/guests/" + userid + "/action", body)

	return err
}

type createGuestWrapper struct {
	Guest		CreateGuestParams `json:"guest"`
}

type addGuestDisksWrapper struct {
	DiskInfo	AddGuestDisksParams `json:"disk_info"`
}

type configureGuestDisksWrapper struct {
	DiskInfo	ConfigureGuestDisksParams `json:"disk_info"`
}

type deleteGuestDisksWrapper struct {
	VDevInfo	DeleteGuestDisksParams `json:"vdev_info"`
}

type attachGuestVolumeWrapper2 struct {
	Connection	AttachGuestVolumeParams `json:"connection"`
}

type attachGuestVolumeWrapper struct {
	Info		attachGuestVolumeWrapper2 `json:"info"`
}

type detachGuestVolumeWrapper2 struct {
	Connection	DetachGuestVolumeParams `json:"connection"`
}

type detachGuestVolumeWrapper struct {
	Info		detachGuestVolumeWrapper2 `json:"info"`
}

type createGuestNICWrapper struct {
	NIC		CreateGuestNICParams `json:"nic"`
}

type createGuestNetworkInterfaceWrapper struct {
	Interface	CreateGuestNetworkInterfaceParams `json:"interface"`
}

type deleteGuestNetworkInterfaceWrapper struct {
	Interface	DeleteGuestNetworkInterfaceParams `json:"interface"`
}

type updateGuestNICWrapper struct {
	Info		UpdateGuestNICParams `json:"info"`
}
