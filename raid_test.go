package main

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"strings"
	"testing"
)

func TestCreateRaidArray(t *testing.T) {
	buf := new(bytes.Buffer)
	log.SetOutput(buf)

	fakeEbsVols := []EbsVol{
		EbsVol{
			EbsVolId:     "1",
			VolumeName:   "raid_test_1",
			RaidLevel:    0,
			VolumeSize:   2,
			AttachedName: "drive_test_1",
			MountPath:    "/mnt/testy",
			FsType:       "ext999",
		},
		EbsVol{
			EbsVolId:     "2",
			VolumeName:   "raid_test_1",
			RaidLevel:    0,
			VolumeSize:   2,
			AttachedName: "drive_test_2",
			MountPath:    "/mnt/testy",
			FsType:       "ext999",
		},
	}

	if raidDrive := CreateRaidArray(fakeEbsVols, "raid_test_1", true); !strings.Contains(raidDrive, "/dev/md") {
		t.Errorf("Raid drive should be /dev/md*, actual: %s", raidDrive)
	}

	bufString := buf.String()

	if !strings.Contains(bufString, "RAID: Creating RAID drive: mdadm [--create /dev/md0 --level=0 --name='KRKN-raid_test_1' --raid-devices=2 drive_test_1 drive_test_2]") {
		t.Errorf("logged wrong thing. Actual: %s", bufString)
	}
}
