[![Build Status](https://api.travis-ci.org/sevagh/goat.svg?branch=master)](https://travis-ci.org/sevagh/goat) [![ReportCard](http://goreportcard.com/badge/sevagh/goat)](http://goreportcard.com/report/sevagh/goat) [![GitHub tag](https://img.shields.io/github/tag/sevagh/goat.svg)](https://github.com/sevagh/goat/releases) [![GitPitch](https://gitpitch.com/assets/badge.svg)](https://gitpitch.com/sevagh/goat/gitpitch?grs=github&t=white)

*VERY EARLY ALPHA - USE AT YOUR OWN RISK*

# goat :goat:
### Attach & mount EBS volumes and RAID arrays from inside EC2 instances

`goat` is a Go application which runs from inside the EC2 instance (it's necessary for the instance to have an IAM Role with full EC2 access).

By setting your tags correctly, `goat` can discover, attach, RAID (with mdadm), mkfs, and mount EBS volumes to the EC2 instance where it's running.

### RPM-based install

Goat is systemd-based (you can download the binary and run it yourself for alternative cases) and has been developed for CentOS.

Install the rpm from the releases page, which enables it in systemd automatically:

```
$ wget https://github.com/sevagh/goat/releases/download/0.2.0/goat.rpm
$ sudo yum install ./goat.rpm
$ sudo systemctl start goat
$ journalctl -u goat
```

### Behavior

`goat` should behave correctly with no parameters. It is configured entirely with tags (explained [below](#tags)). It logs to `stderr` by default.

It takes some options:

* `--dry` - dry run, don't execute any commands
* `--log-level=<level>` - logrus log levels (i.e. debug, info, warn, error, fatal, panic)
* `--debug` - an interactive debug mode which prompts to continue after every phase so you can explore the state between phases

#### Fresh run

The event flow is roughly the following:

* Get EC2 metadata on the running instance
* Use metadata to establish an EC2 client and scan EBS volumes
* Attach the volumes it needs based on their tags
* Use typical Linux disk tools to mount the drives correctly:
    * `mdadm` for RAID volumes (if needed)
    * `blkid` to check for an existing filesystem
    * `mkfs` to make the filesystem
    * `/etc/fstab` entries to preserve behavior on reboot

The filesystem and fstab entries are created with the label `GOAT-{volumeName}` for convenience. Running `goat` multiple times will result in it detecting the existing label it intended to create and not proceeding.

#### Attaching old disks to a new instance

The specific use-case that `goat` was developed to solve is the following. Say you have 3 instances with their own respective disks, and you receive a termination notice for instance 1. I want the `goat` workflow to be:

* Terminate instance 1
* Create a new one with the same GOAT tags (to indicate that it's the same as the machine it is replacing)
* Everything works magically

After `goat` ran on the previous instance with fresh disks, the disks have the correct filesystems, labels, and in the case of RAID, `mdadm` metadata on them.

The event flow on a re-created instance (with disks that were previously attached to another instance) is:

* Get EC2 metadata
* Use metadata to establish an EC2 client and scan EBS volumes
* Attach the volumes it needs based on their tags
* Discover that `/dev/disk/by-label` already contains the correct disks
    * From `mdadm` magic, after the EBS attachment the RAID array is already detected correctly
* Proceed to perform the `fstab` and `mount` phases - skip `mdadm`, `mkfs`

**CAVEAT**: the mdadm metadata will have the hostname of the previous EC2 instance:

```
[centos@ip-172-31-29-69 ~]$ sudo mdadm --detail --scan --verbose
ARRAY /dev/md127 level=raid0 num-devices=3 metadata=1.2 name="ip-172-31-25-105:'GOAT-data'" UUID=2d08b310:fd13bd21:bc2417a4:56a1ec57
   devices=/dev/xvdb,/dev/xvdc,/dev/xvdd
[centos@ip-172-31-29-69 ~]$
```

To avoid this, define a good/persistent hostname for EC2 instance, that you will then re-apply to any instance taking over this instance's disks.

### Run phase

In production I run `goat` at the EC2 user-data phase (executed from a bash script). Further exploration is needed to perhaps embed it properly into `systemd` or `cloud-init`.

### Tags

These are the tags you need:

| Tag Name             | Description             | EC2     | EBS    | Tag Value (examples)                                             |
| -------------------- | ----------------------- | ------- | -----  | ---------------------------------------------------------------- |
| GOAT-IN:Prefix       | Logical stack name      | *Yes*   | *Yes*  | `my_app_v1.3.4`                                                  |
| GOAT-IN:NodeId       | EC2 id within stack     | *Yes*   | *Yes*  | `0`, `1`, `2` for 3-node kafka                                   |
| GOAT-IN:VolumeName   | Distinct volume name    |         | *Yes*  | `data`, `log` - RAID disks must use the same VolumeName          |
| GOAT-IN:VolumeSize   | # of disks in vol group |         | *Yes*  | 2 for 2-disk RAID, 1 for single disk/no RAID                     |
| GOAT-IN:RaidLevel    | level of RAID (0 or 1)  |         | *Yes*  | 0 or 1 for RAID, ignored if VolumeSize == 1                      |
| GOAT-IN:MountPath    | Linux path to mount vol |         | *Yes*  | `/var/kafka_data`                                                |
| GOAT-IN:FsType       | Linux filesystem type   |         | *Yes*  | `ext4`, `vfat`                                                   |

### Examples

[Link to the example Terraform HCL scripts](https://github.com/sevagh/goat/tree/example). Also, the [Gitpitch presentation](https://gitpitch.com/sevagh/goat/gitpitch#) has a partial demonstration.

Here's an example to clarify better.

2-instance Vertica cluster:

* EC2 instance:
    * GOAT-IN:Prefix: adgear_vertica_v0.5 
    * GOAT-IN:NodeId: 0

* EC2 instance:
    * GOAT-IN:Prefix: adgear_vertica_v0.5 
    * GOAT-IN:NodeId: 1

We want 2 RAID0 volumes (2 disks each), so 4 total disks (2 per instance):

* 2x EBS volume for Node 0:
    * GOAT-IN:Prefix: adgear_vertica_v0.5
    * GOAT-IN:NodeId: 0
    * GOAT-IN:VolumeName: vdata
    * GOAT-IN:VolumeSize: 2
    * GOAT-IN:RaidLevel: 0
    * GOAT-IN:MountPath: /vertica/data
    * GOAT-IN:FsType: ext4

* 2x EBS volume for Node 1:
    * GOAT-IN:Prefix: adgear_vertica_v0.5
    * GOAT-IN:NodeId: 1
    * GOAT-IN:VolumeName: vdata
    * GOAT-IN:VolumeSize: 2
    * GOAT-IN:RaidLevel: 0
    * GOAT-IN:MountPath: /vertica/data
    * GOAT-IN:FsType: ext4

Additionally, we want 1 extra disk (single disks, no RAID) per node, for logs:

* 1x EBS volume for Node 0:
    * GOAT-IN:Prefix: adgear_vertica_v0.5
    * GOAT-IN:NodeId: 0
    * GOAT-IN:VolumeName: vlog
    * GOAT-IN:VolumeSize: 1
    * GOAT-IN:RaidLevel: 0
    * GOAT-IN:MountPath: /vertica/log
    * GOAT-IN:FsType: ext4

* 1x EBS volume for Node 1:
    * GOAT-IN:Prefix: adgear_vertica_v0.5
    * GOAT-IN:NodeId: 1
    * GOAT-IN:VolumeName: vlog
    * GOAT-IN:VolumeSize: 1
    * GOAT-IN:RaidLevel: 0
    * GOAT-IN:MountPath: /vertica/log
    * GOAT-IN:FsType: ext4

Run `goat` from the EC2 instance (ideally at the user-data phase) to automatically mount the associated EBS volumes with the above properties:

```
[dbadmin@ip-172-31-46-84 ~]$ ls /dev/disk/by-label/
GOAT-vdata  GOAT-vlog
[dbadmin@ip-172-31-46-84 ~]$
[dbadmin@ip-172-31-46-84 ~]$
[dbadmin@ip-172-31-46-84 ~]$ tail -n1 /etc/fstab
LABEL=GOAT-log /vertica/log ext4 defaults 0 1
LABEL=GOAT-vdata /vertica/data ext4 defaults 0 1
[dbadmin@ip-172-31-46-84 ~]$ sudo mdadm /dev/md0
...
           Name : "GOAT-vdata"
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

With `goat` we can avoid needing to pass parameters or hardcoding values. All the required information comes from the EC2 instance and EBS volume tags.
