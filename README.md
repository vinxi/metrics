# metrics [![Build Status](https://travis-ci.org/vinxi/metrics.png)](https://travis-ci.org/vinxi/metrics) [![GoDoc](https://godoc.org/github.com/vinxi/metrics?status.svg)](https://godoc.org/github.com/vinxi/metrics) [![Coverage Status](https://coveralls.io/repos/github/vinxi/metrics/badge.svg?branch=master)](https://coveralls.io/github/vinxi/metrics?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/vinxi/metrics)](https://goreportcard.com/report/github.com/vinxi/metrics)

Simple, extensible metrics instrumentation for your proxies. 
Collects useful and versatile metrics of duplex HTTP traffic flow and Go runtime stats.

Uses [codahale/metrics](https://github.com/codahale/metrics) under the hood.

## Reporters

- [x] InfluxDB
- [ ] Statsd
- [ ] Prometheus

You can write and plug in your own reporter. 
See [how to write reporters](#writting-reporters) section.

## Installation

```bash
go get -u gopkg.in/vinxi/metrics.v0
```

## API

See [godoc](https://godoc.org/github.com/vinxi/metrics) reference.

## Examples

#### Report metrics to InfluxDB

```go
package main

import (
  "fmt"
  "gopkg.in/vinxi/metrics.v0"
  "gopkg.in/vinxi/metrics.v0/reporters/influx"
  "gopkg.in/vinxi/vinxi.v0"
)

const port = 3100

func main() {
  // Create a new vinxi proxy
  vs := vinxi.NewServer(vinxi.ServerOptions{Port: port})

  // Attach the metrics middleware
  config := influx.Config{
    URL:      "http://localhost:8086",
    Username: "root",
    Password: "root",
    Database: "metrics",
  }
  vs.Use(metrics.New(influx.New(config)))

  // Target server to forward
  vs.Forward("http://httpbin.org")

  fmt.Printf("Server listening on port: %d\n", port)
  err := vs.Listen()
  if err != nil {
    fmt.Errorf("Error: %s\n", err)
  }
}
```

#### Report metrics only for certain scenarios via multiplexer

```go
package main

import (
  "fmt"
  "gopkg.in/vinxi/metrics.v0"
  "gopkg.in/vinxi/metrics.v0/reporters/influx"
  "gopkg.in/vinxi/mux.v0"
  "gopkg.in/vinxi/vinxi.v0"
)

const port = 3100

func main() {
  // Create a new vinxi proxy
  vs := vinxi.NewServer(vinxi.ServerOptions{Port: port})

  // InfluxDB reporter config
  config := influx.Config{
    URL:      "http://localhost:8086",
    Username: "root",
    Password: "root",
    Database: "metrics",
  }

  // Attach the metrics middleware via muxer
  mx := mux.If(mux.Method("GET", "POST"), mux.Path("/"))
  mx.Use(metrics.New(influx.New(config)))
  vs.Use(mx)

  // Target server to forward
  vs.Forward("http://httpbin.org")

  fmt.Printf("Server listening on port: %d\n", port)
  err := vs.Listen()
  if err != nil {
    fmt.Errorf("Error: %s\n", err)
  }
}
```

## Writting reporters

`metrics` package allows you to write and plug in custom reporters in order to send data to third-party
metrics and storage providers. 

Reporters must implement the `Reporter` interface, which consists is a single method:

```go
type Reporter interface {
  Report(Report)
}
```

The metrics publisher will call the `Report` method passing the `Report` struct, which exports 
two fields containing the counters and gauges as a simple map data type.

#### Reporter example

```go
import (
  "fmt"
  "gopkg.in/vinxi/metrics.v0"  
)

type MyCustomReporter struct {
  // reporter specific fields  
}

func (m *MyCustomReporter) Report(r metrics.Report) {
  // Print maps
  fmt.Printf("Counters: %#v \n", r.Counters)
  fmt.Printf("Gauges: %#v \n", r.Gauges)
  
  // Here you should usually map and transform the metrics report
  // into reporter specific data structures.
  data := mapReport(r)
    
  // Finally send the metrics, tipically via network to another server
  reporterClient.Send(data)
}
```

## License

MIT
