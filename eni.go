package main

import (
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

//GoatEni runs Goat for your ENIs - attach
func GoatEni(dryRun bool, debug bool) {
	log.Printf("WELCOME TO GOAT")
	log.Printf("1: COLLECTING EC2 INFO")
	ec2Instance := GetEC2InstanceData()

	log.Printf("2: COLLECTING ENI INFO")
	ec2Instance.FindEnis()

	log.Printf("3: ATTACHING ENIS")

	if len(ec2Instance.Enis) == 0 {
		log.Warn("Empty enis, nothing to do")
		os.Exit(0)
	}

	ec2Instance.AttachEnis(dryRun)
}

//FindEnis returns a list of all ENIs that should be attached to this EC2 instance
func (e *EC2Instance) FindEnis() {
	log.Info("Searching for ENIs")

	params := &ec2.DescribeNetworkInterfacesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:GOAT-IN:Prefix"),
				Values: []*string{
					aws.String(e.Prefix),
				},
			},
			{
				Name: aws.String("tag:GOAT-IN:NodeId"),
				Values: []*string{
					aws.String(e.NodeID),
				},
			},
			{
				Name: aws.String("availability-zone"),
				Values: []*string{
					aws.String(e.Az),
				},
			},
		},
	}

	enis := []string{}

	result, err := e.EC2Client.DescribeNetworkInterfaces(params)
	if err != nil {
		log.Fatalf("Error when searching for ENIs: %v", err)
	}

	for _, eni := range result.NetworkInterfaces {
		attachedID := ""
		if eni.Attachment != nil {
			attachedID = *eni.Attachment.InstanceId
		}
		if attachedID != "" {
			if attachedID != e.InstanceID {
				log.Fatalf("Eni %s attached to different instance-id: %s", *eni.NetworkInterfaceId, attachedID)
			} else {
				log.Infof("Eni %s already attached to this instance, skipping", *eni.NetworkInterfaceId)
				continue
			}
		}
		enis = append(enis, *eni.NetworkInterfaceId)
	}

	e.Enis = enis
}
