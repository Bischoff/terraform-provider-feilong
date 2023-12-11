## Guest Resources

The guests are the virtual machines provided by z/VM.

Here is an example of `feilong_guest` resource:

```terraform
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
  vswitch          = feilong_vswitch.switch.vswitch
}
```


### Guest Sections

The `feilong_guest` resource sections are optional. They may be used to define the following options:

 * `name` (mandatory): any arbitrary name to identify this resource. Please try to make it unique.
 * `memory` (mandatory): desired memory size, as an integer number followed by B, K, M, or G.
 * `disk` (mandatory): desired disk size, as an integer number followed by B, K, M, or G.
 * `image` (mandatory): the imaged used to create the guest. This image has to be prepared as explained in Feilong documentation.
 * `userid` (optional): the desired name of the guest on the z/VM side, maximum 8 characters, all capital letters. If omitted, it will be derived from the `name`.
 * `vcpus` (optional): the desired number of virtual CPUs on the guest. If omitted, it will be set to `1`.
 * `mac` (optional): the desired MAC address of the guest, as 6 hexadecimal digits separed by colons. Only last 3 bytes will be used, the first 3 will be ignored by Feilong. Feilong will set these first 3 bytes arbitrarily.
 * `cloudinit_params` (optional): the path to a local file containing an ISO 9660 image containing cloud-init parameters in the format used by openstack.
 * `network_params` (optional): the path to a local file containing a uncompressed tarball containing network parameters in "doscript" format.
 * `vswitch` (optional): the name of the virtual switch to connect to. If omitted, it will be set to `DEVNET`.

You can prepare the cloud-init parameters file yourself, taking your inspiration from the contents of the `profider/files/cfgdrive/` directory in this project. Alternatively, you can use a `feilong_cloudinit_params` section to prepare it automatically. If you do so, use `feilong_cloudinit_params.<CLOUDINIT_RESOURCE_NAME>.file` instead of a hardcoded path.
In both cases, you must declare the user and hostname of your local machine in `local_user` field of the provider, and accept Feilong's public SSH key.

You can prepare the network parameters file yourself, taking your inspiration from the contents of the `provider/files/network.config/` directory in this project. Alternatively, you can use a `feilong_network_params` section to prepare it automatically. If you do so, use `feilong_network_params.<NETWORK_RESOURCE_NAME>.file` instead of a hardcoded path.
In both cases, you must declare the user and hostname of your local machine in `local_user` field of the provider, and accept Feilong's public SSH key.

You can use any already existing vswitch, or use a `feilong_vswitch` section to define your own vswitch. If you do so, use `feilong_vswitch.<VSWITCH_RESOURCE_NAME>.vswitch` instead of a hardcoded name.
