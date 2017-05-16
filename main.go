package main

import (
	"log"
	"os"
	"io/ioutil"
	"time"
	"github.com/docopt/docopt-go"
)

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

	dryRun := arguments["--dry"].(bool)
	DryRun = dryRun

	log.Printf("RUNNING KRAKEN: %s", currTime.Format(time.RFC850))

	var ec2Instance Ec2Instance
	var ebsVolumes []EbsVol
	var attachedVolumes map[int][]EbsVol
	var err error

	if ec2Instance, err = GetEc2InstanceData(); err != nil {
		log.Fatalf("%v", err)
	}
	if ebsVolumes, err = FindEbsVolumes(&ec2Instance); err != nil {
		log.Fatalf("%v", err)
	}
	if attachedVolumes, err = AttachEbsVolumes(ec2Instance, ebsVolumes, dryRun); err != nil {
		log.Fatalf("%v", err)
	}
	for volId, vols := range attachedVolumes {
		log.Printf("Now mounting for volume %d", volId)
		if len(vols) == 1 {
			if err := MountSingleVolume(vols[0]); err != nil {
				log.Fatalf("%v", err)
			}
		} else {
			if err := MountRaidDrives(vols, volId); err != nil {
				log.Fatalf("%v", err)
			}
		}
	}
}
