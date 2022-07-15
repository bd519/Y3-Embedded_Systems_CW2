package main

import (
	"github.com/gorilla/mux"
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"os"
)

func main() {

	c := cron.New()
	_, err := c.AddFunc("@every 15m", continuousMonitor)
	if err != nil {
		log.Fatal(err)
	}
	c.Start()
	handleRequests()
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/neighbours", returnUserNeighbours).Methods("POST")
	myRouter.HandleFunc("/test", alertNeighbours).Methods("POST")
	httpPort := os.Getenv("HTTP_PORT")

	log.Fatal(http.ListenAndServe(":"+httpPort, myRouter))
}
