package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

const host = "0.beevik-ntp.pool.ntp.org" // этот параметр для удобства можно вынести в качестве флага

func main() {
	ntpTime, err := ntp.Time(host)
	if err != nil {
		log.Fatalf("can't get time from remote host %s: %v", host, err)
	}
	fmt.Printf("current time: %v\n", time.Now().Round(0))
	fmt.Printf("exact time: %v\n", ntpTime.Round(0))
}
