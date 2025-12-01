/**
  Copyright Contributors to the Feilong Project.

  SPDX-License-Identifier: Apache-2.0
**/

package feilong


// Common structures

type GuestDisk struct {
	Size		string		`json:"size"`
	Format		string		`json:"format,omitempty"`
	IsBootDisk	*bool		`json:"is_boot_disk,omitempty"`
	VDev		string		`json:"vdev,omitempty"`
	DiskPool	string		`json:"disk_pool,omitempty"`
}

type GuestNetwork struct {
	Method		string		`json:"method,omitempty"`
	IPAddress	string		`json:"ip_addr,omitempty"`
	DNSAddresses	[]string	`json:"dns_addr,omitempty"`
	GatewayAddress	string		`json:"gateway_addr,omitempty"`
	CIDR		string		`json:"cidr,omitempty"`
	NICVDev		string		`json:"nic_vdev,omitempty"`
	MACAddress	string		`json:"mac_addr,omitempty"`
	NICId		string		`json:"nic_id,omitempty"`
	OSADevice	string		`json:"osa_device,omitempty"`
	Hostname	string		`json:"hostname,omitempty"`
}
