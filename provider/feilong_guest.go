/**
  Copyright Contributors to the Feilong Project.

  SPDX-License-Identifier: Apache-2.0
**/

package provider

import (
	"context"
	"strings"
	"strconv"
	"errors"
	"fmt"
	"time"

	// There is no replacement for the old resource.StateChangeConf
	// (see https://discuss.hashicorp.com/t/terraform-plugin-framework-what-is-the-replacement-for-waitforstate-or-retrycontext/45538)
        oldresource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Bischoff/feilong-client-go"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &FeilongGuest{}
var _ resource.ResourceWithImportState = &FeilongGuest{}

func NewFeilongGuest() resource.Resource {
	return &FeilongGuest{}
}

// FeilongGuest defines the resource implementation.
type FeilongGuest struct {
	Client *feilong.Client
	LocalUser string
}

// FeilongGuestModel describes the resource data model.
type FeilongGuestModel struct {
	Name		types.String	`tfsdk:"name"`
	UserId		types.String	`tfsdk:"userid"`
	VCPUs		types.Int64	`tfsdk:"vcpus"`
	Memory		types.String	`tfsdk:"memory"`
	Disk		types.String	`tfsdk:"disk"`
	Image		types.String	`tfsdk:"image"`
	MAC		types.String	`tfsdk:"mac"`
	VSwitch		types.String	`tfsdk:"vswitch"`
	CloudinitParams	types.String	`tfsdk:"cloudinit_params"`
	MACAddress	types.String	`tfsdk:"mac_address"`
	IPAddress	types.String	`tfsdk:"ip_address"`
}

func (guest *FeilongGuest) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_guest"
}

func (guest *FeilongGuest) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema {
		MarkdownDescription: "Feilong guest VM resource",

		Attributes: map[string]schema.Attribute {
			"name": schema.StringAttribute {
				MarkdownDescription:	"Arbitrary name for the resource",
				Required:		true,
			},
			"userid": schema.StringAttribute {
				MarkdownDescription:	"System name for z/VM",
				Optional:		true,
				Computed:		true,
			},
			"vcpus": schema.Int64Attribute {
				MarkdownDescription:	"Virtual CPUs count",
				Optional:		true,
				Computed:		true,
				Default:		int64default.StaticInt64(1),
			},
			"memory": schema.StringAttribute {
				MarkdownDescription:	"Memory size with unit (G, M, K, B)",
				Optional:		true,
				Computed:		true,
				Default:		stringdefault.StaticString("512M"),
			},
			"disk": schema.StringAttribute {
				MarkdownDescription:	"Disk size of first disk with unit (T, G, M, K, B)",
				Optional:		true,
				Computed:		true,
				Default:		stringdefault.StaticString("10G"),
				// we could use the size of the image file instead
			},
			"image": schema.StringAttribute {
				MarkdownDescription:	"Image name",
				Required:		true,
			},
			"mac": schema.StringAttribute {
				MarkdownDescription:	"Desired MAC address of first interface",
				Optional:		true,
				Computed:		true,
				Default:		stringdefault.StaticString(""),
			},
			"vswitch": schema.StringAttribute {
				MarkdownDescription:	"Name of virtual switch to connect to",
				Optional:		true,
				Computed:		true,
				Default:		stringdefault.StaticString("DEVNET"),
			},
			"cloudinit_params": schema.StringAttribute {
				MarkdownDescription:	"Path to cloud-init parameters file",
				Optional:		true,
			},
			"mac_address": schema.StringAttribute {
				MarkdownDescription:	"MAC address of first interface after deployment",
				Computed:		true,
			},
			"ip_address": schema.StringAttribute {
				MarkdownDescription:	"IP address of first interface after deployment",
				Computed:		true,
			},
		},
	}
}

func (guest *FeilongGuest) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	guest.Client = &req.ProviderData.(*apiClient).Client
	guest.LocalUser = req.ProviderData.(*apiClient).LocalUser
}

