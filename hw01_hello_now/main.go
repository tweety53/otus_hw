package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	const timeLayout = "2006-01-02 15:04:05 -0700 MST"

	ntpTime, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		log.Fatalf("ntp.Time() error: %s", err)
	}

	fmt.Printf("current time: %s\n", time.Now().Format(timeLayout))
	fmt.Printf("exact time: %s\n", ntpTime.Format(timeLayout))
}
