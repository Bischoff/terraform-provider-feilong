# Terraform Feilong Provider

This terraform provider enables to deploy s390 virtual machines on z/VM via Feilong.

**NOTE:** this is the branch for terraform 0.13.4 (protocol version 5).
The code for terraform 1.0.10 (protocol version 6) is in `main` branch.


## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 0.13.4 or [OpenTofu](https://opentofu.org/docs/intro/install/) >= 1.6.0
- [Go](https://golang.org/doc/install) >= 1.21


## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

The provider will be installed into `$GOPATH/bin`.

The Feilong provider is in HashiCorp's or Opentofu's registry. To bypass the registries, you can do one of these:

Create a system-wide symbolic link:

```bash
# -- Terraform --
# mkdir -p /usr/share/terraform/plugins/registry.terraform.io/bischoff/feilong/0.0.9/linux_amd64/
# cd /usr/share/terraform/plugins/registry.terraform.io/bischoff/feilong/0.0.9/linux_amd64/
# ln -s <GOPATH>/bin/terraform-provider-feilong
# -- OpenTofu --
# mkdir -p /usr/share/terraform/plugins/registry.opentofu.org/bischoff/feilong/0.0.9/linux_amd64/
# cd /usr/share/terraform/plugins/registry.opentofu.org/bischoff/feilong/0.0.9/linux_amd64/
# ln -s <GOPATH>/bin/terraform-provider-feilong
```

Or define this override in the `.terraformrc` or `.tofurc` file in your home directory:

```terraform
provider_installation {

  dev_overrides {
    "registry.opentofu.org/bischoff/feilong" = "<GOPATH>/bin/"
  }

  direct {}
}
```

Replace `<GOPATH>` with the value of your `$GOPATH` environment variable.


## Using the Provider

In your `main.tf` file, use:

```terraform
terraform {
  required_version = ">= 0.13.4"
  required_providers {
    feilong = {
      source = "bischoff/feilong"
      version = "0.0.9"
    }
  }
}

provider "feilong" {
  connector = "http://1.2.3.4"     // URL of your z/VM cloud connector
                                   // (i.e. the VM where Feilong runs)
}

resource "feilong_guest" "opensuse" {
  name   = "leap"                  // arbitrary name for the resource
  memory = "2G"                    // memory size
  disk   = "20G"                   // disk size of first disk
  image  = "opensuse155"           // image

  // optional parameters:
  userid = "LINUX097"              // name for z/VM
  vcpus  = 2                       // virtual CPUs count
  mac    = "12:34:56:78:9a:bc"     // MAC address of first interface
                                   // (first 3 bytes may be changed by Feilong)
}
```

Then use Terraform or OpenTofu commands:

```bash
$ terraform init
$ terraform apply
(use the VMs)
$ terraform destroy
```

```bash
$ tofu init
$ tofu apply
(use the VMs)
$ tofu destroy
```

For more details, refer to the [documentation](docs/README.md).


## To Do

* Write missing CRUD functions:
  * network configuration Update()
  * finish vswitch Update()
  * cloudinit Read()
  * cloudinit Update()
* Support more z/VM resources:
  * additional network interfaces
  * minidisk
  * fiber channel
  * other?
* Resurrect acceptance tests

Your help is welcome!


## License

Apache 2.0, See LICENSE file