func (guest *FeilongGuest) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FeilongGuestModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Compute computed fields
	resourceName := data.Name.ValueString()
	userid := data.UserId.ValueString()
	if userid == "" {
		userid = strings.ToUpper(resourceName)
		if (len(userid) > 8) {
			userid = userid[:8]
		}
		data.UserId = types.StringValue(userid)
	}

	// Compute values passed to Feilong API but not part of the data model
	size := data.Disk.ValueString()
	vcpus := int(data.VCPUs.ValueInt64())
	memory, err := convertToMegabytes(data.Memory.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Conversion Error", fmt.Sprintf("Got error: %s", err))
		return
	}
	image := data.Image.ValueString()
	mac := data.MAC.ValueString()
	vswitch := data.VSwitch.ValueString()
	cloudinitParams := data.CloudinitParams.ValueString()
	localUser := guest.LocalUser

	// Create the guest
	client := guest.Client
	diskList := []feilong.GuestDisk {
		{
			Size:		size,
			IsBootDisk:	true,
		},
	}
	createParams := feilong.CreateGuestParams {
		UserId:		userid,
		VCPUs:		vcpus,
		Memory:		memory,
		DiskList:	diskList,
	}
	_, err = client.CreateGuest(&createParams)
	if err != nil {
		resp.Diagnostics.AddError("Creation Error", fmt.Sprintf("Got error: %s", err))
		return
	}

	// Create the first network interface
	createNICParams := feilong.CreateGuestNICParams {
		MACAddress:	mac,
	}
	err = client.CreateGuestNIC(userid, &createNICParams)
	if err != nil {
		resp.Diagnostics.AddError("NIC Creation Error", fmt.Sprintf("Got error: %s", err))
		return
	}

	// Couple the first network interface to the virtual switch
	updateNICParams := feilong.UpdateGuestNICParams {
		Couple:		true,
		VSwitch:	vswitch,
	}
	err = client.UpdateGuestNIC(userid, "1000", &updateNICParams)
	if err != nil {
		resp.Diagnostics.AddError("NIC Coupling Error", fmt.Sprintf("Got error: %s", err))
		return
	}

	// Deploy the guest
	deployParams := feilong.DeployGuestParams {
		Image:		image,
		TransportFiles:	cloudinitParams,
		RemoteHost:	localUser,
	}
	err = client.DeployGuest(userid, &deployParams)
	if err != nil {
		resp.Diagnostics.AddError("Deployment Error", fmt.Sprintf("Got error: %s", err))
		return
	}

	// Start the guest
	err = client.StartGuest(userid)
	if err != nil {
		resp.Diagnostics.AddError("Startup Error", fmt.Sprintf("Got error: %s", err))
		return
	}

	// Wait until the guest gets an IP address
	var macAddress string
	var ipAddress string
	err = waitForLease(ctx, client, userid, &macAddress, &ipAddress)
	if err != nil {
		resp.Diagnostics.AddError("Error Waiting for an IP Address", fmt.Sprintf("Got error: %s", err))
		return
	}
	data.MACAddress = types.StringValue(macAddress)
	data.IPAddress = types.StringValue(ipAddress)

	// Write logs using the tflog package
	tflog.Trace(ctx, "created a Feilong guest resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (guest *FeilongGuest) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FeilongGuestModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := guest.Client

	userid := data.UserId.ValueString()

	// Obtain info about this guest
	guestInfo, err := client.GetGuestInfo(userid)
	if err != nil {
		resp.Diagnostics.AddError("Guest Querying Error", fmt.Sprintf("Got error: %s", err))
		return
	}

	// Read number of vCPUs
	data.VCPUs = types.Int64Value(int64(guestInfo.Output.NumCPUs))

	// Read memory
	declaredMemory, err := convertToMegabytes(data.Memory.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Conversion Error", fmt.Sprintf("Got error: %s", err))
		return
	}
	obtainedMemory := guestInfo.Output.MaxMemKB / 1_024
	if declaredMemory == obtainedMemory {
		// do not overwrite memory if value equal but a different unit
		tflog.Info(ctx, "Not replacing memory " +  data.Memory.ValueString() + " with equal value " + strconv.Itoa(obtainedMemory) + "M")
	} else {
		data.Memory = types.StringValue(strconv.Itoa(obtainedMemory) + "M")
	}

	// Obtain first minidisk info
	minidisksInfo, err := client.GetGuestMinidisksInfo(userid)
	if err != nil {
		resp.Diagnostics.AddError("Minidisks Querying Error", fmt.Sprintf("Got error: %s", err))
		return
	}
	if len(minidisksInfo.Output.Minidisks) < 1 {
		resp.Diagnostics.AddError("Minidisk Not Found Error", fmt.Sprintf("Got number: %d", len(minidisksInfo.Output.Minidisks)))
		return
	}
	firstMinidisk := minidisksInfo.Output.Minidisks[0]

	// Read disk size
	declaredDiskSize, err := convertToMegabytes(data.Disk.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Conversion Error", fmt.Sprintf("Got error: %s", err))
		return
	}
	if firstMinidisk.DeviceUnits != "Cylinders" {
		resp.Diagnostics.AddError("Unknown Minidisk Unit Error", fmt.Sprintf("Got unit: %s", firstMinidisk.DeviceUnits))
		return
	}
        // tracks/cylinder=15  blocks/track=12  kilobytes/block=4  15*12*4=720
	obtainedDiskSize := (firstMinidisk.DeviceSize * 720) / 1_024
	if declaredDiskSize == obtainedDiskSize {
		// do not overwrite disk size if value equal but a different unit
		tflog.Info(ctx, "Not replacing disk " + data.Disk.ValueString() + " with equal value " + strconv.Itoa(obtainedDiskSize) + "M")
	} else {
		data.Disk = types.StringValue(strconv.Itoa(obtainedDiskSize) + "M")
	}

	// Obtain first network adapter info
	adaptersInfo, err := client.GetGuestAdaptersInfo(userid)
	if err != nil {
		resp.Diagnostics.AddError("Network Adapter Info Querying Error", fmt.Sprintf("Got error: %s", err))
		return
	}
	if len(adaptersInfo.Output.Adapters) < 1 {
		resp.Diagnostics.AddError("Network Adapter Not Found Error", fmt.Sprintf("Got number: %d", len(adaptersInfo.Output.Adapters)))
		return
	}
	firstAdapter := adaptersInfo.Output.Adapters[0]

	// Read virtual switch name
	data.VSwitch = types.StringValue(firstAdapter.LANName)

	// Read MAC address
	declaredMAC := data.MAC.ValueString()
	obtainedMAC := firstAdapter.MACAddress
	if strings.ToLower(declaredMAC[8:]) == strings.ToLower(obtainedMAC[8:]) {
		// do not overwrite a MAC address if last 3 hex bytes are the same
		tflog.Info(ctx, "Not replacing MAC address " + declaredMAC + " with other MAC address with same last 3 hex bytes " + obtainedMAC)
	} else {
		data.MAC = types.StringValue(obtainedMAC)
	}

	// Read IP address
	declaredIPAddress := data.IPAddress.ValueString()
	obtainedIPAddress := firstAdapter.IPAddress
	if firstAdapter.IPVersion == "6" && strings.HasPrefix(obtainedIPAddress, "fe80:") {
		// do not overwrite an IPv4 address with a link-local IPv6 address
		tflog.Info(ctx, "Not replacing IP address " + declaredIPAddress + " with an IPv6 link-local address " + obtainedIPAddress)
	} else {
		data.IPAddress = types.StringValue(obtainedIPAddress)
	}

	// CAVEATS:
	//  - the image used during the deployment cannot be determined after the deployment
	//  - the cloud init image used during the deployment cannot be determined after the deployment

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (guest *FeilongGuest) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data FeilongGuestModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//	return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (guest *FeilongGuest) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data FeilongGuestModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := guest.Client
	userid := data.UserId.ValueString()

	err := client.DeleteGuest(userid)
	if err != nil {
		resp.Diagnostics.AddError("Deletion Error", fmt.Sprintf("Got error: %s", err))
		return
	}
}

