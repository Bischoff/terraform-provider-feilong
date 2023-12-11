## Virtual Switch Resources

The z/VM Virtual switches allow to connect a virtual network directly to an external physical LAN, and then to act as routers. Furthermore, a vswitch can participate in a IEEE 802.1Q VLAN.

You do not need to define such resources if, for example, the `DEVNET` vswitch is enough for your needs.

Here is an example of `feilong_vswitch` resource:

```terraform
resource "feilong_vswitch" "switch" {
  name             = "my vswitch"

  // optional parameters:
  vswitch          = "MYSWITCH"
  real_device      = "0906"
  controller       = "*"
  connection_type  = "CONNECT"
  network_type     = "ETHERNET"
  router           = "NONROUTER"
  vlan_id          = 2100
  port_type        = "ACCESS"
  gvrp             = "NOGVRP"
  queue_mem        = 8
  native_vlan_id   = 1
  persist          = false
}
```

You can then reference such vswitches from your Feilong guests declarations.


### VSwitch Sections

The `feilong_vswitch` resource sections are optional. They may be used to define the following options: 

 * `name` (mandatory): any arbitrary name to identify this resource. Please try to make it unique.
 * `vswitch` (optional): the desired name of the virtual switch on the z/OS side, maximum 8 characters, all capital letters. If omitted, it will be derived from the `name`.
 * `real_device` (optional): the real device to connect to.
 * `controller` (optional): the controller to use, or `*` for any.
 * `connection_type` (optional): `CONNECT`, `DISCONNECT`, or `NOUPLINK`.
 * `network_type` (optional): `IP` or `ETHERNET`.
 * `router` (optional): `NONROUTER` or `PRIROUTER`.
 * `vlan_id` (optional): VLAN identifier (1 to 4094).
 * `port_type` (optional): `ACCESS` or `TRUNK`.
 * `gvrp` (optional): `GVRP` or `NOGVRP`.
 * `queue_mem` (optional): 1 to 8 (megabytes).
 * `native_vlan_id` (optional): native VLAN identifier (1 to 4094).
 * `persist` (optional): whether the switch is permanent.
