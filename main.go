package main

import (
	"github.com/docopt/docopt-go"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var DryRun = false
var PREFIX = "KRKN"

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
	log.SetPrefix(PREFIX + ": ")
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

	log.Printf("%s", DrawAsciiBanner("1: COLLECTING EC2 INFO"))
	if ec2Instance, err = GetEc2InstanceData(); err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("%s", DrawAsciiBanner("2: COLLECTING EBS INFO"))
	if ebsVolumes, err = MapEbsVolumes(&ec2Instance); err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("%s", DrawAsciiBanner("3: ATTACHING EBS VOLS"))
	if ebsVolumes, err = AttachEbsVolumes(ec2Instance, ebsVolumes); err != nil {
		log.Fatalf("%v", err)
	}

	log.Printf("%s", DrawAsciiBanner("4: MOUNTING ATTACHED VOLS"))
	for volName, vols := range ebsVolumes {
		if len(vols) == 1 {
			if err := MountSingleDrive(vols[0].AttachedName, vols[0].MountPath, vols[0].FsType, vols[0].VolumeName); err != nil {
				log.Fatalf("%v", err)
			}
		} else {
			if err := MountRaidDrives(vols, volName); err != nil {
				log.Fatalf("%v", err)
			}
		}
	}
}
