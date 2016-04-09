package influx

import (
	"fmt"
	"log"
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

type reporter struct {
	config Config
	client client.Client
}

// New creates a new InfluxDB reporter which will post the metrics to the specified server.
func New(c Config) *reporter {
	re := &reporter{config: c}
	if err := re.makeClient(); err != nil {
		log.Printf("unable to make InfluxDB client. err=%v", err)
		return nil
	}
	return re
}

func (r *reporter) makeClient() (err error) {
	r.client, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     r.config.URL,
		Username: r.config.Username,
		Password: r.config.Password,
	})
	return
}

func (r *reporter) Report(re metrics.Report) {
	if err := r.send(re); err != nil {
		log.Printf("error sending metrics to InfluxDB: err=%v", err)
	}
}

func (r *reporter) send(re metrics.Report) error {
	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Precision: "s",
		Database:  r.config.Database,
	})
	if err != nil {
		return err
	}

	// TODO: aggregate histogram as unique point

	now := time.Now()
	for key, value := range re.Counters {
		fields := map[string]interface{}{"value": int64(value)}
		pt, err := client.NewPoint(fmt.Sprintf("%s.count", key), r.config.Tags, fields, now)
		if err != nil {
			return err
		}
		bp.AddPoint(pt)
	}

	for key, value := range re.Gauges {
		fields := map[string]interface{}{"value": int64(value)}
		pt, err := client.NewPoint(fmt.Sprintf("%s.gauge", key), r.config.Tags, fields, now)
		if err != nil {
			return err
		}
		bp.AddPoint(pt)
	}

	// Send data to InfluxDB server
	return r.client.Write(bp)
}
