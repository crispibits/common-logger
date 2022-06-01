# common-logger
Provides a consistent logger across projects.
This is a very thin wrapper for [Uber's Zap](https://github.com/uber-go/zap).

It just provides two profiles which work for my projects - `console` and `server`.  `console` is easily read on the CLI, whereas `server` provides json formatted logs suitable for use with GKE Stackdriver logging

Usage:

```go
package main

import (
    log "github.com/crispibits/common-logger"
)

func main() {
    l := log.New()
    l.Debug("Hello")
}
```

The logging can then be changed by running:
```sh
PROFILE=console go run main.go
```
