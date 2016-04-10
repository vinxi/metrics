package influx

import (
	"fmt"
	"log"
	"strings"
	"time"

	client "github.com/influxdata/influxdb/client/v2"
	"gopkg.in/vinxi/metrics.v0"
)

// Config stores the InfluxDB connection params.
type Config struct {
	URL      string
	Database string
	Username string
	Password string
	Tags     map[string]string
}

// Reporter implements an InfluxDB metrics reporter who send data to a InfluxDB server via HTTP.
type Reporter struct {
	config Config
	client client.Client
}

// New creates a new InfluxDB reporter which will post the metrics to the specified server.
func New(c Config) *Reporter {
	re := &reporter{config: c}
	if err := re.makeClient(); err != nil {
		log.Printf("unable to make InfluxDB client. err=%v", err)
	}
	return re
}

// Report implements the metrisc.Reporter interface.
func (r *Reporter) Report(re metrics.Report) {
	if err := r.send(re); err != nil {
		log.Printf("infludb: error sending metrics err=%v", err)
	}
}

func (r *Reporter) makeClient() (err error) {
	r.client, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     r.config.URL,
		Username: r.config.Username,
		Password: r.config.Password,
	})
	return
}

func (r *Reporter) send(re metrics.Report) error {
	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Precision: "s",
		Database:  r.config.Database,
	})
	if err != nil {
		return err
	}

	// Map metrics report into influx points
	if err := r.mapReport(re, bp); err != nil {
		return err
	}

	// Send data to InfluxDB server
	return r.client.Write(bp)
}

func (r *Reporter) mapReport(re metrics.Report, bp client.BatchPoints) error {
	now := time.Now()

	// Extract histograms properly
	gauges, histograms := extractHistograms(re.Gauges)

	// Add histograms
	for key, hg := range histograms {
		fields := map[string]interface{}{
			"p50":  hg["P50"],
			"p75":  hg["P75"],
			"p90":  hg["P90"],
			"p95":  hg["P95"],
			"p99":  hg["P99"],
			"p999": hg["P999"],
		}
		pt, err := client.NewPoint(fmt.Sprintf("%s.histogram", key), r.config.Tags, fields, now)
		if err != nil {
			return err
		}
		bp.AddPoint(pt)
	}

	// Add gauges
	for key, value := range gauges {
		fields := map[string]interface{}{"value": int64(value)}
		pt, err := client.NewPoint(fmt.Sprintf("%s.gauge", key), r.config.Tags, fields, now)
		if err != nil {
			return err
		}
		bp.AddPoint(pt)
	}

	// Add counters
	for key, value := range re.Counters {
		fields := map[string]interface{}{"value": int64(value)}
		pt, err := client.NewPoint(fmt.Sprintf("%s.count", key), r.config.Tags, fields, now)
		if err != nil {
			return err
		}
		bp.AddPoint(pt)
	}

	return nil
}

// extractHistograms is used to split standalone gauge metrics from histograms.
// This function could be generalized in the future.
func extractHistograms(records map[string]int64) (map[string]int64, map[string]map[string]int64) {
	gauges := make(map[string]int64)
	histograms := make(map[string]map[string]int64)

	for key, value := range records {
		parts := strings.Split(key, ".")
		perc := parts[len(parts)-1]

		// If not percentile, store as unique gauge
		if !isPercentile(perc) {
			gauges[key] = value
			continue
		}

		// Aggregate histogram percentiles
		name := strings.Join(parts[0:len(parts)-1], ".")
		store, ok := histograms[name]
		if !ok {
			store = make(map[string]int64)
			histograms[name] = store
		}
		store[perc] = value
	}

	return gauges, histograms
}

func isPercentile(key string) bool {
	return key == "P50" || key == "P75" || key == "P90" ||
		key == "P95" || key == "P99" || key == "P999"
}
