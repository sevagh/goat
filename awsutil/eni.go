package awsutil

import (
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

//FindEnis returns a list of all ENIs that should be attached to this EC2 instance
func FindEnis(ec2Instance *EC2Instance) []string {
	log.Info("Searching for ENIs")

	params := &ec2.DescribeNetworkInterfacesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:GOAT-IN:Prefix"),
				Values: []*string{
					aws.String(ec2Instance.Prefix),
				},
			},
			{
				Name: aws.String("tag:GOAT-IN:NodeId"),
				Values: []*string{
					aws.String(ec2Instance.NodeID),
				},
			},
			{
				Name: aws.String("availability-zone"),
				Values: []*string{
					aws.String(ec2Instance.Az),
				},
			},
		},
	}

	enis := []string{}

	result, err := ec2Instance.EC2Client.DescribeNetworkInterfaces(params)
	if err != nil {
		log.Fatalf("Error when searching for ENIs: %v", err)
	}

	for _, eni := range result.NetworkInterfaces {
		attachedID := *eni.Attachment.InstanceId
		if attachedID != "" {
			if attachedID != ec2Instance.InstanceID {
				log.Fatalf("Eni %s attached to different instance-id: %s", *eni.NetworkInterfaceId, attachedID)
			}
		} else {
			continue
		}
		enis = append(enis, *eni.NetworkInterfaceId)
	}

	return enis
}
