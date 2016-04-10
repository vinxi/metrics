package metrics

import (
	"net/http"
	"strconv"
	"time"
)

type collector func(*Info)

// metricWriter implements a http.ResponseWriter capable interface used
// to intercept the HTTP response, register useful data and finally
// call data collector function.
type metricWriter struct {
	info      *Info
	collector collector
	w         http.ResponseWriter
}

// newMetricsWriter creates a new metrics writer
func newMetricsWriter(w http.ResponseWriter, r *http.Request, collector collector) *metricWriter {
	info := &Info{TimeStart: time.Now(), Request: r, Header: w.Header()}
	return &metricWriter{w: w, info: info, collector: collector}
}

// Header implements http.ResponseWriter Header method.
func (l *metricWriter) Header() http.Header {
	return l.w.Header()
}

// WriteHeader implements http.ResponseWriter WriteHeader method.
func (l *metricWriter) WriteHeader(code int) {
	if l.info.Status != 0 {
		return
	}

	l.info.Status = code
	l.info.TimeEnd = time.Now()
	l.info.BodyLength, _ = strconv.ParseInt(l.w.Header().Get("Content-Length"), 10, 64)
	defer l.collector(l.info)

	l.w.WriteHeader(code)
}

// Write implements http.ResponseWriter Write method.
func (l *metricWriter) Write(buf []byte) (int, error) {
	return l.w.Write(buf)
}
