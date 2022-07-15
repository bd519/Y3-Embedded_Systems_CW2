package main

import (
	"context"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"log"
	"os"
)

func continuousMonitor() {
	client := influxdb2.NewClient("http://"+os.Getenv("INFLUXDB_HOSTNAME")+":8086", os.Getenv("INFLUXDB_API_KEY"))
	queryAPI := client.QueryAPI("jbmn")
	window, _ := getEnvUint("TEMP_WINDOW")
	threshold, _ := getEnvUint("TEMP_THRESHOLD")
	query := fmt.Sprintf(`from(bucket:"iot")
									|> range(start: -%dm)
									|> filter(fn: (r) => r._measurement == "temperature")
									|> filter(fn : (r) => r._field == "temperature")
									|> filter(fn : (r) => r._value >= %d)`, window, threshold)
	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	set := mapset.NewSet()
	for result.Next() {
		set.Add(result.Record().Values()["host"])
	}
	sendAlerts(set)
}
