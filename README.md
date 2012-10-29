go-config
=========
The package implements a basic INI configuration file reader. Its interface is heavily influenced by `encoding/csv` Go package.

The configuration file consists of sections, led by a `[section]` header and followed by key-value options (`name = value` or `name=value`). Leading whitespace is removed from values.

For example:

    [Some section]
    foo: bar
    fur = foo

Installation
==========

    go get Claymore/go-config/config

Examples
==========
```go
    package main

    import (
        "github.com/Claymore/go-config/config"
        "fmt"
        "os"
    )

    func main() {
        file, err := os.Open("example.cfg")
        if err != nil {
            fmt.Println("Error:", err)
            return
        }
        defer file.Close()
        reader := config.NewReader(file)
        sections, err := reader.ReadAll()
        if err != nil {
            fmt.Println("Error:", err)
            return
        }
        for _, section := range sections {
            fmt.Println(section)
        }
    }
```

See also
==========
[kless/goconfig](https://github.com/kless/goconfig) for Python ConfigParser flavoured API.