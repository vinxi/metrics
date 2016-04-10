package metrics

import (
	"testing"

	"github.com/nbio/st"
)

func TestMetrics(t *testing.T) {
	metrics := NewMetrics()
	defer metrics.Reset()

	metrics.Counter("foo").Add()
	st.Expect(t, metrics.Snapshot().Counters["foo"], uint64(1))

	metrics.Guage("foo").Set(1)
	st.Expect(t, metrics.Snapshot().Gauges["foo"], int64(1))

	metrics.Histogram("foo").RecordValue(100)
	st.Expect(t, metrics.Snapshot().Gauges["foo.P50"], int64(100))
	st.Expect(t, metrics.Snapshot().Gauges["foo.P75"], int64(100))
	st.Expect(t, metrics.Snapshot().Gauges["foo.P90"], int64(100))
	st.Expect(t, metrics.Snapshot().Gauges["foo.P95"], int64(100))
	st.Expect(t, metrics.Snapshot().Gauges["foo.P99"], int64(100))
	st.Expect(t, metrics.Snapshot().Gauges["foo.P999"], int64(100))
}
