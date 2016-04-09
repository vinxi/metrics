# metrics [![Build Status](https://travis-ci.org/vinxi/metrics.png)](https://travis-ci.org/vinxi/metrics) [![GoDoc](https://godoc.org/github.com/vinxi/metrics?status.svg)](https://godoc.org/github.com/vinxi/metrics) [![Coverage Status](https://coveralls.io/repos/github/vinxi/metrics/badge.svg?branch=master)](https://coveralls.io/github/vinxi/metrics?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/vinxi/metrics)](https://goreportcard.com/report/github.com/vinxi/metrics)

Simple metrics instrumentation for your proxy. 
Supports counters, gauges and histograms analyzing multiple scopes of HTTP traffic.

Uses [codahale/metrics](https://github.com/codahale/metrics) under the hood.

**Work in progress**

## Installation

```bash
go get -u gopkg.in/vinxi/metrics.v0
```

## API

See [godoc](https://godoc.org/github.com/vinxi/metrics) reference.

## Example

#### Default log to stdout

```go
package main

import (
  "fmt"
  "gopkg.in/vinxi/metrics.v0"
  "gopkg.in/vinxi/vinxi.v0"
)

const port = 3100

func main() {
  // Create a new vinxi proxy
  vs := vinxi.NewServer(vinxi.ServerOptions{Port: port})
  
  // Intrument with default metrics reporter 
  vs.Use(metrics.Default)
  
  // Target server to forward by default
  vs.Forward("http://httpbin.org")

  fmt.Printf("Server listening on port: %d\n", port)
  err := vs.Listen()
  if err != nil {
    fmt.Errorf("Error: %s\n", err)
  }
}
```

## License

MIT
