package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type EbsVol struct {
	EbsVolId     string
	VolumeId     int
	RaidLevel    int
	VolumeSize   int
	AttachedName string
	MountPath    string
	FsType       string
}

func FindEbsVolumes(ec2Instance *Ec2Instance) ([]EbsVol, error) {
	params := &ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag:KRAKEN-IN:Prefix"),
				Values: []*string{
					aws.String(ec2Instance.Prefix),
				},
			},
			&ec2.Filter{
				Name: aws.String("tag:KRAKEN-IN:NodeId"),
				Values: []*string{
					aws.String(ec2Instance.NodeId),
				},
			},
			&ec2.Filter{
				Name: aws.String("availability-zone"),
				Values: []*string{
					aws.String(ec2Instance.Az),
				},
			},
		},
	}

	volumes := []EbsVol{}

	result, err := ec2Instance.Ec2Client.DescribeVolumes(params)
	if err != nil {
		return volumes, err
	}

	for _, volume := range result.Volumes {
		ebsVolume := EbsVol{
			EbsVolId: *volume.VolumeId,
		}
		if len(volume.Attachments) > 0 {
			log.Printf("Active attachments on volume %s, investigating...", *volume.VolumeId)
			for _, attachment := range volume.Attachments {
				if *attachment.InstanceId != ec2Instance.InstanceId {
					return volumes, fmt.Errorf("Volume %s attached to different instance-id: %s", *volume.VolumeId, attachment.InstanceId)
				}
				log.Printf("Active attachment is on current instance-id, continuing")
				ebsVolume.AttachedName = *attachment.Device
			}
		} else {
			ebsVolume.AttachedName = ""
		}
		for _, tag := range volume.Tags {
			switch *tag.Key {
			case "KRAKEN-IN:VolumeId":
				if ebsVolume.VolumeId, err = strconv.Atoi(*tag.Value); err != nil {
					log.Printf("Couldn't parse tag VolumeId for vol %s as int", *volume.VolumeId)
					return volumes, err
				}
			case "KRAKEN-IN:RaidLevel":
				if ebsVolume.RaidLevel, err = strconv.Atoi(*tag.Value); err != nil {
					log.Printf("Couldn't parse tag RaidLevel for vol %s as int", *volume.VolumeId)
					return volumes, err
				}
			case "KRAKEN-IN:VolumeSize":
				if ebsVolume.VolumeSize, err = strconv.Atoi(*tag.Value); err != nil {
					log.Printf("Couldn't parse tag VolumeSize for vol %s as int", *volume.VolumeId)
					return volumes, err
				}
			case "KRAKEN-IN:MountPath":
				ebsVolume.MountPath = *tag.Value
			case "KRAKEN-IN:FsType":
				ebsVolume.FsType = *tag.Value
			case "KRAKEN-IN:NodeId": //do nothing
			case "KRAKEN-IN:Prefix": //do nothing
			default:
				log.Printf("Unrecognized tag %s for vol %s, ignoring...", *tag.Key, *volume.VolumeId)
			}
		}
		volumes = append(volumes, ebsVolume)
	}
	return volumes, nil
}
