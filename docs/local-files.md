## Local File Resources

Here is an example of `feilong_cloudinit_params` resource definition:

```terraform
resource "feilong_cloudinit_params" "cloudinit" {
  name       = "opensuse_cloudinit"
  hostname   = "zvm.example.org"
  public_key = "ssh-rsa AAAAB3Nz(...)L5yvQjrN johndoe@client.example.org"
}
```

It creates a file inside `/tmp/terraform-provider-feilong/` on the local machine. When Feilong needs it, it can upload it.

This upload mechanism implies that you store the SSH public key of the `zvm_guest` user on the Feilong connector into your file `.ssh/authorized_keys`. See the [Global Options](global-options.md) chapter for more details on how to declare the local user to Feilong.

You can then reference those files from your Feilong guests declarations.


### Cloud-init Parameters Sections

The `feilong_cloudinit_params` sections are optional. They may contain the following options:

 * `name` (mandatory): any arbitrary name to identify this resource. Please try to make it unique.
 * `hostname` (mandatory): the desired fully qualified domain name of the z/VM guest.
 * `public_key` (mandatory): the desired public key of the default user on the z/VM guest. The name of the default user depends on the cloud-init settings at the time of the preparation of the image. For SUSE distributions, it is usually `sles`.
