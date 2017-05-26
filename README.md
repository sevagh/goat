[![Build Status](https://api.travis-ci.org/sevagh/kraken.svg?branch=master)](https://travis-ci.org/sevagh/kraken) [![ReportCard](http://goreportcard.com/badge/sevagh/kraken)](http://goreportcard.com/report/sevagh/kraken) [![GitHub tag](https://img.shields.io/github/tag/sevagh/kraken.svg)](https://github.com/sevagh/kraken/releases) [![GitPitch](https://gitpitch.com/assets/badge.svg)](https://gitpitch.com/sevagh/kraken/gitpitch?grs=github&t=white)

*VERY EARLY ALPHA - USE AT YOUR OWN RISK*

# kraken
### Attach & mount EBS volumes from inside EC2 instances

<img src="https://raw.githubusercontent.com/sevagh/kraken/gitpitch/assets/logo.png" width="300">

### Introduction

Kraken is a Go application which runs from inside the EC2 instance (it's necessary for the instance to have an IAM Role with full EC2 access).

By setting your tags correctly, Kraken can discover, attach, RAID, mkfs, and mount EBS volumes to the EC2 instance where it's running.

### Behavior

Kraken should behave correctly with no parameters. It is configured entirely with tags (explained [below](#tags)). It logs to `stderr` by default.

It takes some options:

* `--dry` - dry run, don't execute any commands
* `--log-level=<level>` - logrus log levels (i.e. debug, info, warn, error, fatal, panic)

Kraken's workflow is roughly the following:

* Get EC2 metadata on the running instance
* Use metadata to establish an EC2 client and scan EBS volumes
* Attach the volumes it needs based on their tags
* Use typical Linux disk tools to mount the drives correctly:
    * `mdadm` for RAID volumes (if needed)
    * `blkid` to check for an existing filesystem
    * `mkfs` to make the filesystem
    * `/etc/fstab` entries to preserve behavior on reboot

The filesystem and fstab entries are created with the label `KRKN-{volumeName}` for convenience. Running Kraken multiple times will result in Kraken detecting the existing label it intended to create and not proceeding.

### Run phase

In production I run `kraken` at the EC2 user-data phase (executed from a bash script). Further exploration is needed to perhaps embed kraken properly into `systemd` or `cloud-init`.

### Tags

These are the tags you need:

| Tag Name             | Description             | EC2     | EBS    | Tag Value (examples)                                             |
| -------------------- | ----------------------- | ------- | -----  | ---------------------------------------------------------------- |
| KRKN-IN:Prefix       | Logical stack name      | *Yes*   | *Yes*  | `my_app_v1.3.4`                                                  |
| KRKN-IN:NodeId       | EC2 id within stack     | *Yes*   | *Yes*  | `0`, `1`, `2` for 3-node kafka                                   |
| KRKN-IN:VolumeName   | Distinct volume name    |         | *Yes*  | `data`, `log` - RAID disks must use the same VolumeName          |
| KRKN-IN:VolumeSize   | # of disks in vol group |         | *Yes*  | 2 for 2-disk RAID, 1 for single disk/no RAID                     |
| KRKN-IN:RaidLevel    | level of RAID (0 or 1)  |         | *Yes*  | 0 or 1 for RAID, ignored if VolumeSize == 1                      |
| KRKN-IN:MountPath    | Linux path to mount vol |         | *Yes*  | `/var/kafka_data`                                                |
| KRKN-IN:FsType       | Linux filesystem type   |         | *Yes*  | `ext4`, `vfat`                                                   |

### Examples

[Link to my official example repo](https://github.com/sevagh/kraken-example). Also, the [Gitpitch presentation](https://gitpitch.com/sevagh/kraken/gitpitch#) has a partial demonstration.

Here's an example to clarify better.

2-instance Vertica cluster:

* EC2 instance:
    * KRKN-IN:Prefix: adgear_vertica_v0.5 
    * KRKN-IN:NodeId: 0

* EC2 instance:
    * KRKN-IN:Prefix: adgear_vertica_v0.5 
    * KRKN-IN:NodeId: 1

We want 2 RAID0 volumes (2 disks each), so 4 total disks (2 per instance):

* 2x EBS volume for Node 0:
    * KRKN-IN:Prefix: adgear_vertica_v0.5
    * KRKN-IN:NodeId: 0
    * KRKN-IN:VolumeName: vdata
    * KRKN-IN:VolumeSize: 2
    * KRKN-IN:RaidLevel: 0
    * KRKN-IN:MountPath: /vertica/data
    * KRKN-IN:FsType: ext4

* 2x EBS volume for Node 1:
    * KRKN-IN:Prefix: adgear_vertica_v0.5
    * KRKN-IN:NodeId: 1
    * KRKN-IN:VolumeName: vdata
    * KRKN-IN:VolumeSize: 2
    * KRKN-IN:RaidLevel: 0
    * KRKN-IN:MountPath: /vertica/data
    * KRKN-IN:FsType: ext4

Additionally, we want 1 extra disk (single disks, no RAID) per node, for logs:

* 1x EBS volume for Node 0:
    * KRKN-IN:Prefix: adgear_vertica_v0.5
    * KRKN-IN:NodeId: 0
    * KRKN-IN:VolumeName: vlog
    * KRKN-IN:VolumeSize: 1
    * KRKN-IN:RaidLevel: 0
    * KRKN-IN:MountPath: /vertica/log
    * KRKN-IN:FsType: ext4

* 1x EBS volume for Node 1:
    * KRKN-IN:Prefix: adgear_vertica_v0.5
    * KRKN-IN:NodeId: 1
    * KRKN-IN:VolumeName: vlog
    * KRKN-IN:VolumeSize: 1
    * KRKN-IN:RaidLevel: 0
    * KRKN-IN:MountPath: /vertica/log
    * KRKN-IN:FsType: ext4

Run `kraken` from the EC2 instance (ideally at the user-data phase) to automatically mount the associated EBS volumes with the above properties:

```
[dbadmin@ip-172-31-46-84 ~]$ ls /dev/disk/by-label/
KRKN-vdata  KRKN-vlog
[dbadmin@ip-172-31-46-84 ~]$
[dbadmin@ip-172-31-46-84 ~]$
[dbadmin@ip-172-31-46-84 ~]$ tail -n1 /etc/fstab
LABEL=KRKN-log /vertica/log ext4 defaults 0 1
LABEL=KRKN-vdata /vertica/data ext4 defaults 0 1
[dbadmin@ip-172-31-46-84 ~]$ sudo mdadm /dev/md0
...
           Name : "KRKN-vdata"
           UUID : ddf766d8:ab73e885:e6e35b57:3dc7eb28
         Events : 0
...
[dbadmin@ip-172-31-46-84 ~]$ df -h | tail -n2
/dev/xvdc        20G   45M   19G   1% /vertica/log
/dev/md0      222G   61M  210G   1% /vertica/data
```

### Motivation

The Terraform resource `aws_volume_attachment` isn't handled well when destroying a stack. See [here](https://github.com/hashicorp/terraform/issues/9000) for some discussion on the matter.

We initially wrote instance-specific user-data shell scripts with hardcoded values (e.g. `mkfs.ext4 /dev/xvdb`, `mount /dev/xvdb /var/kafka_data`).

With Kraken we can avoid needing to pass parameters or hardcoding values. All the required information comes from the EC2 instance and EBS volume tags.
