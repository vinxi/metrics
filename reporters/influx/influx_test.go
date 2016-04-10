package influx

import (
	"testing"

	client "github.com/influxdata/influxdb/client/v2"
	"github.com/nbio/st"
	"gopkg.in/vinxi/metrics.v0"
)

var testConfig = Config{URL: "http://foo"}

func TestMapReportCounters(t *testing.T) {
	counters := make(map[string]uint64)
	counters["foo"] = 100
	report := metrics.Report{Counters: counters}
	reporter := New(testConfig)

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{})
	reporter.mapReport(report, bp)

	st.Expect(t, len(bp.Points()), 1)
	st.Expect(t, bp.Points()[0].Name(), "foo.count")
	st.Expect(t, bp.Points()[0].Fields()["value"], int64(100))
}

func TestMapReportGauges(t *testing.T) {
	gauges := make(map[string]int64)
	gauges["foo"] = 100
	report := metrics.Report{Gauges: gauges}
	reporter := New(testConfig)

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{})
	reporter.mapReport(report, bp)

	st.Expect(t, len(bp.Points()), 1)
	st.Expect(t, bp.Points()[0].Name(), "foo.gauge")
	st.Expect(t, bp.Points()[0].Fields()["value"], int64(100))
}

func TestMapReportHistograms(t *testing.T) {
	gauges := make(map[string]int64)
	gauges["foo.P50"] = 50
	gauges["foo.P75"] = 75
	gauges["foo.P90"] = 90
	gauges["foo.P95"] = 95
	gauges["foo.P99"] = 99
	gauges["foo.P999"] = 999
	report := metrics.Report{Gauges: gauges}
	reporter := New(testConfig)

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{})
	reporter.mapReport(report, bp)

	st.Expect(t, len(bp.Points()), 1)
	histogram := bp.Points()[0]
	st.Expect(t, histogram.Name(), "foo.histogram")
	fields := histogram.Fields()
	st.Expect(t, fields["p50"], int64(50))
	st.Expect(t, fields["p75"], int64(75))
	st.Expect(t, fields["p90"], int64(90))
	st.Expect(t, fields["p95"], int64(95))
	st.Expect(t, fields["p99"], int64(99))
	st.Expect(t, fields["p999"], int64(999))
}

func TestInfluxDataReport(t *testing.T) {
	// TODO: mock InfluxDB server and assert reported JSON data
}
