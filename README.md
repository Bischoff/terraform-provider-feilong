# Terraform Feilong Provider

This terraform provider enables to deploy s390 virtual machines on z/VM via Feilong.

**NOTE:** this is the branch for terraform 0.13.4 (protocol version 5).
The code for terraform 1.0.10 (protocol version 6) is in `main` branch.


## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 0.13.4
- [Go](https://golang.org/doc/install) >= 1.21


## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

The provider will be installed into `$GOPATH/bin`.

The Feilong provider is in HashiCorp's registry. To bypass the registry, you can
create a system-wide symbolic link:

```bash
# mkdir -p /usr/share/terraform/plugins/registry.terraform.io/bischoff/feilong/0.0.4/linux_amd64/
# cd /usr/share/terraform/plugins/registry.terraform.io/bischoff/feilong/0.0.4/linux_amd64/
# ln -s <GOPATH>/bin/terraform-provider-feilong
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
      version = "0.0.4"
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

Then use Terraform commands:

```bash
$ terraform init
$ terraform apply
(use the VMs)
$ terraform destroy
```

For more details, refer to the [documentation](docs/README.md).


## License

Apache 2.0, See LICENSE file
