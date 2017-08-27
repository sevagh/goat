package eni

import (
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/sevagh/goat/awsutil"
)

//GoatEni runs Goat for your ENIs - attach
func GoatEni(dryRun bool, debug bool) {
	log.Printf("WELCOME TO GOAT")
	log.Printf("1: COLLECTING EC2 INFO")
	ec2Instance := awsutil.GetEC2InstanceData()

	log.Printf("2: COLLECTING ENI INFO")
	ec2Instance.FindEnis()

	log.Printf("3: ATTACHING ENIS")

	if len(ec2Instance.Enis) == 0 {
		log.Warn("Empty enis, nothing to do")
		os.Exit(0)
	}

	ec2Instance.AttachEnis(dryRun)
}
