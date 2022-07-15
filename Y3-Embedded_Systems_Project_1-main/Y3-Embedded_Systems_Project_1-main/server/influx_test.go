package main

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"math/rand"
	"testing"
	"time"
)

var hostnames = []string{
	"11111111-2222-3333-1012-555555555555",
	"11111111-2222-3333-1355-555555555555",
	"11111111-2222-3333-1549-555555555555",
}

func generatePoints() []*write.Point {
	points := make([]*write.Point, 1000)
	for i := 0; i < 1000; i++ {
		points[i] = influxdb2.NewPoint("mem",
			map[string]string{"host": hostnames[i%3], "neighbour": hostnames[(i+1)%3]},
			map[string]interface{}{"RSSI": -rand.Intn(60)},
			time.Now().Add(time.Duration(-rand.Intn(30))*time.Minute))
	}
	return points
}

func TestInflux(t *testing.T) {
	// Create a client
	// You can generate an API Token from the "API Tokens Tab" in the UI
	client := influxdb2.NewClient("http://localhost:8086", "2zJivHhc53i9RtFelszigZv9_W9kUabRpoCSJgwP-EvtoWznQwaovALqg8uC6cpBRGosqFrDDuA9omccfk6HJA==")
	// get non-blocking write client
	writeAPI := client.WriteAPI("jbmn", "iot")

	p := generatePoints()
	// write point asynchronously
	for _, point := range p {
		writeAPI.WritePoint(point)
	}
	// Flush writes
	writeAPI.Flush()
	defer client.Close()
}

func TestClearInflux(t *testing.T) {
	// Create a client
	// You can generate an API Token from the "API Tokens Tab" in the UI
	client := influxdb2.NewClient("http://localhost:8086", "4Hn1YgPWnb64zYjjqCsKD19If-vuif4UfX80qUvruvbj3CkrsekzAoqLMtrP-on9yU6RW3LzsdoDbOO7lVQ1IA==")
	// get non-blocking write client
	deleteAPI := client.DeleteAPI()
	_ = deleteAPI.DeleteWithName(context.Background(), "jbmn", "iot", time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Now(), "")
	defer client.Close()
}
