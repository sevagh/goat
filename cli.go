package main

import (
	"github.com/docopt/docopt-go"
	"log"
	"os"
	"time"
)

func main() {
	usage := `kraken - EC2/EBS utility

Usage:
  kraken <mount_path>
  kraken -h | --help
  kraken --version

Options:
  -h --help     Show this screen.
  --version     Show version.`
	arguments, _ := docopt.Parse(usage, nil, true, "Kraken 0.1", false)

	currTime := time.Now().UTC()
	logger := log.New(os.Stderr, "kraken: ", log.Lshortfile)
	logger.Printf("RUNNING KRAKEN: %s", currTime.Format(time.RFC850))
	_, ok := arguments["<mount_path>"].(string)

	if ok {
		deviceNames, err := AttachEbsVolumes(logger)
		if err != nil {
			logger.Println(err)
			os.Exit(-1)
		}
		logger.Printf("Attached: %s\n", deviceNames)
	}
}
