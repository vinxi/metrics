package metrics

import (
	"net/http"
	"testing"
	"time"

	"github.com/nbio/st"
)

func TestMeterNumberOfRequests(t *testing.T) {
	info, metrics := createMetrics()
	defer metrics.Reset()
	MeterNumberOfRequests(info, metrics)
	st.Expect(t, len(metrics.Counters), 1)
	st.Expect(t, metrics.Snapshot().Counters["req.total"], uint64(1))
}

func TestMeterResponseStatus(t *testing.T) {
	info, metrics := createMetrics()
	MeterResponseStatus(info, metrics)
	st.Expect(t, len(metrics.Counters), 1)
	st.Expect(t, metrics.Snapshot().Counters["res.status.ok"], uint64(1))
	metrics.Reset()

	info, metrics = createMetrics()
	info.Status = 400
	MeterResponseStatus(info, metrics)
	st.Expect(t, len(metrics.Counters), 1)
	st.Expect(t, metrics.Snapshot().Counters["res.status.ok"], uint64(0))
	st.Expect(t, metrics.Snapshot().Counters["res.status.bad"], uint64(1))
	metrics.Reset()

	info, metrics = createMetrics()
	info.Status = 500
	MeterResponseStatus(info, metrics)
	st.Expect(t, len(metrics.Counters), 1)
	st.Expect(t, metrics.Snapshot().Counters["res.status.ok"], uint64(0))
	st.Expect(t, metrics.Snapshot().Counters["res.status.bad"], uint64(0))
	st.Expect(t, metrics.Snapshot().Counters["res.status.error"], uint64(1))
	metrics.Reset()
}

func TestMeterRequestOperation(t *testing.T) {
	info, metrics := createMetrics()
	MeterRequestOperation(info, metrics)
	st.Expect(t, len(metrics.Counters), 1)
	st.Expect(t, metrics.Snapshot().Counters["req.reads"], uint64(1))
	metrics.Reset()

	info, metrics = createMetrics()
	info.Request.Method = "POST"
	MeterRequestOperation(info, metrics)
	st.Expect(t, len(metrics.Counters), 1)
	st.Expect(t, metrics.Snapshot().Counters["req.reads"], uint64(0))
	st.Expect(t, metrics.Snapshot().Counters["req.writes"], uint64(1))
	metrics.Reset()
}

func TestMeterResponseTime(t *testing.T) {
	info, metrics := createMetrics()
	defer metrics.Reset()
	MeterResponseTime(info, metrics)
	st.Expect(t, metrics.Snapshot().Gauges["res.time.P50"] >= 100, true)
	st.Expect(t, metrics.Snapshot().Gauges["res.time.P99"] >= 100, true)
	st.Expect(t, metrics.Snapshot().Gauges["res.time.P999"] >= 100, true)
}

func TestMeterResponseBodySize(t *testing.T) {
	info, metrics := createMetrics()
	defer metrics.Reset()
	MeterResponseBodySize(info, metrics)
	st.Expect(t, metrics.Snapshot().Gauges["res.body.size.P50"], int64(10))
	st.Expect(t, metrics.Snapshot().Gauges["res.body.size.P99"], int64(10))
	st.Expect(t, metrics.Snapshot().Gauges["res.body.size.P999"], int64(10))
}

func TestMeterRequestBodySize(t *testing.T) {
	info, metrics := createMetrics()
	defer metrics.Reset()
	MeterRequestBodySize(info, metrics)
	st.Expect(t, metrics.Snapshot().Gauges["req.body.size.P50"], int64(10))
	st.Expect(t, metrics.Snapshot().Gauges["req.body.size.P99"], int64(10))
	st.Expect(t, metrics.Snapshot().Gauges["req.body.size.P999"], int64(10))
}

func createMetrics() (*Info, *Metrics) {
	metrics := NewMetrics()
	info := &Info{
		Status:     200,
		BodyLength: 10 * 1024,
		TimeStart:  time.Now(),
		TimeEnd:    time.Now().Add(100 * time.Millisecond),
		Request:    &http.Request{Method: "GET", ContentLength: 10 * 1024},
	}
	return info, metrics
}
