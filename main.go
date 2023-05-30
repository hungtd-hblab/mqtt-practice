package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {
	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().
		AddBroker("tcp://broker.emqx.io:1883").
		SetClientID("emqx_test_pub")

	opts.SetKeepAlive(60 * time.Second)
	// Set the message callback handler
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Handle Disconnect
	chn := make(chan os.Signal)
	signal.Notify(chn, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-chn
		c.Disconnect(250)
		os.Exit(1)
	}()

	count := 0
	for {
		time.Sleep(10 * time.Second)
		fmt.Println("------------------")
		// Publish a message
		token := c.Publish(
			"testtopic/random_shino_19472613",
			0,
			false,
			fmt.Sprintf("Hello World %d", count),
		)
		token.Wait()
		count++
	}
}
