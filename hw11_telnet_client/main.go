package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	flag "github.com/spf13/pflag"
)

var timeoutVal = flag.String("timeout", "10s", "")

func main() {
	flag.Parse()

	errLog := log.New(os.Stderr, "", 0)

	timeout, err := time.ParseDuration(*timeoutVal)
	if err != nil {
		log.Fatalf("can't to parse flag: %v", err)
	}

	c := len(os.Args)
	if c < 3 {
		log.Fatalf("invalid input arguments")
	}

	args := os.Args[c-2 : c]
	address := net.JoinHostPort(args[0], args[1])

	ctx, cancel := context.WithCancel(context.Background())

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout, os.Stderr, cancel)

	err = client.Connect()
	if err != nil {
		log.Fatalf("can't connect: %v", err)
	}
	defer client.Close()

	go func() {
		err := client.Receive()
		if err != nil {
			errLog.Printf("can't receieve: %v", err)
			return
		}
	}()

	go func() {
		err := client.Send()
		if err != nil {
			errLog.Printf("can't send: %v", err)
			return
		}
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	select {
	case <-sigint:
		cancel()
	case <-ctx.Done():
		close(sigint)
	}
}
