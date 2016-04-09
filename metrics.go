package metrics

import (
	"net/http"
	"sync"
	"time"

	"github.com/codahale/metrics"
)

// Info is used in meter functions to access to collected data from the response writer.
type Info struct {
	// Status stores the response HTTP status.
	Status int
	// BodyLength stores the response body length in bytes.
	BodyLength int64
	// TimeStart stores when the request was received by the server.
	TimeStart time.Time
	// TimeEnd stores when the response is written.
	TimeEnd time.Time
	// Header stores the response HTTP header.
	Header http.Header
	// Request points to the original http.Request instance.
	Request *http.Request
}

// Report is used to expose Counters and Gauges collected via Metrics.
type Report struct {
	// Gauges stores metrics gauges values accesible by key.
	Gauges map[string]int64
	// Counters stores the metrics counters accesible by key.
	Counters map[string]uint64
}

// Metrics is used to temporary store metrics data of multiple origins and nature.
// Provides a simple interface to write and read metric values.
//
// Metrics is designed to be safety used by multiple goroutines.
type Metrics struct {
	// Mutex provides synchronization for thead safety.
	sync.Mutex
	// Gauges stores gauges
	Gauges map[string]metrics.Gauge
	// Counters stores counters by key
	Counters map[string]metrics.Counter
	// Histograms stores histograms by key.
	Histograms map[string]*metrics.Histogram
}

// NewMetrics creates a new metrics object for reporting.
func NewMetrics() *Metrics {
	m := &Metrics{}
	m.Reset()
	return m
}

// Counter returns a counter metric by key.
// If the counter doesn't exists, it will be transparently created.
func (m *Metrics) Counter(key string) metrics.Counter {
	m.Lock()
	defer m.Unlock()
	counter, ok := m.Counters[key]
	if !ok {
		counter = metrics.Counter(key)
		m.Counters[key] = counter
	}
	return counter
}

// Guige returns a guide metric by key.
// If the guige doesn't exists, it will be transparently created.
func (m *Metrics) Guage(key string) metrics.Gauge {
	m.Lock()
	defer m.Unlock()
	gauge, ok := m.Gauges[key]
	if !ok {
		gauge = metrics.Gauge(key)
		m.Gauges[key] = gauge
	}
	return gauge
}

// Histogram returns an histrogram by key.
// If the histogram doesn't exists, it will be transparently created.
func (m *Metrics) Histogram(key string) *metrics.Histogram {
	m.Lock()
	defer m.Unlock()
	hist, ok := m.Histograms[key]
	if !ok {
		hist = metrics.NewHistogram(key, 0, 1e8, 5)
		m.Histograms[key] = hist
	}
	return hist
}

// Snapshot collects and returns a report of the existent counters and gauges metrics
// to be consumed by metrics publishers and listeners.
func (m *Metrics) Snapshot() Report {
	c, g := metrics.Snapshot()
	return Report{Gauges: g, Counters: c}
}

// Reset resets all the metrics (counters, gauges & histograms) to zero.
// You should collect them first with Snapshop(), otherview the collected data will be lost.
func (m *Metrics) Reset() {
	metrics.Reset()
	m.Lock()
	m.Gauges = make(map[string]metrics.Gauge)
	m.Counters = make(map[string]metrics.Counter)
	m.Histograms = make(map[string]*metrics.Histogram)
	m.Unlock()
}
