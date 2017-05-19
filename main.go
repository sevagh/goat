package main

import (
	"github.com/docopt/docopt-go"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var DryRun = false

func main() {
	usage := `kraken - EC2/EBS utility

Usage:
  kraken [-q | --quiet] [-d | --dry]
  kraken -h | --help
  kraken --version

Options:
  -q, --quiet   Suppress output
  -d, --dry     Dry run
  -h --help     Show this screen.
  --version     Show version.`
	arguments, _ := docopt.Parse(usage, nil, true, "kraken 0.1", false)

	currTime := time.Now().UTC()
	log.SetPrefix("KRAKEN: ")
	log.SetFlags(log.Lshortfile)

	if arguments["--quiet"].(bool) {
		log.SetOutput(ioutil.Discard)
	} else {
		log.SetOutput(os.Stderr)
	}

	DryRun = arguments["--dry"].(bool)

	log.Printf("RUNNING KRAKEN: %s", currTime.Format(time.RFC850))

	var ec2Instance Ec2Instance
	var ebsVolumes map[string][]EbsVol
	var err error

	if ec2Instance, err = GetEc2InstanceData(); err != nil {
		log.Fatalf("%v", err)
	}
	if ebsVolumes, err = MapEbsVolumes(&ec2Instance); err != nil {
		log.Fatalf("%v", err)
	}
	if err = AttachEbsVolumes(ec2Instance, ebsVolumes); err != nil {
		log.Fatalf("%v", err)
	}
	for volName, vols := range ebsVolumes {
		log.Printf("Now mounting for volume %s", volName)
		if len(vols) == 1 {
			if err := MountSingleVolume(vols[0]); err != nil {
				log.Fatalf("%v", err)
			}
		} else {
			if err := MountRaidDrives(vols, volName); err != nil {
				log.Fatalf("%v", err)
			}
		}
	}
}
