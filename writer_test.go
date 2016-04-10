package metrics

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/nbio/st"
	"gopkg.in/vinxi/utils.v0"
)

func TestMetricWriter(t *testing.T) {
	w := utils.NewWriterStub()
	req := &http.Request{URL: &url.URL{}, Header: make(http.Header)}

	var info *Info
	collector := func(i *Info) {
		info = i
	}

	writer := newMetricsWriter(w, req, collector)
	writer.Header().Set("foo", "bar")
	writer.Header().Set("Content-Length", "11")
	writer.WriteHeader(200) // collect
	writer.Write([]byte("hello world"))

	st.Expect(t, w.Code, 200)
	st.Expect(t, w.Header().Get("foo"), "bar")
	st.Expect(t, string(w.Body), "hello world")

	st.Expect(t, info.Status, w.Code)
	st.Expect(t, info.Header.Get("foo"), "bar")
	st.Expect(t, info.BodyLength, int64(11))
	st.Expect(t, info.Request, req)
}
