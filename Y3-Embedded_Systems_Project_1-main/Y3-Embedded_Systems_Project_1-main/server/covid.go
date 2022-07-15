package main

import (
	"encoding/json"
	"fmt"
	mapset "github.com/deckarep/golang-set"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type alertRequest struct {
	UUID     string `json:"uuid"`
	Positive bool   `json:"positive,omitempty"`
}

func alertNeighbours(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var req alertRequest
	err := json.Unmarshal(reqBody, &req)
	if err != nil {
		log.Fatal(err)
	}
	isLookupTable, err := getEnvBool("LOOKUP_TABLE")
	if err != nil {
		log.Fatal(err)
	}

	if req.Positive {
		client := influxdb2.NewClient("http://"+os.Getenv("INFLUXDB_HOSTNAME")+":8086", os.Getenv("INFLUXDB_API_KEY"))
		queryAPI := client.QueryAPI("jbmn")
		if isLookupTable == true {
			var lookUpTable Users
			lookupTableBytes, err := ioutil.ReadFile("lookup_table.yml")
			if err != nil {
				panic(err)
			}
			err = yaml.Unmarshal(lookupTableBytes, &lookUpTable)
			userNameUuidMap := reverseMap(lookUpTable.UuidUserNames)
			sendAlerts(getNeighbours(userNameUuidMap[req.UUID], 2, queryAPI))
		} else {
			sendAlerts(getNeighbours(req.UUID, 2, queryAPI))
		}
		client.Close()
	} else {
		if isLookupTable == true {
			var lookUpTable Users
			lookupTableBytes, err := ioutil.ReadFile("lookup_table.yml")
			if err != nil {
				panic(err)
			}
			err = yaml.Unmarshal(lookupTableBytes, &lookUpTable)
			userNameUuidMap := reverseMap(lookUpTable.UuidUserNames)
			disableAlert(userNameUuidMap[req.UUID])
		} else {
			disableAlert(req.UUID)
		}
	}
	w.WriteHeader(http.StatusOK)
}

func sendAlerts(toAlert mapset.Set) {
	mqttClient := setupClient()
	toAlert.Each(func(neighbour interface{}) bool {
		publish(mqttClient, "IC.es/JBMNsystems/"+fmt.Sprintf("%s", neighbour), "True")
		return false
	})
	mqttClient.Disconnect(250)
}

func disableAlert(uuid string) {
	mqttClient := setupClient()
	publish(mqttClient, "IC.es/JBMNsystems/"+fmt.Sprintf("%s", uuid), "False")
	mqttClient.Disconnect(250)
}
