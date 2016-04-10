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
