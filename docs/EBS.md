### Goat for EBS

#### Behavior

`goat ebs` should behave correctly with no parameters. It is configured entirely with tags (explained [below](#tags)). It logs to `stderr` by default.

It takes some options:

* `-logLevel=<level>` - logrus log levels (i.e. debug, info, warn, error, fatal, panic)
* `-debug` - an interactive debug mode which prompts to continue after every phase so you can explore the state between phases

#### Tags

These are the tags you need:

| Tag Name             | Description             | EC2     | EBS    | Required | Tag Value (examples)                                             |
| -------------------- | ----------------------- | ------- | -----  | -------- | ---------------------------------------------------------------- |
| GOAT-IN:Prefix       | Logical stack name      | *Yes*   | *Yes*  | *YES*    | `my_app_v1.3.4`                                                  |
| GOAT-IN:NodeId       | EC2 id within stack     | *Yes*   | *Yes*  | *YES*    | `0`, `1`, `2` for 3-node kafka                                   |
| GOAT-IN:VolumeName   | Distinct volume name    |         | *Yes*  | no       | `data`, `log` - RAID disks must use the same VolumeName          |
| GOAT-IN:VolumeSize   | # of disks in vol group |         | *Yes*  | no       | 2 for 2-disk RAID, 1 for single disk/no RAID                     |
| GOAT-IN:RaidLevel    | level of RAID (0 or 1)  |         | *Yes*  | no       | 0 or 1 for RAID, ignored if VolumeSize == 1                      |
| GOAT-IN:MountPath    | Linux path to mount vol |         | *Yes*  | no       | `/var/kafka_data`                                                |
| GOAT-IN:FsType       | Linux filesystem type   |         | *Yes*  | no       | `ext4`, `vfat`                                                   |

#### Missing tags

If a non-required tag is missing, that step will be skipped. E.g. without RaidLevel or VolumeSize, `mdadm` won't be run. Without a filesystem, `mkfs` won't be run. Without a mount path, `mount` won't be run. Without a volume name, nothing will be run.

*The barest case will simply attach the EBS volumes and perform no further actions*

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
