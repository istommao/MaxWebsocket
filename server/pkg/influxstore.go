package pkg

import (
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type InfluxStore struct {
	client   influxdb2.Client
	writeAPI api.WriteAPI
}

func NewInfluxDB() *InfluxStore {
	// Create a new client using an InfluxDB server base URL and an authentication token
	client := influxdb2.NewClient("http://localhost:8086", "nkzlXjiT87wDWj0CuzcswWmELyaTMZbrjxDarGai3dh3bsseOEflIOkrPx3Rp4YCj4tfOFfEwFowbIj6h0QU9A==")
	writeAPI := client.WriteAPI("CodeMax", "CodeMax")

	return &InfluxStore{client: client, writeAPI: writeAPI}
}

func (w *InfluxStore) SaveMessage(key string, value []byte) {
	// Use blocking write client for writes to desired bucket

	//  WriteAPIBlocking("my-org", "my-bucket")
	// Create point using full params constructor
	p := influxdb2.NewPoint("stat",
		map[string]string{"unit": "temperature", "key": key},
		map[string]interface{}{"avg": 24.5, "max": 45.0, "value": value},
		time.Now())
	// write point immediately

	w.writeAPI.WritePoint(p)
}

func (w *InfluxStore) Close() {
	w.client.Close()
}
