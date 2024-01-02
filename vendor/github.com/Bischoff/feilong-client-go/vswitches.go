/**
  Copyright Contributors to the Feilong Project.

  SPDX-License-Identifier: Apache-2.0
**/

package feilong

import (
	"encoding/json"
)


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#create-vswitch

type CreateVSwitchParams struct {
	Name		string		`json:"name"`
	RealDev		string		`json:"rdev,omitempty"`
	Controller	string		`json:"controller,omitempty"`
	Connection	string		`json:"connection,omitempty"`
	NetworkType	string		`json:"network_type,omitempty"`
	Router		string		`json:"router,omitempty"`
	VLANId		interface{}	`json:"vid,omitempty"`
	PortType	string		`json:"port_type,omitempty"`
	GVRP		string		`json:"gvrp,omitempty"`
	QueueMem	int		`json:"queue_mem,omitempty"`
	NativeVLANId	int		`json:"native_vid,omitempty"`
	Persist		bool		`json:"persist,omitempty"`
}

func (c *Client) CreateVSwitch(params *CreateVSwitchParams) error {
	wrapper := createVSwitchWrapper { VSwitch: *params }
	body, err := json.Marshal(&wrapper)
	if err != nil {
		return err
	}

	_, err = c.doRequest("POST", "/vswitches", body)
	return err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#list-vswitches

type ListVSwitchesResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
	Output		[]string	`json:"output"`
}

func (c *Client) ListVSwitches() (*ListVSwitchesResult, error) {
	var result ListVSwitchesResult

	body, err := c.doRequest("GET", "/vswitches", nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#get-vswitch-details

type VSwitchRealDevice struct {
	VDev		string		`json:"vdev"`
	Controller	string		`json:"controller"`
	PortName	string		`json:"port_name"`
	DevStatus	string		`json:"dev_status"`
	DevErr		string		`json:"dev_err"`
}

type VSwitchAuthorizedUser struct {
	PortNum		string		`json:"port_num"`
	PromMode	string		`json:"prom_mode"`
	OSDSim		string		`json:"osd_sim"`
	VLANCount	string		`json:"vlan_count"`
	VLANIds		[]string	`json:"vlan_ids"`
}

type VSwitchAdapter struct {
	MAC		string		`json:"mac"`
	Type		string		`json:"type"`
}

type GetVSwitchDetailsOutput struct {
	SwitchName	string		`json:"switch_name"`
	SwitchType	string		`json:"switch_type"`
	SwitchStatus	string		`json:"switch_status"`
	PortType	string		`json:"port_type"`
	TransportType	string		`json:"transport_type"`
	IsolationStatus	string		`json:"isolation_status"`
	UserPortBased	string		`json:"user_port_based"`
	VLANAwareness	string		`json:"vlan_awareness"`
	VLANId		string		`json:"vlan_id"`
	NativeVLANId	string		`json:"native_vlan_id"`
	QueueMemoryLimit string		`json:"queue_memory_limit"`
	MACAddress	string		`json:"mac_address"`
	MACProtect	string		`json:"MAC_protect"`
	GVRPRequestAttribute string	`json:"gvrp_request_attribute"`
	GVRPEnabledAttribute string	`json:"gvrp_enabled_attribute"`
	LinkAg		string		`json:"link_ag"`
	LAgGroup	string		`json:"lag_group"`
	LAgInterval	string		`json:"lag_interval"`
	IPTimeout	string		`json:"IP_timeout"`
	RoutingValue	string		`json:"routing_value"`
	VLANCounters	string		`json:"VLAN_counters"`
	RealDevices	map[string]VSwitchRealDevice `json:"real_devices"`
	AuthorizedUsers	map[string]VSwitchAuthorizedUser `json:"authorized_users"`
	Adapters	map[string]VSwitchAdapter `json:"adapters"`
}

type GetVSwitchDetailsResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
	Output		GetVSwitchDetailsOutput `json:"output"`
}

func (c *Client) GetVSwitchDetails(name string) (*GetVSwitchDetailsResult, error) {
	var result GetVSwitchDetailsResult

	body, err := c.doRequest("GET", "/vswitches/" + name, nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#grant-user-to-vswitch

type GrantUserToVSwitchParams struct {
	GrantUserId	string		`json:"grant_userid"`
}

func (c *Client) GrantUserToVSwitch(name string, params *GrantUserToVSwitchParams) error {
	wrapper := grantUserToVSwitchWrapper { VSwitch: *params }

	body, err := json.Marshal(&wrapper)
	if err != nil {
		return err
	}

	_, err = c.doRequest("PUT", "/vswitches/" + name, body)

	return err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#revoke-user-from-vswitch

type RevokeUserFromVSwitchParams struct {
	RevokeUserId	string		`json:"revoke_userid"`
}

func (c *Client) RevokeUserFromVSwitch(name string, params *RevokeUserFromVSwitchParams) error {
	wrapper := revokeUserFromVSwitchWrapper { VSwitch: *params }

	body, err := json.Marshal(&wrapper)
	if err != nil {
		return err
	}

	_, err = c.doRequest("PUT", "/vswitches/" + name, body)

	return err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#set-user-vlanid-to-vswitch

type UserVLANId struct {
	UserId		string		`json:"userid"`
	VLANId		int		`json:"vlanid"`
}

type SetUserVLANIdToVSwitchParams struct {
	UserVLANId	UserVLANId	`json:"user_vlan_id"`
}

func (c *Client) SetUserVLANIdToVSwitch(name string, params *SetUserVLANIdToVSwitchParams) error {
	wrapper := setUserVLANIdToVSwitchWrapper { VSwitch: *params }

	body, err := json.Marshal(&wrapper)
	if err != nil {
		return err
	}

	_, err = c.doRequest("PUT", "/vswitches/" + name, body)

	return err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#delete-vswitch

func (c *Client) DeleteVSwitch(name string) error {
	_, err := c.doRequest("DELETE", "/vswitches/" + name, nil)

	return err
}


// For internal use

type createVSwitchWrapper struct {
	VSwitch		CreateVSwitchParams `json:"vswitch"`
}

type grantUserToVSwitchWrapper struct {
	VSwitch		GrantUserToVSwitchParams `json:"vswitch"`
}

type revokeUserFromVSwitchWrapper struct {
	VSwitch		RevokeUserFromVSwitchParams `json:"vswitch"`
}

type setUserVLANIdToVSwitchWrapper struct {
	VSwitch		SetUserVLANIdToVSwitchParams `json:"vswitch"`
}
