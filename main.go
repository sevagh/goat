package main

import (
	"log"
	"os"
	"time"
)

func main() {
	currTime := time.Now().UTC()
	logger := log.New(os.Stderr, "kraken: ", log.Lshortfile)
	logger.Printf("RUNNING KRAKEN: %s", currTime.Format(time.RFC850))

	ec2Instance, err := GetEc2InstanceData(logger)
	if err != nil {
		logger.Fatalf("%v", err)
		os.Exit(-1)
	}
	ebsVolumes, err := FindEbsVolumes(&ec2Instance, logger)
	if err != nil {
		logger.Fatalf("%v", err)
		os.Exit(-1)
	}
	attachedVolumes, err := AttachEbsVolumes(ec2Instance, ebsVolumes, logger)
	if err != nil {
		logger.Fatalf("%v", err)
		os.Exit(-1)
	}

	for volId, vols := range attachedVolumes {
		logger.Printf("Now mounting for volume %d", volId)
		if len(vols) == 1 {
			if err := MountSingleDrive(vols[0], logger); err != nil {
				os.Exit(-1)
			}
		} else {
			if err := MountRaidDrives(vols, volId, logger); err != nil {
				os.Exit(-1)
			}
		}
	}
}
