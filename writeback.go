package main

import (
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
)

//WritebackTag writes a tag describing when goat touched the disk, helpful when re-assembling the mdadm array
func WritebackTag(ebsVols []EbsVol, ec2Instance *EC2Instance, dryRun bool) error {
	for _, vol := range ebsVols {
		touchedKey := PREFIX + "-OUT:Touched"
		touchedVal := strconv.Itoa(int(time.Now().Unix()))

		params := &ec2.CreateTagsInput{
			DryRun:    &dryRun,
			Resources: []*string{&vol.EbsVolID},
			Tags: []*ec2.Tag{
				&ec2.Tag{
					Key:   &touchedKey,
					Value: &touchedVal,
				},
			},
		}
		result, err := ec2Instance.EC2Client.CreateTags(params)
		log.WithFields(log.Fields{"vol": vol}).Infof("Result when writing back tags to vol: %s", result)
		if err != nil {
			return err
		}
	}
	return nil
}
