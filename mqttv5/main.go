package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eclipse/paho.golang/paho"
)

var qos *int

func main() {
	// read flag clientNo
	qos = flag.Int("qos", 0, "QoS")
	flag.Parse()
	fmt.Println("Running test with QoS: ", *qos)
	logger := log.New(os.Stdout, "PUBLISHER: ", log.LstdFlags)
	ctx := context.Background()

	conn, err := net.Dial("tcp", "broker.emqx.io:1883")
	if err != nil {
		err := fmt.Errorf("failed to connect to %s: %w", "broker.emqx.io:1883", err)
		panic(err)
	}

	cfg := paho.ClientConfig{
		Conn: conn,
	}

	c := paho.NewClient(cfg)
	cp := &paho.Connect{
		KeepAlive:  30,
		ClientID:   "emqx_test_pub",
		CleanStart: true,
		Username:   "hungtd",
		Password:   []byte("hungtd"),
	}
	c.SetDebugLogger(logger)
	c.SetErrorLogger(logger)

	ca, err := c.Connect(ctx, cp)
	if err != nil {
		fmt.Println(err)
	}
	if ca.ReasonCode != 0 {
		log.Fatalf("Failed to connect to %s : %d - %s", "broker.emqx.io:1883", ca.ReasonCode, ca.Properties.ReasonString)
	}

	// Handle Disconnect
	chn := make(chan os.Signal)
	signal.Notify(chn, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-chn
		fmt.Println("signal received, exiting")
		if c != nil {
			d := &paho.Disconnect{ReasonCode: 0}
			if err := c.Disconnect(d); err != nil {
				fmt.Println(err)
			}
		}
		os.Exit(1)
	}()

	count := 0
	for {
		fmt.Println("------------------")
		msg := fmt.Sprintf("Hello World %d", count)
		mqttPub := paho.Publish{
			Payload: []byte(msg),
			Topic:   "testtopic/random_shino_19472613",
			QoS:     1,
			Retain:  false,
		}

		_, err := c.Publish(ctx, &mqttPub)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Published message: ", msg)
		count++
		time.Sleep(5 * time.Second)
	}
}
