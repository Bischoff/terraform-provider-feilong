## Terraform Provider for Feilong

Here is a complete configuration example:

```terraform
terraform {
  required_version = ">= 0.13.4"
  required_providers {
    feilong = {
      source = "bischoff/feilong"
      version = "0.0.2"
    }
  }
}

provider "feilong" {
  connector  = "feilong.example.org"
  local_user = "johndoe@client.example.org"
}

resource "feilong_cloudinit_params" "cloudinit" {
  name       = "opensuse_cloudinit"
  hostname   = "zvm.example.org"
  public_key = "ssh-rsa AAAAB3Nz(...)L5yvQjrN johndoe@client.example.org"
}

resource "feilong_network_params" "network" {
  name       = "opensuse_network"
  os_distro  = "sles15"
}

resource "feilong_guest" "opensuse" {
  name       = "leap"
  memory     = "2G"
  disk       = "20G"
  image      = "opensuse155"

  // optional parameters:
  userid     = "LINUX097"
  vcpus      = 2
  mac        = "12:34:56:78:9a:bc"
  cloudinit_params = feilong_cloudinit_params.cloudinit.file
  network_params   = feilong_network_params.network.file
}
```

Feilong provider section (mandatory):

 * `connector` (optional): the IP address or domain name of the z/VM connector, i.e. the VM where Feilong runs. If you omit it, the value will be taken from environment variable `$ZVM_CONNECTOR`.
 * `local_user` (optional): user name and IP address or domain name of the workstation where you run terraform. You need to specify it if you intend to use cloud-init parameters and/or network parameters. In that case, you must drop the public SSH key of the z/VM connector into file `.shh/authorized_keys` in the home directory of that user. This will allow Feilong to upload the cloud-init parameters file and/or the network parameters file.


Cloud-init parameters sections (optional):

 * `name` (mandatory): any arbitrary name to identify this resource. Please try to make it unique.
 * `hostname` (mandatory): the desired fully qualified domain name of the z/VM guest.
 * `public_key` (mandatory): the desired public key of the default user on the z/VM guest. The name of the default user depends on the cloud-init settings at the time of the preparation of the image. For SUSE distributions, it is usually `sles`.


Network parameters sections (optional):

 * `name` (mandatory): any arbitrary name to identify this resource. Please try to make it unique.
 * `os_distro` (mandatory): the OS and the distribution used to select network parameters like udev definitions, interface network definitions, network service to start, etc. The only value currently defined is "sles15".


Guest sections (optional):

 * `name` (mandatory): any arbitrary name to identify this resource. Please try to make it unique.
 * `memory` (mandatory): desired memory size, as an integer number followed by B, K, M, G
 * `disk` (mandatory): desired disk size, as an integer number followed by B, K, M, G
 * `image` (mandatory): the imaged used to create the guest. This image has to be prepared as explained in Feilong documentation.
 * `userid` (optional): the desired name of the guest on the z/VM side, maximum 8 characters, all capital letters. If omitted, it will be derived from the `name`.
 * `vcpus` (optional): the desired number of virtual CPUs on the guest. If omitted, it will be set to `1`.
 * `mac` (optional): the desired MAC address of the guest, as 6 hexadecimal digits separed by colons. Only last 3 bytes will be used, the first 3 will be ignored by Feilong. Feilong will set these first 3 bytes arbitrarily.
 * `cloudinit_params` (optional): the path to a local file containing an ISO 9660 image containing cloud-init parameters in the format used by openstack. You can:
    * prepare this file yourself, taking your inspiration from the contents of the `profider/files/cfgdrive/` directory in this project, or
    * use a cloud-init parameters section to prepare it automatically. If you do so, use `feilong_cloudinit_params.<CLOUDINIT_RESOURCE_NAME>.file` instead of a hardcoded path.
   In both cases, you must declare the user and hostname of your local machine in `local_user` field of the provider, and accept Feilong's public SSH key.
 * `network_params` (optional): the path to a local file containing a uncompressed tarball containing network parameters in "doscript" format. You can:
    * prepare this file yourself, taking your inspiration from the contents of the `provider/files/network.config/` directory in this project, or
    * use a network parameters section to prepare it automatically. If you do so, use `feilong_network_params.<NETWORK_RESOURCE_NAME>.file` instead of a hardcoded path.
   In both cases, you must declare the user and hostname of your local machine in `local_user` field of the provider, and accept Feilong's public SSH key.
