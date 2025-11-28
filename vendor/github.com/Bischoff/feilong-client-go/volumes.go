/**
  Copyright Contributors to the Feilong Project.

  SPDX-License-Identifier: Apache-2.0
**/

package feilong

import (
	"strings"
	"encoding/json"
)


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#refresh-volume-bootmap-info

type RefreshVolumeBootmapInfoParams struct {
	FCPChannel	[]string	`json:"fcpchannel"`
	WWPN		[]string	`json:"wwpn"`
	LUN		string		`json:"lun"`
	TransportFiles	string		`json:"transportfiles,omitempty"`
	GuestNetworks	[]GuestNetwork	`json:"guest_networks,omitempty"`
}

func (c *Client) RefreshVolumeBootmapInfo(params *RefreshVolumeBootmapInfoParams) error {
	wrapper := refreshVolumeBootmapInfoWrapper { Info: *params }
	body, err := json.Marshal(&wrapper)
	if err != nil {
		return err
	}
	_, err = c.doRequest("PUT", "/volumes/volume_refresh_bootmap", body)

	return err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#get-volume-connector

type GetVolumeConnectorParams struct {
	Reserve		bool		`json:"reserve"`
	FCPTemplateId	string		`json:"fcp_template_id,omitempty"`
	StorageProvider	string		`json:"storage_provider,omitempty"`
	PCHIdInfo	string		`json:"pchid_info,omitempy"`
}

type GetVolumeConnectorOutput struct {
	FCP		[]string	`json:"zvm_fcp"`
	WWPNs		[]string	`json:"wwpns"`
	Host		string		`json:"host"`
}

type GetVolumeConnectorResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
	Output		GetVolumeConnectorOutput `json:"output"`
}

