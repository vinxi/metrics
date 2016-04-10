package metrics

import (
	"math"
	"time"
)

// Meters stores the built-in function meters used by default for metrics collection.
// You can define your custom meter functions via metrics.AddMeter() or metrics.SetMeters().
var Meters = []MeterFunc{
	MeterNumberOfRequests,
	MeterResponseStatus,
	MeterRequestOperation,
	MeterResponseTime,
	MeterResponseBodySize,
	MeterRequestBodySize,
}

// MeterNumberOfRequests is used to register the total number of served requests.
func MeterNumberOfRequests(i *Info, m *Metrics) {
	m.Counter("req.total").Add()
}

// MeterResponseStatus is used to count the response status code by range (2xx, 4xx, 5xx).
func MeterResponseStatus(i *Info, m *Metrics) {
	s := i.Status / 100
	if s >= 2 && s < 4 {
		m.Counter("res.status.ok").Add()
	} else if s == 5 {
		m.Counter("res.status.error").Add()
	} else if s == 4 {
		m.Counter("res.status.bad").Add()
	}
}

// MeterRequestOperation is used to count the number of request by HTTP operation.
// Operation in inferred by HTTP verb:
//
// - GET, HEAD = read operation
// - POST, PUT, PATCH, DELETE = write operation
func MeterRequestOperation(i *Info, m *Metrics) {
	if i.Request.Method == "GET" || i.Request.Method == "HEAD" {
		m.Counter("req.reads").Add()
		return
	}

	if i.Request.Method == "POST" || i.Request.Method == "PUT" ||
		i.Request.Method == "PATCH" || i.Request.Method == "DELETE" {
		m.Counter("req.writes").Add()
	}
}

// MeterResponseTime is used to measure the HTTP request/response time.
// Data will be stored in a histogram.
func MeterResponseTime(i *Info, m *Metrics) {
	resTime := i.TimeEnd.Sub(i.TimeStart).Nanoseconds() / int64(time.Millisecond)
	m.Histogram("res.time").RecordValue(resTime)
}

// MeterResponseBodySize is used to measure the HTTP response body length.
// Data will be stored in a histogram.
func MeterResponseBodySize(i *Info, m *Metrics) {
	if i.BodyLength > 0 {
		m.Histogram("res.body.size").RecordValue(toKB(i.BodyLength))
	}
}

// MeterRequestBodySize is used to measure the HTTP request body length.
// Data will be stored in a histogram.
func MeterRequestBodySize(i *Info, m *Metrics) {
	if i.Request.ContentLength > 0 {
		m.Histogram("req.body.size").RecordValue(toKB(i.Request.ContentLength))
	}
}

// toKB converts n bytes into KB.
func toKB(n int64) int64 {
	return int64(math.Floor((float64(n) / 1024) + 0.5))
}
