package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var clientNo *int

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {
	// read flag clientNo
	clientNo = flag.Int("clientNo", 0, "client number")
	flag.Parse()
	fmt.Println(*clientNo)

	// Create mqtt client
	mqtt.DEBUG = log.New(os.Stdout, fmt.Sprintf("[Shino-%d]", *clientNo), 0)
	mqtt.ERROR = log.New(os.Stdout, "[ERROR]", 0)
	opts := mqtt.NewClientOptions().
		AddBroker("tcp://broker.emqx.io:1883").
		SetClientID(fmt.Sprintf("emqx_test_sub_%d", *clientNo))

	opts.SetKeepAlive(60 * time.Second)
	// Set the message callback handler
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Subscribe to a topic
	if token := c.Subscribe("testtopic/random_shino_19472613", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println("errorrrrr", token.Error())
		os.Exit(1)
	}

	// time.Sleep(20 * time.Second)
	// if token := c.Unsubscribe("testtopic/random_lasdkfjlasdj"); token.Wait() && token.Error() != nil {
	// 	fmt.Println(token.Error())
	// 	os.Exit(1)
	// }
	// Disconnect from broker
	chn := make(chan os.Signal)
	signal.Notify(chn, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-chn
		c.Disconnect(250)
		os.Exit(1)
	}()

	for {
		time.Sleep(1 * time.Second)
		fmt.Println("sleeping")
	}
}