func (guest *FeilongGuest) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// For internal use

func convertToMegabytes(sizeWithUnit string) (int, error) {
	lastButOne := len(sizeWithUnit) - 1

	size, err := strconv.Atoi(sizeWithUnit[:lastButOne])
	if (err != nil) {
		return 0, err
	}

	unit := sizeWithUnit[lastButOne:]

	switch unit {
		case "B":
			size = size / 1_048_576
		case "K":
			size = size / 1_024
		case "M":
			// (nothing to do)
		case "G":
			size = size * 1_024
		case "T":
			size = size * 1_048_576
		default:
			return 0, errors.New("Unit must be one of B K M G T")
	}
	return size, nil
}

const waitingMsg string = "Still waiting for IP address"
const obtainedMsg string = "IP address obtained"

func waitForLease(ctx context.Context, client *feilong.Client, userid string, macAddress *string, ipAddress *string) error {
	waitFunction := func() (interface{}, string, error) {
		err := getAddresses(client, userid, macAddress, ipAddress)
		if err != nil {
			return false, "", err
		}
		if *ipAddress == "" {
			return false, waitingMsg, nil
		}
		return true, obtainedMsg, nil
	}

	stateConf := &oldresource.StateChangeConf {
		Pending:	[]string { waitingMsg },
		Target:		[]string { obtainedMsg },
		Refresh:	waitFunction,
		Timeout:	1 * time.Minute,
		MinTimeout:	3 * time.Second,
		Delay:		5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func getAddresses(client *feilong.Client, userid string, macAddress *string, ipAddress *string) error {
	result, err := client.GetGuestAdaptersInfo(userid)
	if err != nil {
		return err
	}

	if len(result.Output.Adapters) == 0 {
		return errors.New("Adapters Query Error")
	}
	*macAddress = result.Output.Adapters[0].MACAddress
	*ipAddress = result.Output.Adapters[0].IPAddress
	return nil
}
