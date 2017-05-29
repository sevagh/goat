package main

import (
	"fmt"
	"github.com/docopt/docopt-go"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

//PREFIX defines the prefix to use for all tags and labels
var PREFIX = "EWIZ"

func main() {
	usage := `ebs_raid_wizard - EC2/EBS utility

Usage:
  ebs_raid_wizard [--log-level=<log-level>] [--dry]
  ebs_raid_wizard -h | --help
  ebs_raid_wizard --version

Options:
  --log-level=<level>  Log level (debug, info, warn, error, fatal) [default: info]
  --dry                Dry run
  -h --help            Show this screen.
  --version            Show version.`
	arguments, _ := docopt.Parse(usage, nil, true, "ebs_raid_wizard 0.1", false)

	log.SetOutput(os.Stderr)
	logLevel := arguments["--log-level"].(string)
	if level, err := log.ParseLevel(logLevel); err != nil {
		log.Fatalf("%v", err)
	} else {
		log.SetLevel(level)
	}

	log.SetFormatter(&log.TextFormatter{})

	dryRun := arguments["--dry"].(bool)

	log.Printf("%s", drawASCIIBanner("WELCOME TO EBS_RAID_WIZARD"))

	log.Printf("%s", drawASCIIBanner("1: COLLECTING EC2 INFO"))
	ec2Instance := GetEC2InstanceData()

	log.Printf("%s", drawASCIIBanner("2: COLLECTING EBS INFO"))
	ebsVolumes := MapEbsVolumes(&ec2Instance)

	log.Printf("%s", drawASCIIBanner("3: ATTACHING EBS VOLS"))
	ebsVolumes = AttachEbsVolumes(ec2Instance, ebsVolumes, dryRun)

	log.Printf("%s", drawASCIIBanner("4: MOUNTING ATTACHED VOLS"))

	if len(ebsVolumes) == 0 {
		log.Warn("Empty vols, nothing to do")
		os.Exit(0)
	}

	for volName, vols := range ebsVolumes {
		PrepAndMountDrives(volName, vols, dryRun)
	}
}

func drawASCIIBanner(headLine string) string {
	return fmt.Sprintf("\n%[1]s\n# %[2]s #\n%[1]s\n",
		strings.Repeat("#", len(headLine)+4),
		headLine)
}
