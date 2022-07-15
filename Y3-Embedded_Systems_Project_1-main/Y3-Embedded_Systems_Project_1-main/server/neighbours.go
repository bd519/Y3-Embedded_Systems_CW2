package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/deckarep/golang-set"
	"github.com/goccy/go-graphviz"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type graphRequest struct {
	UUID   string `json:"uuid"`
	Window uint   `json:"window,string"`
}

type Users struct {
	UuidUserNames map[string]string `yaml:"users"`
}

func returnUserNeighbours(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	reqBody, _ := ioutil.ReadAll(r.Body)
	var req graphRequest
	err := json.Unmarshal(reqBody, &req)
	if err != nil {
		log.Fatal(err)
	}
	client := influxdb2.NewClient("http://"+os.Getenv("INFLUXDB_HOSTNAME")+":8086", os.Getenv("INFLUXDB_API_KEY"))
	queryAPI := client.QueryAPI("jbmn")

	isLookupTable, err := getEnvBool("LOOKUP_TABLE")
	if err != nil {
		log.Fatal(err)
	}
	var resp bytes.Buffer
	if isLookupTable == true {
		var lookUpTable Users
		lookupTableBytes, err := ioutil.ReadFile("lookup_table.yml")
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(lookupTableBytes, &lookUpTable)

		userNameUuidMap := reverseMap(lookUpTable.UuidUserNames)

		if err != nil {
			return
		}
		resp = generateGraph(userNameUuidMap[req.UUID], getNeighbours(userNameUuidMap[req.UUID], req.Window, queryAPI), lookUpTable.UuidUserNames)
	} else {
		resp = generateGraph(req.UUID, getNeighbours(req.UUID, req.Window, queryAPI), map[string]string{})
	}

	client.Close()
	_, err = w.Write(resp.Bytes())
	if err != nil {
		log.Fatal(err)
	}
}

func getNeighbours(uuid string, timeRange uint, queryAPI api.QueryAPI) mapset.Set {

	//Forward Lookup
	queryA := fmt.Sprintf(`from(bucket:"iot")
									|> range(start: -%dd)
									|> filter(fn: (r) => r._measurement == "mem")
									|> filter(fn: (r) => r["host"] == "%s")
									|> filter(fn : (r) => r._field == "RSSI")
									|> filter(fn : (r) => r._value >= -70.0)`, timeRange, uuid)

	//Backward Lookup
	queryB := fmt.Sprintf(`from(bucket:"iot")
									|> range(start: -%dd)
									|> filter(fn: (r) => r._measurement == "mem")
									|> filter(fn: (r) => r["neighbour"] == "%s")
									|> filter(fn : (r) => r._field == "RSSI")
									|> filter(fn : (r) => r._value >= -70.0)`, timeRange, uuid)

	resultA, err := queryAPI.Query(context.Background(), queryA)
	if err != nil {
		panic(err)
	}
	resultB, err := queryAPI.Query(context.Background(), queryB)
	if err != nil {
		panic(err)
	}
	set := mapset.NewSet()

	// Iterate over queryA response
	for resultA.Next() {
		set.Add(resultA.Record().Values()["neighbour"])
	}

	for resultB.Next() {
		set.Add(resultB.Record().Values()["host"])
	}

	return set
}

func generateGraph(uuid string, s mapset.Set, lookupTable map[string]string) bytes.Buffer {

	newset := mapset.NewSet()

	s.Each(func(neighbourUuid interface{}) bool {
		if value, ok := lookupTable[fmt.Sprintf("%s", neighbourUuid)]; ok {
			newset.Add(value)
		} else {
			newset.Add(neighbourUuid)
		}
		return false
	})

	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := graph.Close(); err != nil {
			log.Fatal(err)
		}
		err := g.Close()
		if err != nil {
			log.Fatal(err)
			return
		}
	}()
	self, err := graph.CreateNode(lookupTable[uuid])
	if err != nil {
		log.Fatal(err)
	}
	newset.Each(func(neighbour interface{}) bool {
		tmp, err := graph.CreateNode(fmt.Sprintf("%s", neighbour))
		if err != nil {
			log.Fatal(err)
		}
		_, err = graph.CreateEdge("", self, tmp)
		if err != nil {
			log.Fatal(err)
		}
		return false
	})
	var buf bytes.Buffer
	if err := g.Render(graph, "svg", &buf); err != nil {
		log.Fatal(err)
	}
	return buf
}
