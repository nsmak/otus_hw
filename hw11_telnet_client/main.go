package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

var timeoutVal = flag.String("timeout", "10s", "")

func main() {
	flag.Parse()

	errLog := log.New(os.Stderr, "", 0)

	timeout, err := time.ParseDuration(*timeoutVal)
	if err != nil {
		log.Fatalf("can't to parse flag: %v", err)
	}

	if len(os.Args) < 4 {
		log.Fatalf("invalid input arguments")
	}

	args := os.Args[2:]
	if len(args) < 2 {
		log.Fatalln("invalid incoming address arguments")
	}

	address := net.JoinHostPort(args[0], args[1])

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	err = client.Connect()
	if err != nil {
		log.Fatalf("can't connect: %v", err)
	}

	go func() {
		defer client.Close()

		err := client.Receive()
		if err != nil {
			errLog.Printf("can't receieve: %v", err)
			return
		}
		log.Println("...Connection was closed by peer")
	}()

	go func() {
		defer client.Close()

		err := client.Send()
		if err != nil {
			errLog.Printf("can't send: %v", err)
			return
		}
		fmt.Println("...EOF")
	}()

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		client.Close()
	}()

	<-client.Done()
}
