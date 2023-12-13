## Global Parameters

Here is an example of global parameters:

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
  connector   = "feilong.example.org"
  admin_token = "zvX2mFxuj8HcrYkAacLReV0RTQ0K5IIEighOR9F8AG"
  local_user  = "johndoe@client.example.org"
}
```

### Terraform Section

The `terraform` section is mandatory. The possible variables are:

 * `required_version`: the version of Terraform itself.
 * `require_providers`: a map of all providers to use in this deployment. For each of them, the `source` (from the terraform registry) and the `version` are needed.

You can bypass the terraform registry by compiling the Terraform provider yourself, and by creating the a symbolic link like this:

```bash
# cd /usr/share/terraform/plugins/registry.terraform.io/bischoff/feilong/0.0.4/linux_amd64/
# ln -s /home/someuser/bin/terraform-provider-feilong
```


### Feilong Provider Section

The `provider` section is mandatory. The possible variables are:

 * `connector` (optional): the IP address or domain name of the z/VM connector, i.e. the VM where Feilong runs. If you omit it, the value will be taken from environment variable `$ZVM_CONNECTOR`.
 * `admin_token` (optional): the secret shared with the z/VM connector for authentication, in case this was set up. See the [Token Usage](https://cloudlib4zvm.readthedocs.io/en/latest/setuphttpd.html#token-usage) chapter of the Feilong documentation for more information on how to set this up. If you don't want to store it in the `main.tf` file, you can use [Terraform variables](https://developer.hashicorp.com/terraform/language/values/variables) to pass it at run time.
 * `local_user` (optional): user name and IP address or domain name of the workstation where you run terraform. You need to specify it if you intend to use cloud-init parameters and/or network parameters. In that case, you must drop the public SSH key of the z/VM connector into file `.shh/authorized_keys` in the home directory of that user. This will allow Feilong to upload the cloud-init parameters file and/or the network parameters file.
