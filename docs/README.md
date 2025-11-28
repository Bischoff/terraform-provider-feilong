## Terraform Provider for Feilong

### Introduction and Ecosystem

The Terraform provider for Feilong enables to dynamically deploy guests on a s/390 system that runs z/VM operating system. It leverages [Feilong](https://openmainframeproject.org/projects/feilong/) connector, which is at the core of the [Cloud Infrastructure Center](https://www.ibm.com/products/cloud-infrastructure-center) of IBM. However, Feilong can be used alone, without a full CIC deployment.

Feilong offers an HTTP REST API that manages VMs by transmitting requests to the more complex [Systems Management API](https://www.ibm.com/docs/en/zvm/7.2?topic=introduction-smapi-quick-start-guide), also known as SMAPI. The [Go library for Feilong](https://github.com/Bischoff/feilong-client-go) allows to call that REST API in a simple manner from a Go program. The Terraform provider for Feilong relies on that library and is seen by the [Terraform](https://www.terraform.io/) automation tool as a plugin.

Terraform relies on `main.tf` files for deploying VMs. The Feilong provider extends the syntax of those configuration files to take into account the specificities of z/VM.

A normal cycle of commands to use the `main.tf` file is:
```bash
$ tofu init
$ tofu apply
(use the VMS)
$ tofu destroy
```
or
```bash
$ terraform init
$ terraform apply
(use the VMS)
$ terraform destroy
```


### Global Structure of the Configuration File

The `main.tf` file is made of the following sections:

```terraform
terraform {
  (...)
}

provider "feilong" {
  (...)
}

(other providers)


resource "feilong_cloudinit_params" "(some name)" {
  (...)
}

resource "feilong_vswitch" "(some name)" {
  (...)
}

resource "feilong_guest" "(some name)" {
  (...)
}

(other resources for Feilong)

(other resources for the other providers)


output "feilong_guest_mac_address" {
  value = feilong_guest.(some name).mac_address
}

output "feilong_guest_ip_address" {
  value = feilong_guest.(some name).ip_address
}

(other values for output)
```

The `terraform` section allows to define the required versions of terraform and of the various providers. The `provider` section allows to define global settings, like the IP address or the domain name of the Feilong connector. Both are described more in details in [Global Parameters](global-options.md) chapter.

The `feilong_cloudinit_params` resource sections allow to create locally a file that can be used to store parameters for [cloud-init](https://github.com/canonical/cloud-init) during the initial deployment. It is described more in details in [Local Files](local-files.md) chapter.

The `feilong_vswitch` resource sections allow to create s/390 [virtual switches](https://www.redbooks.ibm.com/redbooks/pdfs/sg247023.pdf), in the case that the existing vswitches do not match your needs. They are decribed more in details in [Virtual Switches](virtual-switches.md) chapter.

The `feilong_guest` resource sections allow to create s/390 guest VMs (`userid`s in z/VM parlance). They are described more in details in [Guests](guests.md) chapter.

The `output` sections allow to display computed values at the end of the terraform deployment. These are values that were unknown at the start of the deployment.

Terraform also has the notion of "data sources". They are currently not used in the Terraform provider for Feilong.
