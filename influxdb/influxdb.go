package influxdb

import (
	"context"
	"time"

	"github.com/Scrin/scanmytesla-influxdb/config"
	influxdb "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/rs/zerolog/log"
)

type Value struct {
	Field string
	Value float64
}

var client influxdb.Client
var writeAPI api.WriteAPIBlocking
var measurement string

func Setup(conf config.InfluxDB) {
	url := conf.Url
	if url == "" {
		url = "https://localhost:8086"
	}
	bucket := conf.Bucket
	if bucket == "" {
		bucket = "tesla"
	}
	measurement = conf.Measurement
	if measurement == "" {
		measurement = "tesla"
	}

	client = influxdb.NewClient(url, conf.AuthToken)
	writeAPI = client.WriteAPIBlocking(conf.Org, bucket)
	writeAPI.EnableBatching()
	log.Info().Str("url", url).Str("bucket", bucket).Str("measurement", measurement).Msg("Setup InfluxDB client")
}

func Close() {
	writeAPI.Flush(context.Background())
	client.Close()
}

func Send(timestamp time.Time, values []Value) error {
	p := influxdb.NewPointWithMeasurement(measurement)
	p.SetTime(timestamp)
	for _, v := range values {
		p.AddField(v.Field, v.Value)
	}
	return writeAPI.WritePoint(context.Background(), p)
}
