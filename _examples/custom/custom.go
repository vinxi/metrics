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

	// Attach the metrics middleware
	vs.Use(metrics.New(reporter(collect)))

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
