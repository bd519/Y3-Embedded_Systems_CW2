package main

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"os"
)

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func setupClient() mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", os.Getenv("BROKER_HOSTNAME"), 1883))
	opts.SetClientID("go_alerting_service")
	//	opts.SetUsername("emqx")
	// opts.SetPassword("public")
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return mqttClient
}
func publish(client mqtt.Client, topic string, msg string) {
	token := client.Publish(topic, 0, false, msg)
	fmt.Printf("Sending %s\n", msg)
	token.Wait()
}