func (c *Client) GetVolumeConnector(userid string, params *GetVolumeConnectorParams) (*GetVolumeConnectorResult, error) {
	var result GetVolumeConnectorResult

	wrapper := getVolumeConnectorWrapper { Info: *params }
	body, err := json.Marshal(&wrapper)
	if err != nil {
		return nil, err
	}

	body, err = c.doRequest("GET", "/volumes/conn/" + userid, body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#create-fcp-template

type CreateFCPTemplateParams struct {
	Name		string		`json:"name"`
	Description	string		`json:"description,omitempty"`
	FCPDevices	string		`json:"fcp_devices,omitempty"`
	HostDefault	bool		`json:"host_default,omitempty"`
	StorageProviders []string	`json:"storage_providers,omitempty"`
	MinFCPPathsCount int		`json:"min_fcp_paths_count,omitempty"`
}

type CreateFCPTemplate struct {
	Id		string		`json:"id"`
	Name		string		`json:"name"`
	Description	string		`json:"description,omitempty"`
	HostDefault	bool		`json:"host_default,omitempty"`
	StorageProviders []string	`json:"storage_providers,omitempty"`
	MinFCPPathsCount int		`json:"min_fcp_paths_count,omitempty"`
}

type CreateFCPTemplateOutput struct {
	FCPTemplate	CreateFCPTemplate `json:"fcp_template"`
}

type CreateFCPTemplateResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
	Output		CreateFCPTemplateOutput `json:"output"`
}

func (c *Client) CreateFCPTemplate(params *CreateFCPTemplateParams) (*CreateFCPTemplateResult, error) {
	var result CreateFCPTemplateResult

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	body, err = c.doRequest("POST", "/volumes/fcptemplates", body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#delete-fcp-template

func (c *Client) DeleteFCPTemplate(templateId string) error {
	_, err := c.doRequest("DELETE", "/volumes/fcptemplates/" + templateId, nil)

	return err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#get-fcp-templates

type GetFCPTemplate struct {
	Name		string		`json:"name"`
	Id		string		`json:"id"`
	Description	string		`json:"description"`
	HostDefault	bool		`json:"host_default"`
	StorageProviderDefault []string	`json:"sp_default"`
	MinFCPPathsCount int		`json:"min_fcp_paths_count"`
	CPCSerialNumber	string		`json:"cpc_sn"`
	CPCName		string		`json:"cpc_name"`
	LogicalPartition string		`json:"lpar"`
	HypervisorHostname string	`json:"hypervisor_hostname"`
	PhysicalChannelIds []string	`json:"pchids"`
}

type GetFCPTemplatesOutput struct {
	FCPTemplates	[]GetFCPTemplate `json:"fcp_templates"`
}

type GetFCPTemplatesResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
	Output		GetFCPTemplatesOutput `json:"output"`
}

func (c *Client) GetFCPTemplates(templateIdList []string, assignerId *string, defaultSPList *string, hostDefault *string) (*GetFCPTemplatesResult, error) {
	var result GetFCPTemplatesResult

	req := "/volumes/fcptemplates"
	sep := "?"
	if templateIdList != nil {
		req = req + sep + "template_id_list=['" + strings.Join(templateIdList, "','") + "']"
		sep = "&"
	}
	if assignerId != nil {
		req = req + sep + "assigner_id=" + *assignerId
		sep = "&"
	}
	if defaultSPList != nil {
		req = req + sep + "default_sp_list=" + *defaultSPList
		sep = "&"
	}
	if hostDefault != nil {
		req = req + sep + "host_default=" + *hostDefault
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


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#edit-fcp-template

type EditFCPTemplateParams struct {
	Name		*string		`json:"name,omitempty"`
	Description	*string		`json:"description,omitempty"`
	FCPDevices	*string		`json:"fcp_devices,omitempty"`
	HostDefault	*bool		`json:"host_default,omitempty"`
	DefaultSPList	[]string	`json:"default_sp_list,omitempty"`
	MinFCPPathsCount *int		`json:"min_fcp_paths_count,omitempty"`
}

func (c *Client) EditFCPTemplate(templateId string, params *EditFCPTemplateParams) error {
	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	_, err = c.doRequest("PUT", "/volumes/fcptemplates/" + templateId, body)

	return err
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#get-fcp-templates-details

type FCPTemplateStatistics struct {
	Total		string		`json:"total"`
	TotalCount	map[string]int	`json:"total_count"`
	SingleFCP	string		`json:"single_fcp"`
	RangeFCP	string		`json:"range_fcp"`
	Available	string		`json:"available"`
	AvailableCount	map[string]int	`json:"available_count"`
	Allocated	string		`json:"allocated"`
	ReserveOnly	string		`json:"reserve_only"`
	ConnectionOnly	string		`json:"connection_only"`
	UnallocatedButActive map[string]string `json:"unallocated_but_active"`
	AllocatedButFree string		`json:"allocated_but_free"`
	NotFound	string		`json:"notfound"`
	Offline		string		`json:"offline"`
	ChannelIds	map[string]string `json:"chids"`
	PhysicalChannelIds map[string]string `json:"pchids"`
}

type GetFCPTemplateDetails struct {
	Id		string		`json:"id"`
	Name		string		`json:"name"`
	Description	string		`json:"description"`
	HostDefault	bool		`json:"host_default"`
	StorageProviders []string	`json:"storage_providers"`
	MinFCPPathsCount int		`json:"min_fcp_paths_count"`
	Raw		map[string][][]interface{} `json:"raw,omitempty"`
	Statistics	map[string]FCPTemplateStatistics `json:"statistics,omitempty"`
	PhysicalChannelIds[]string	`json:"pchids"`
	CPCSerialNumber	string		`json:"cpc_sn"`
	CPCName		string		`json:"cpc_name"`
	LogicalPartition string		`json:"lpar"`
	HypervisorHostname string	`json:"hypervisor_hostname"`
}

type GetFCPTemplatesDetailsOutput struct {
	FCPTemplates	[]GetFCPTemplateDetails `json:"fcp_templates"`
}

type GetFCPTemplatesDetailsResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
	Output		GetFCPTemplatesDetailsOutput `json:"output"`
}

func (c *Client) GetFCPTemplatesDetails(templateIdList []string, raw bool, statistics bool, syncWithZVM bool) (*GetFCPTemplatesDetailsResult, error) {
	var result GetFCPTemplatesDetailsResult

	req := "/volumes/fcptemplates/detail?"
	if templateIdList != nil {
		req += "template_id_list=['" + strings.Join(templateIdList, "','") + "']&"
	}
	if raw {
		req += "raw=true&"
	} else {
		req += "raw=false&"
	}
	if statistics {
		req += "statistics=true&"
	} else {
		req += "statistics=false&"
	}
	if syncWithZVM {
		req += "sync_with_zvm=true"
	} else {
		req += "sync_with_zvm=false"
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


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#get-fcp-usage

type GetVolumeFCPUsageResult struct {
	OverallRC	int		`json:"overallRC"`
	ReturnCode	int		`json:"rc"`
	Reason		int		`json:"rs"`
	ErrorMsg	string		`json:"errmsg"`
	ModuleId	int		`json:"modID"`
	Output		[]interface{}	`json:"output"`
}

func (c *Client) GetVolumeFCPUsage(fcpid string) (*GetVolumeFCPUsageResult, error) {
	var result GetVolumeFCPUsageResult

	body, err := c.doRequest("GET", "/volumes/fcp/" + fcpid, nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}


// https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#set-fcp-usage

type SetVolumeFCPUsageParams struct {
	UserId		string		`json:"userid"`
	Reserved	int		`json:"reserved"`
	Connections	int		`json:"connections"`
	FCPTemplateId	string		`json:"fcp_template_id"`
}

func (c *Client) SetVolumeFCPUsage(fcpid string, params *SetVolumeFCPUsageParams) error {
	wrapper := setVolumeFCPUsageWrapper { Info: *params }
	body, err := json.Marshal(&wrapper)
	if err != nil {
		return err
	}
	_, err = c.doRequest("PUT", "/volumes/fcp/" + fcpid, body)

	return err
}


// For internal use

type refreshVolumeBootmapInfoWrapper struct {
	Info		RefreshVolumeBootmapInfoParams `json:"info"`
}

type getVolumeConnectorWrapper struct {
	Info		GetVolumeConnectorParams `json:"info"`
}

type setVolumeFCPUsageWrapper struct {
	Info		SetVolumeFCPUsageParams `json:"info"`
}
