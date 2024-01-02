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

For more examples, look at the [examples](examples/README.md) directory.

For an overview of all the functions, see the [documentation](docs/README.md).


## Compliance

The library implements the [Feilong API](https://cloudlib4zvm.readthedocs.io/en/latest/restapi.html#) version 1.0. It has been tested to work with Feilong server version 1.6.7.


## To Do

The following are not done yet:

 * detailed documentation
 * acceptance tests

Also, volume operations are probably not sufficiently tested because of missing opportunities to test it.

Your help is welcome!


## License

Apache 2.0, See LICENSE file
