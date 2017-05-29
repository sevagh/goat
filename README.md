[![Build Status](https://api.travis-ci.org/sevagh/ebs_raid_wizard.svg?branch=master)](https://travis-ci.org/sevagh/ebs_raid_wizard) [![ReportCard](http://goreportcard.com/badge/sevagh/ebs_raid_wizard)](http://goreportcard.com/report/sevagh/ebs_raid_wizard) [![GitHub tag](https://img.shields.io/github/tag/sevagh/ebs_raid_wizard.svg)](https://github.com/sevagh/ebs_raid_wizard/releases) [![GitPitch](https://gitpitch.com/assets/badge.svg)](https://gitpitch.com/sevagh/ebs_raid_wizard/gitpitch?grs=github&t=white)

*VERY EARLY ALPHA - USE AT YOUR OWN RISK*

# ebs_raid_wizard
### Attach & mount EBS volumes and RAID arrays from inside EC2 instances

`ebs_raid_wizard` is a Go application which runs from inside the EC2 instance (it's necessary for the instance to have an IAM Role with full EC2 access).

By setting your tags correctly, the wizard can discover, attach, RAID (with mdadm), mkfs, and mount EBS volumes to the EC2 instance where it's running.

### Behavior

The wizard should behave correctly with no parameters. It is configured entirely with tags (explained [below](#tags)). It logs to `stderr` by default.

It takes some options:

* `--dry` - dry run, don't execute any commands
* `--log-level=<level>` - logrus log levels (i.e. debug, info, warn, error, fatal, panic)

The workflow is roughly the following:

* Get EC2 metadata on the running instance
* Use metadata to establish an EC2 client and scan EBS volumes
* Attach the volumes it needs based on their tags
* Use typical Linux disk tools to mount the drives correctly:
    * `mdadm` for RAID volumes (if needed)
    * `blkid` to check for an existing filesystem
    * `mkfs` to make the filesystem
    * `/etc/fstab` entries to preserve behavior on reboot

The filesystem and fstab entries are created with the label `EWIZ-{volumeName}` for convenience. Running the wizard multiple times will result in it detecting the existing label it intended to create and not proceeding.

### Run phase

In production I run `ebs_raid_wizard` at the EC2 user-data phase (executed from a bash script). Further exploration is needed to perhaps embed it properly into `systemd` or `cloud-init`.

### Tags

These are the tags you need:

| Tag Name             | Description             | EC2     | EBS    | Tag Value (examples)                                             |
| -------------------- | ----------------------- | ------- | -----  | ---------------------------------------------------------------- |
| EWIZ-IN:Prefix       | Logical stack name      | *Yes*   | *Yes*  | `my_app_v1.3.4`                                                  |
| EWIZ-IN:NodeId       | EC2 id within stack     | *Yes*   | *Yes*  | `0`, `1`, `2` for 3-node kafka                                   |
| EWIZ-IN:VolumeName   | Distinct volume name    |         | *Yes*  | `data`, `log` - RAID disks must use the same VolumeName          |
| EWIZ-IN:VolumeSize   | # of disks in vol group |         | *Yes*  | 2 for 2-disk RAID, 1 for single disk/no RAID                     |
| EWIZ-IN:RaidLevel    | level of RAID (0 or 1)  |         | *Yes*  | 0 or 1 for RAID, ignored if VolumeSize == 1                      |
| EWIZ-IN:MountPath    | Linux path to mount vol |         | *Yes*  | `/var/kafka_data`                                                |
| EWIZ-IN:FsType       | Linux filesystem type   |         | *Yes*  | `ext4`, `vfat`                                                   |

### Examples

[Link to the example Terraform HCL scripts](https://github.com/sevagh/ebs_raid_wizard/tree/example). Also, the [Gitpitch presentation](https://gitpitch.com/sevagh/ebs_raid_wizard/gitpitch#) has a partial demonstration.

Here's an example to clarify better.

2-instance Vertica cluster:

* EC2 instance:
    * EWIZ-IN:Prefix: adgear_vertica_v0.5 
    * EWIZ-IN:NodeId: 0

* EC2 instance:
    * EWIZ-IN:Prefix: adgear_vertica_v0.5 
    * EWIZ-IN:NodeId: 1

We want 2 RAID0 volumes (2 disks each), so 4 total disks (2 per instance):

* 2x EBS volume for Node 0:
    * EWIZ-IN:Prefix: adgear_vertica_v0.5
    * EWIZ-IN:NodeId: 0
    * EWIZ-IN:VolumeName: vdata
    * EWIZ-IN:VolumeSize: 2
    * EWIZ-IN:RaidLevel: 0
    * EWIZ-IN:MountPath: /vertica/data
    * EWIZ-IN:FsType: ext4

* 2x EBS volume for Node 1:
    * EWIZ-IN:Prefix: adgear_vertica_v0.5
    * EWIZ-IN:NodeId: 1
    * EWIZ-IN:VolumeName: vdata
    * EWIZ-IN:VolumeSize: 2
    * EWIZ-IN:RaidLevel: 0
    * EWIZ-IN:MountPath: /vertica/data
    * EWIZ-IN:FsType: ext4

Additionally, we want 1 extra disk (single disks, no RAID) per node, for logs:

* 1x EBS volume for Node 0:
    * EWIZ-IN:Prefix: adgear_vertica_v0.5
    * EWIZ-IN:NodeId: 0
    * EWIZ-IN:VolumeName: vlog
    * EWIZ-IN:VolumeSize: 1
    * EWIZ-IN:RaidLevel: 0
    * EWIZ-IN:MountPath: /vertica/log
    * EWIZ-IN:FsType: ext4

* 1x EBS volume for Node 1:
    * EWIZ-IN:Prefix: adgear_vertica_v0.5
    * EWIZ-IN:NodeId: 1
    * EWIZ-IN:VolumeName: vlog
    * EWIZ-IN:VolumeSize: 1
    * EWIZ-IN:RaidLevel: 0
    * EWIZ-IN:MountPath: /vertica/log
    * EWIZ-IN:FsType: ext4

Run `ebs_raid_wizard` from the EC2 instance (ideally at the user-data phase) to automatically mount the associated EBS volumes with the above properties:

```
[dbadmin@ip-172-31-46-84 ~]$ ls /dev/disk/by-label/
EWIZ-vdata  EWIZ-vlog
[dbadmin@ip-172-31-46-84 ~]$
[dbadmin@ip-172-31-46-84 ~]$
[dbadmin@ip-172-31-46-84 ~]$ tail -n1 /etc/fstab
LABEL=EWIZ-log /vertica/log ext4 defaults 0 1
LABEL=EWIZ-vdata /vertica/data ext4 defaults 0 1
[dbadmin@ip-172-31-46-84 ~]$ sudo mdadm /dev/md0
...
           Name : "EWIZ-vdata"
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

With this wizard we can avoid needing to pass parameters or hardcoding values. All the required information comes from the EC2 instance and EBS volume tags.
