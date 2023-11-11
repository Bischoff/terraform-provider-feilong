# Terraform Feilong Provider

This terraform provider enables to deploy s390 virtual machines on z/VM via Feilong.

**NOTE:** this is the branch for terraform 1.5.5 (protocol version 6).
The code for terraform 0.13.4 (protocol version 5) is in `terraform-protocol-5` branch.


## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.5.5
- [Go](https://golang.org/doc/install) >= 1.21


## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

The provider will be installed into `$GOPATH/bin`.


## Using the Provider

At this point, the Feilong provider is not yet in HashiCorp's registry. To access it, you must do one of these:

Create a system-wide symbolic link:

```bash
# mkdir -p /usr/share/terraform/plugins/registry.terraform.io/bischoff/feilong/0.0.1/linux_amd64/
# cd /usr/share/terraform/plugins/registry.terraform.io/bischoff/feilong/0.0.1/linux_amd64/
# ln -s <GOPATH>/bin/terraform-provider-feilong
```

Or define this override in the `.terraformrc` file in your home directory:

```terraform
provider_installation {

  dev_overrides {
      "registry.terraform.io/bischoff/feilong" = "<GOPATH>/bin/"
  }

  direct {}
}
```

Replace `<GOPATH>` with the value of your `$GOPATH` environment variable.

In your `main.tf` file, use:

```terraform
terraform {
  required_version = ">= 1.5.5"
  required_providers {
    feilong = {
      source = "bischoff/feilong"
      version = "0.0.1"
    }
  }
}

provider "feilong" {
  connector = "1.2.3.4" // the IP address of your z/VM cloud connector
                        // (i.e. the VM where Feilong runs)
}

resource "feilong_guest" "opensuse" {
  name   = "leap"       // system name for Linux
  userid = "LINUX097"   // system name for system/Z
  vcpus  = 2            // virtual CPUs count
  memory = "2G"         // memory size
  disk = "20G"          // disk size
  image = "opensuse155" // image
}
```


## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

In order to run the full suite of acceptance tests, you need a real Feilong deployment. You also need to upload an image named `testacc`. Once that is done, run `make testacc`:

```bash
$ # specify the address of the Feilong connector
$ export ZVM_CONNECTOR="1.2.3.4"
$ # run the tests
$ make testacc
```

Note: currently hitting this bug: https://github.com/hashicorp/terraform-plugin-testing/issues/185


## License

Apache 2.0, See LICENSE file
