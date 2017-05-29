package main

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
)

//WritebackTag writes a tag describing when goat touched the disk, helpful when re-assembling the mdadm array
func WritebackTag(ebsVols []EbsVol, ec2Instance *EC2Instance, dryRun bool) error {
	volIds := []*string{}
	for _, vol := range ebsVols {
		volIds = append(volIds, &vol.EbsVolID)
	}

	touchedKey := PREFIX + "-OUT:Touched"
	touchedVal := strconv.Itoa(int(time.Now().Unix()))

	params := &ec2.CreateTagsInput{
		DryRun:    &dryRun,
		Resources: volIds,
		Tags: []*ec2.Tag{
			&ec2.Tag{
				Key:   &touchedKey,
				Value: &touchedVal,
			},
		},
	}
	_, err := ec2Instance.EC2Client.CreateTags(params)
	if err != nil {
		return err
	}
	return nil
}
