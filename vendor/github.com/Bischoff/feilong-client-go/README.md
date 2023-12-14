# Go Client Library for Feilong

This library enables to deploy s390 virtual machines on z/VM via [Feilong](https://openmainframeproject.org/projects/feilong/) from a Go program.


## Requirements

- [Go](https://golang.org/doc/install) >= 1.21


## Using the Library

Here is a small sample program using the library:

```go
package main

import (
        "fmt"
        "os"

        "github.com/Bischoff/feilong-client-go"
)

func main() {
        connector := os.Getenv("ZVM_CONNECTOR")

        client := feilong.NewClient(&connector, nil)

        result, err := client.GetFeilongVersion()
        if err != nil {
                fmt.Println(err.Error())
                return
        }

        fmt.Printf("API version: %s\n", result.Output.APIVersion)
}
```

For more examples, look in [snippets](examples/snippets/) directory.


## Naming Conventions

The following conventions are used both for function names and for input and output structures:

 * names use camel case, with no underscores: `user_profile` in Feilong JSON API becomes `UserProfile` in this Go library
 * acronyms are completly capitalized: `api_version` becomes `APIVersion`
 * abbreviations are sometimes completed, but never fully capitalized: `max_version` becomes `MaxVersion`, and `modID` becomes `ModuleId`
 * the Feilong name is not completely respected if it does not respect English grammar or is otherwise wrong: `GetGuestsList` becomes `GetGuestList`.

Please refer to the individual definitions to know the exact names used in this library.


## Completeness

The library implements the [Feilong API](https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#) version 1.0.

The following are not implemented yet:

 * documentation
 * acceptance tests
 * many API functions (see below).


### Implemented Functions

The numbers below refer to the section numbers in the Feilong documentation.

 * 7.2 - [Version](version.go)
   * 7.2.1 - `GetVersion()`
 * 7.3 - [Token](token.go)
   * 7.3.1 - `CreateToken()`
 * 7.4 - [SMAPI](smapi.go)
   * 7.4.1 - `SMAPIHealth()`
 * 7.5 - [Guests](guests.go)
   * 7.5.1 - `ListGuests()`
   * 7.5.2 - `CreateGuest()`
   * 7.5.3 - `GuestAddDisks()`
   * 7.5.5 - `GuestDeleteDisks()`
   * 7.5.15 - `ShowGuestDefinition()`
   * 7.5.16 - `DeleteGuest()`
   * 7.5.18 - `GetGuestInfo()`
   * 7.5.20 - `GetGuestAdaptersInfo()`
   * 7.5.21 - `CreateGuestNIC()`
   * 7.5.24 - `StartGuest()`
   * 7.5.25 - `StopGuest()`
   * 7.5.39 - `DeployGuest()`
   * 7.5.43 - `UpdateGuestNIC()`
 * 7.6 - [Host](host.go)
   * 7.6.1 - `GetGuestList()`
   * 7.6.2 - `GetHostInfo()`
   * 7.6.3 - `GetHostDiskPoolInfo()`, `GetHostDiskPoolDetails()`
 * 7.7 - [Images](images.go)
   * 7.7.1 - `ListImages()`
   * 7.7.2 - `CreateImage()`
   * 7.7.3 - `ExportImage()`
   * 7.7.4 - `GetRootDiskSize()`
   * 7.7.5 - `DeleteImage()`
 * 7.8 - [VSwitches](vswitches.go)
   * 7.8.1 - `CreateVSwitch()`
   * 7.8.2 - `ListVSwitches()`
   * 7.8.3 - `GetVSwitchDetails()`
   * 7.8.7 - `DeleteVSwitch()`
 * 7.9 - [Files](file.go)
   * 7.9.1 - `ImportFile()`
   * 7.9.2 - `ExportFile()`


## License

Apache 2.0, See LICENSE file
