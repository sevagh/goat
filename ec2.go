package main

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Ec2Instance struct {
	Ec2Client  *ec2.EC2
	InstanceId string
	Prefix     string
	NodeId     string
}

func GetEc2InstanceData() (Ec2Instance, error) {
	var ec2Instance Ec2Instance
	sess := session.New()

	creds := credentials.NewCredentials(
		&ec2rolecreds.EC2RoleProvider{
			Client:       ec2metadata.New(sess),
			ExpiryWindow: 5 * time.Minute,
		},
	)
	sess.Config.Credentials = creds

	svc := ec2metadata.New(sess)
	result, err := svc.GetMetadata("instance-id")
	if err != nil {
		return ec2Instance, err
	}

	region, err := getInstanceRegion(svc)
	if err != nil {
		return ec2Instance, err
	}

	ec2Instance.InstanceId = result
	sess.Config.Region = &region
	ec2Svc := ec2.New(sess)
	ec2Instance.Ec2Client = ec2Svc

	err = getInstanceTags(&ec2Instance)
	if err != nil {
		return ec2Instance, err
	}

	return ec2Instance, nil
}

func getInstanceRegion(svc *ec2metadata.EC2Metadata) (string, error) {
	id, err := svc.GetInstanceIdentityDocument()
	if err != nil {
		return "", err
	}
	return id.Region, nil
}

func getInstanceTags(ec2Instance *Ec2Instance) error {
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("instance-id"),
				Values: []*string{
					aws.String(ec2Instance.InstanceId),
				},
			},
		},
	}

	result, err := ec2Instance.Ec2Client.DescribeInstances(params)
	if err != nil {
		return err
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			for _, tag := range instance.Tags {
				if *tag.Key == "NodeId" {
					ec2Instance.NodeId = *tag.Value
				} else if *tag.Key == "Prefix" {
					ec2Instance.Prefix = *tag.Value
				}
			}
		}
	}
	return nil
}
