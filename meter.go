package metrics

import (
	"net/http"
	"sync"
	"time"

	"gopkg.in/vinxi/layer.v0"
)

// PublishInterval defines the amount of time to wait between metrics publish cycles.
// Defaults to 15 seconds.
var PublishInterval = 15 * time.Second

// MeterFunc represents the function interface to be implemented by metrics meter functions.
type MeterFunc func(*Info, *Metrics)

// Reporter represents the function interface to be implemented by metrics reporters.
// Metric reporters are responsable of reading, filtering and adapting metrics data.
// Also, reporters tipically sends the metrics to an external service.
type Reporter interface {
	Report(Report) error
}

// Meter provides a metrics instrumentation for vinxi
// Supports configurable metrics reporters and meters.
type Meter struct {
	sync.Mutex
	quit      chan bool
	meters    []MeterFunc
	reporters []Reporter
	metrics   *Metrics
	runtime   *RuntimeCollector
}

// New creates a new metrics meter middleware.
func New(l ...Reporter) *Meter {
	m := &Meter{
		reporters: l,
		meters:    Meters,
		metrics:   NewMetrics(),
		quit:      make(chan bool),
	}

	// Bind the Go runtime stats collector
	m.runtime = NewRuntimeCollector(m.gaugeRuntime)

	// Start collector goroutines
	go m.runtime.Start()
	go m.Start()

	return m
}

// AddMeter adds one or multiple meter functions.
func (m *Meter) AddMeter(meters ...MeterFunc) {
	m.Lock()
	m.meters = append(m.meters, meters...)
	m.Unlock()
}

// SetMeters sets a new set of meter functions, replacing the existent ones.
func (m *Meter) SetMeters(meters []MeterFunc) {
	m.Lock()
	m.meters = meters
	m.Unlock()
}

// AddReporter adds one or multiple metrics reporters.
func (m *Meter) AddReporter(reporters ...Reporter) {
	m.Lock()
	m.reporters = append(m.reporters, reporters...)
	m.Unlock()
}

// Register registers the metrics middleware function.
func (m *Meter) Register(mw layer.Middleware) {
	mw.UsePriority("request", layer.TopHead, m.measureHTTP)
}

// Publish publishes the metrics snapshot report to the registered reporters.
func (m *Meter) Publish() {
	report := m.metrics.Snapshot()
	m.metrics.Reset()

	for _, reporter := range m.reporters {
		go reporter.Report(report)
	}
}

// Start starts a time ticker to publish metrics every certain amount of time.
// You should only call Start when you previously called Stop.
// Start is designed to be executed in its own goroutine.
func (m *Meter) Start() {
	tick := time.NewTicker(PublishInterval)
	for {
		select {
		case <-m.quit:
			return
		case <-tick.C:
			m.Publish()
		}
	}
}

// Stop stops the publish goroutine.
func (m *Meter) Stop() {
	close(m.quit)
	// close(m.runtime.Done) // TODO
}

// measureHTTP instruments and logs an incoming HTTP request and response.
func (m *Meter) measureHTTP(h http.Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		mw := newMetricsWriter(w, r, m.gauge)
		h.ServeHTTP(mw, r)
	}
}

// gauge collects metrics and forward them to the registered meters
func (m *Meter) gauge(i *Info) {
	for _, meter := range m.meters {
		meter(i, m.metrics)
	}
}

// gaugeRuntime collects runtime metrics and stores it in a histogram.
func (m *Meter) gaugeRuntime(key string, val uint64) {
	m.metrics.Histogram(key).RecordValue(int64(val))
}
