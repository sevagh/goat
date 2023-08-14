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
	Region     string
	Vols       map[string][]EbsVol
	Enis       []string
}

//GetEC2InstanceData returns a populated EC2Instance struct with the current EC2 instances' metadata
func GetEC2InstanceData(tagPrefix string) EC2Instance {
	var ec2Instance EC2Instance
	sess := session.New(&aws.Config{
    MaxRetries: aws.Int(5),
})

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

	if err := ec2Instance.populateRegionInfo(svc); err != nil {
		log.Fatalf("Couldn't access InstanceIdentityDocument: %v", err)
	}

	sess.Config.Region = &ec2Instance.Region

	log.WithFields(log.Fields{"instance_id": result}).Info("Retrieved metadata successfully")

	ec2Logger := log.WithFields(log.Fields{"instance_id": result, "region": ec2Instance.Region, "az": ec2Instance.Az})
	ec2Logger.Info("Using metadata to initialize EC2 SDK client")

	ec2Svc := ec2.New(sess)
	ec2Instance.EC2Client = ec2Svc

	err = ec2Instance.getInstanceTags(tagPrefix)
	if err != nil {
		ec2Logger.Fatalf("Couldn't get tags: %v", err)
	}

	if ec2Instance.NodeID == "" || ec2Instance.Prefix == "" {
		ec2Logger.Fatalf("This instance is missing required '%s' tags NodeId, Prefix", tagPrefix)
	}

	return ec2Instance
}

func (e *EC2Instance) populateRegionInfo(svc *ec2metadata.EC2Metadata) error {
	id, err := svc.GetInstanceIdentityDocument()
	if err != nil {
		return err
	}
	e.Az = id.AvailabilityZone
	e.Region = id.Region
	return nil
}

func (e *EC2Instance) getInstanceTags(tagPrefix string) error {
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("instance-id"),
				Values: []*string{
					aws.String(e.InstanceID),
				},
			},
		},
	}

	result, err := e.EC2Client.DescribeInstances(params)
	if err != nil {
		return err
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			for _, tag := range instance.Tags {
				if *tag.Key == tagPrefix+":NodeId" {
					e.NodeID = *tag.Value
				} else if *tag.Key == tagPrefix+":Prefix" {
					e.Prefix = *tag.Value
				}
			}
		}
	}
	return nil
}
