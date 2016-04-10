package metrics

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/nbio/st"
	"gopkg.in/vinxi/utils.v0"
)

func TestMeter(t *testing.T) {
	metrics := &Meter{
		meters:  Meters,
		metrics: NewMetrics(),
		quit:    make(chan bool),
	}
	defer metrics.Stop()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("foo"))
	})

	rw := utils.NewWriterStub()
	req := &http.Request{Header: make(http.Header), URL: &url.URL{Host: "foo.com"}}

	metrics.measureHTTP(handler)(rw, req)
	st.Expect(t, rw.Code, 200)

	rw = utils.NewWriterStub()
	req = &http.Request{Header: make(http.Header), URL: &url.URL{Host: "foo.com"}}

	metrics.measureHTTP(handler)(rw, req)
}

type writerStub struct {
	code int
	data string
}

func (w *writerStub) WriteHeader(code int) {
	w.code = code
}

func (w *writerStub) Write(data []byte) (int, error) {
	w.data = string(data)
	return 0, nil
}
