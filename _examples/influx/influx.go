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
