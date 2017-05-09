package main

import (
	"os"
	"log"
	"time"
)

func main() {
	currTime := time.Now().UTC()
	logger := log.New(os.Stderr, "kraken: ", log.Lshortfile)
	logger.Printf("RUNNING KRAKEN: %s", currTime.Format(time.RFC850))
	deviceNames, err := AttachEbsVolumes(logger)
	if err != nil {
		logger.Println(err)
		os.Exit(-1)
	}
	logger.Printf("Attached: %s\n", deviceNames)
}
