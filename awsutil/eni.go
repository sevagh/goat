package awsutil

import (
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

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
