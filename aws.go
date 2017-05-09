package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type attachData struct {
	instanceId    string
	prefix        string
	nodeId        string
	volumes       map[int]string
	attachedNames []string
}

func AttachEbsVolumes(logger *log.Logger) ([]string, error) {
	ret := []string{}
	var ad attachData
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
		return ret, err
	}

	region, err := getRegion(svc)
	if err != nil {
		return ret, err
	}

	ad.instanceId = result
	sess.Config.Region = &region
	ec2_svc := ec2.New(sess)

	err = getTags(ec2_svc, &ad)
	if err != nil {
		return ret, err
	}

	err = findEbsVolumes(ec2_svc, &ad, logger)
	if err != nil {
		return ret, err
	}

	deviceNames, err := attachVolumes(ec2_svc, &ad, logger)
	if err != nil {
		return ret, err
	}

	return deviceNames, nil
}

func getRegion(svc *ec2metadata.EC2Metadata) (string, error) {
	id, err := svc.GetInstanceIdentityDocument()
	if err != nil {
		return "", err
	}
	return id.Region, nil
}

func getTags(svc *ec2.EC2, ad *attachData) error {
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("instance-id"),
				Values: []*string{
					aws.String(ad.instanceId),
				},
			},
		},
	}

	result, err := svc.DescribeInstances(params)
	if err != nil {
		return err
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			for _, tag := range instance.Tags {
				if *tag.Key == "NodeId" {
					ad.nodeId = *tag.Value
				} else if *tag.Key == "Prefix" {
					ad.prefix = *tag.Value
				}
			}
		}
	}
	return nil
}

func findEbsVolumes(svc *ec2.EC2, ad *attachData, logger *log.Logger) error {
	params := &ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag:Prefix"),
				Values: []*string{
					aws.String(ad.prefix),
				},
			},
			&ec2.Filter{
				Name: aws.String("tag:VolumeId"),
				Values: []*string{
					aws.String(ad.nodeId),
				},
			},
		},
	}

	ad.volumes = make(map[int]string)

	result, err := svc.DescribeVolumes(params)
	if err != nil {
		return err
	}

	ad.attachedNames = []string{}

OUTER:
	for _, volume := range result.Volumes {
		if len(volume.Attachments) > 0 {
			log.Printf("Active attachments on volume %s, investigating...", *volume.VolumeId)
			for _, attachment := range volume.Attachments {
				if *attachment.InstanceId != ad.instanceId {
					return fmt.Errorf("Volume %s attached to different instance-id: %s", *volume.VolumeId, attachment.InstanceId)
				}
				log.Printf("Active attachment is on current instance-id, continuing")
				ad.attachedNames = append(ad.attachedNames, *attachment.Device)
				continue OUTER
			}
		}
		for _, tag := range volume.Tags {
			if *tag.Key == "DiskId" {
				volKey, err := strconv.Atoi(*tag.Value)
				if err != nil {
					logger.Printf("Couldn't parse tag DiskId for vol %s as int", *volume.VolumeId)
					return err
				}
				ad.volumes[volKey] = *volume.VolumeId
				break
			}
		}
	}
	return nil
}

func attachVolumes(svc *ec2.EC2, ad *attachData, logger *log.Logger) ([]string, error) {
	deviceNames := []string{}
	if len(ad.volumes) == 0 {
		logger.Println("Nothing to attach, returning existing attached device names")
		return ad.attachedNames, nil
	}
	for diskId, volume := range ad.volumes {
		deviceName := generateDeviceName(diskId)
		attachVolIn := &ec2.AttachVolumeInput{
			Device:     &deviceName,
			InstanceId: &ad.instanceId,
			VolumeId:   &volume,
		}
		volAttachments, err := svc.AttachVolume(attachVolIn)
		if err != nil {
			return deviceNames, err
		}
		logger.Println(volAttachments)
		deviceNames = append(deviceNames, deviceName)
	}
	return deviceNames, nil
}

var diskIdLetterMappings = map[int]string{
	0: "b",
	1: "c",
	2: "d",
	3: "e",
	4: "f",
	5: "g",
	6: "h",
	7: "i",
	8: "j",
	9: "k",
}

func generateDeviceName(diskId int) string {
	return "/dev/xvd" + diskIdLetterMappings[diskId]
}
