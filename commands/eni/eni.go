package eni

import (
	log "github.com/sirupsen/logrus"

	"github.com/sevagh/goat/awsutil"
)

//GoatEni runs Goat for your ENIs - attach
func GoatEni(ec2Instance awsutil.EC2Instance, dryRun bool, debug bool) {
	log.Printf("2: COLLECTING ENI INFO")
	enis := awsutil.FindEnis(&ec2Instance)

	log.Printf("3: ATTACHING ENIS")
	awsutil.AttachEnis(ec2Instance, enis, dryRun)
}
