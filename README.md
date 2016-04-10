# metrics [![Build Status](https://travis-ci.org/vinxi/metrics.png)](https://travis-ci.org/vinxi/metrics) [![GoDoc](https://godoc.org/github.com/vinxi/metrics?status.svg)](https://godoc.org/github.com/vinxi/metrics) [![Coverage Status](https://coveralls.io/repos/github/vinxi/metrics/badge.svg?branch=master)](https://coveralls.io/github/vinxi/metrics?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/vinxi/metrics)](https://goreportcard.com/report/github.com/vinxi/metrics)

Simple, extensible metrics instrumentation for your proxies. 

Collects useful and versatile metrics of duplex HTTP traffic flows and Go runtime stats.

Supports `counters`, `gauges` and `histogram` with `50`, `75`, `90`, `95`, `99` and `999` percentiles.
Uses [codahale/metrics](https://github.com/codahale/metrics) under the hood.

## Reporters

Reporters are pluggable components that reads metric reports and tipically sends 
it to an data ingestor provider. 

You can write and plug in your own reporter. 
See [how to write reporters](#writting-reporters) section.

Built-in supported reporters:

- [x] [InfluxDB](https://github.com/vinxi/metrics/tree/master/reporters/influx)
- [ ] Statsd
- [ ] Prometheus

## Meters

Meters are simple functions that reads HTTP request/response info and generates 
further counters, gauges or histograms based on it.

`metrics` package allows you to easily extend meter function in order to measure 
custom or new properies of the HTTP flow to cover your specific needs.
See [how to write meters](#writting-meters) section.

Default provided meters (listed as: description, measure type, metric name):

- Total requests - `counter` - `req.total.count`
- Total success responses - `counter` - `res.status.ok.count`
- Total error responses - `counter` - `res.status.error.count`
- Total bad responses - `counter` - `res.status.bad.count`
- Total read requests - `counter` - `req.reads.count`
- Total write requests - `counter` - `req.writes.count`
- Response time in milliseconds - `histogram` - `res.time.histogram`
- Response body size in KB - `histogram` - `res.body.size.histogram`
- Request body size in KB - `histogram` - `req.body.size.histogram`

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
  Report(metrics.Report) error
}
```

The metrics publisher will call the `Report` method passing the `Report` struct, 
which exports the fields `Counters` and `Gauges`.

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

## Writting meters

Meters are simple functions implementing the following function signature:

```go
type MeterFunc func(*metrics.Info, *metrics.Metrics)
```

#### Meter example

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

  // Create a custom meter function the increases a counter
  // when the response status is 200
  myMeter := func(i *metrics.Info, m *metrics.Metrics) {
    if i.Status == 200 {
      m.Counter("res.success.total").Add()
    }
  }

  // Create a new metrics middleware
  m := metrics.New(reporter(collect))
  // Add the custom meter
  m.AddMeter(myMeter)
  // Attach the metrics middleware
  vs.Use(m)

  // Target server to forward
  vs.Forward("http://httpbin.org")

  fmt.Printf("Server listening on port: %d\n", port)
  err := vs.Listen()
  if err != nil {
    fmt.Errorf("Error: %s\n", err)
  }
}

// Simple stub reporter
type reporter func(metrics.Report) error

func (c reporter) Report(r metrics.Report) error {
  return c(r)
}

func collect(r metrics.Report) error {
  fmt.Printf("Gaudes: %#v\n", r.Gauges)
  fmt.Printf("Counters: %#v\n", r.Counters)
  return nil
}
```

## License

MIT
