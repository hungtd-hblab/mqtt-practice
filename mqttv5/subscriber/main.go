package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/eclipse/paho.golang/paho"
)

var clientNo *int

func main() {
	// read flag clientNo
	clientNo = flag.Int("clientNo", 0, "client number")
	flag.Parse()
	fmt.Println(*clientNo)

	logger := log.New(os.Stdout, fmt.Sprintf("SUBSCRIBER %d: ", *clientNo), log.LstdFlags)

	conn, err := net.Dial("tcp", "broker.emqx.io:1883")
	if err != nil {
		err := fmt.Errorf("failed to connect to %s: %w", "broker.emqx.io:1883", err)
		panic(err)
	}

	msgChan := make(chan *paho.Publish)

	cfg := paho.ClientConfig{
		Router: paho.NewSingleHandlerRouter(func(m *paho.Publish) {
			msgChan <- m
		}),
		Conn: conn,
	}

	c := paho.NewClient(cfg)
	cp := &paho.Connect{
		KeepAlive:  30,
		ClientID:   fmt.Sprintf("emqx_test_sub_%d", *clientNo),
		CleanStart: true,
	}
	c.SetDebugLogger(logger)
	c.SetErrorLogger(logger)

	ca, err := c.Connect(context.Background(), cp)
	if err != nil {
		fmt.Println(err)
	}
	if ca.ReasonCode != 0 {
		log.Fatalf("Failed to connect to %s : %d - %s", "broker.emqx.io:1883", ca.ReasonCode, ca.Properties.ReasonString)
	}
	sa, err := c.Subscribe(context.Background(), &paho.Subscribe{
		Subscriptions: map[string]paho.SubscribeOptions{
			"testtopic/random_shino_19472613": {QoS: byte(1)},
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	if sa.Reasons[0] != byte(1) {
		log.Fatalf("Failed to subscribe to %s : %d", "testtopic/random_shino_19472613", sa.Reasons[0])
	}
	log.Printf("Subscribed to %s", "testtopic/random_shino_19472613")

	for m := range msgChan {
		fmt.Println("--------------------")
		log.Println("Received message:", string(m.Payload))
	}
}
