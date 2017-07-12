package main

import (
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

//EC2Instance is a struct containing current instance info + a connected EC2 client
type EC2Instance struct {
	EC2Client  *ec2.EC2
	InstanceID string
	Prefix     string
	NodeID     string
	Az         string
}

type ec2Metadata struct {
	region string
	az     string
}

//GetEC2InstanceData returns a populated EC2Instance struct with the current EC2 instances' metadata
func GetEC2InstanceData() EC2Instance {
	var ec2Instance EC2Instance
	sess := session.New()

	creds := credentials.NewCredentials(
		&ec2rolecreds.EC2RoleProvider{
			Client:       ec2metadata.New(sess),
			ExpiryWindow: 5 * time.Minute,
		},
	)

	log.Info("Establishing metadata client")

	sess.Config.Credentials = creds
	svc := ec2metadata.New(sess)

	result, err := svc.GetMetadata("instance-id")
	if err != nil {
		log.Fatalf("Couldn't get self instance-id from metadata: %v", err)
	}

	ec2Instance.InstanceID = result

	meta, err := populateRegionInfo(svc)
	if err != nil {
		log.Fatalf("Couldn't access InstanceIdentityDocument: %v", err)
	}

	ec2Instance.Az = meta.az
	sess.Config.Region = &meta.region

	log.WithFields(log.Fields{"instance_id": result}).Info("Retrieved metadata successfully")

	ec2Logger := log.WithFields(log.Fields{"instance_id": result, "region": meta.region, "az": meta.az})
	ec2Logger.Info("Using metadata to initialize EC2 SDK client")

	ec2Svc := ec2.New(sess)
	ec2Instance.EC2Client = ec2Svc

	err = getInstanceTags(&ec2Instance)
	if err != nil {
		ec2Logger.Fatalf("Couldn't get tags: %v", err)
	}

	if ec2Instance.NodeID == "" || ec2Instance.Prefix == "" {
		ec2Logger.Fatal("This instance is missing required GOAT-IN tags NodeId, Prefix")
	}

	return ec2Instance
}

func populateRegionInfo(svc *ec2metadata.EC2Metadata) (ec2Metadata, error) {
	ret := ec2Metadata{}
	id, err := svc.GetInstanceIdentityDocument()
	if err != nil {
		return ret, err
	}
	ret.az = id.AvailabilityZone
	ret.region = id.Region
	return ret, nil
}

func getInstanceTags(ec2Instance *EC2Instance) error {
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("instance-id"),
				Values: []*string{
					aws.String(ec2Instance.InstanceID),
				},
			},
		},
	}

	result, err := ec2Instance.EC2Client.DescribeInstances(params)
	if err != nil {
		return err
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			for _, tag := range instance.Tags {
				if *tag.Key == PREFIX+"-IN:NodeId" {
					ec2Instance.NodeID = *tag.Value
				} else if *tag.Key == PREFIX+"-IN:Prefix" {
					ec2Instance.Prefix = *tag.Value
				}
			}
		}
	}
	return nil
}
