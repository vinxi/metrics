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
	vs.Use(metrics.New(collector(collect)))

	// Target server to forward
	vs.Forward("http://httpbin.org")

	fmt.Printf("Server listening on port: %d\n", port)
	err := vs.Listen()
	if err != nil {
		fmt.Errorf("Error: %s\n", err)
	}
}

// Simple stub collector
type collector func(metrics.Reporter)

func (c collector) Collect(r metrics.Reporter) {
	return c(r)
}

func collect(r metrics.Report) {
	fmt.Printf("Gaudes: %#v\n", r.Gauges)
	fmt.Printf("Counters: %#v\n", r.Counters)
}
