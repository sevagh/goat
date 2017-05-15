# kraken
### Attach & mount EBS volumes from inside EC2 instances

<img src="./.github/kraken-logo.png" width="300">

### What is it

Kraken is a Go application which runs from inside the EC2 instance (it's necessary for the instance to have an IAM Role with full EC2 access).

By setting your tags correctly, Kraken can discover, attach, RAID, mkfs, and mount EBS volumes to the EC2 instance where it's running.

### Motivation

The Terraform resource `aws_volume_attachment` isn't handled well when destroying a stack. See [here](https://github.com/hashicorp/terraform/issues/9000) for some discussion on the matter.

We initially wrote instance-specific user-data shell scripts with hardcoded values (e.g. `mkfs.ext4 /dev/xvdb`, `mount /dev/xvdb /var/kafka_data`).

With Kraken we can avoid needing to pass parameters or hardcoding values. All the required information comes from the EC2 instance and EBS volume tags.

### How do I use it

These are the tags you need:

| Tag Name   | Description             | EC2     | EBS    | Tag Value (examples)                                             |
| ---------- | ----------------------- | ------- | -----  | ---------------------------------------------------------------- |
| Prefix     | Logical stack name      | *Yes*   | *Yes*  | `my_app_v1.3.4`                                                  |
| NodeId     | EC2 id within stack     | *Yes*   | *Yes*  | `0`, `1`, `2` for 3-node kafka                                   |
| VolumeId   | Distinct volume id      |         | *Yes*  | `0`, `0`, for 2-disk RAID, `0`, `1` for 2 separate single disks  |
| VolumeSize | # of disks in vol group |         | *Yes*  | 2 for 2-disk RAID, 1 for single disk/no RAID                     |
| RaidLevel  | level of RAID (0 or 1)  |         | *Yes*  | 0 or 1 for RAID, ignored if VolumeSize == 1                      |
| MountPath  | Linux path to mount vol |         | *Yes*  | `/var/kafka_data`                                                |
| FsType     | Linux filesystem type   |         | *Yes*  | `ext4`, `vfat`                                                   |

Here's an example to clarify better.

#### 3 EC2 instances

We want a 3-instance Vertica cluster:

* EC2 instance:
    * Prefix: adgear_vertica_v0.5 
    * NodeId: 0

* EC2 instance:
    * Prefix: adgear_vertica_v0.5 
    * NodeId: 1

* EC2 instance:
    * Prefix: adgear_vertica_v0.5 
    * NodeId: 2

#### 9 volumes

We want 3 RAID0 volumes (2 disks each)

* 2x EBS volume for Node 0:
    * Prefix: adgear_vertica_v0.5
    * NodeId: 0
    * VolumeId: 0
    * VolumeSize: 2
    * RaidLevel: 0
    * MountPath: /vertica/data
    * FsType: ext4

* 2x EBS volume for Node 1:
    * Prefix: adgear_vertica_v0.5
    * NodeId: 1
    * VolumeId: 0
    * VolumeSize: 2
    * RaidLevel: 0
    * MountPath: /vertica/data
    * FsType: ext4

* 2x EBS volume for Node 2:
    * Prefix: adgear_vertica_v0.5
    * NodeId: 2
    * VolumeId: 0
    * VolumeSize: 2
    * RaidLevel: 0
    * MountPath: /vertica/data
    * FsType: ext4

Additionally, we want 1 extra disk (single disks, no RAID) per node, for logs:

* 1x EBS volume for Node 0:
    * Prefix: adgear_vertica_v0.5
    * NodeId: 0
    * VolumeId: 1
    * VolumeSize: 1
    * RaidLevel: 0
    * MountPath: /var/log/vertica
    * FsType: ext4

* 1x EBS volume for Node 1:
    * Prefix: adgear_vertica_v0.5
    * NodeId: 1
    * VolumeId: 1
    * VolumeSize: 1
    * RaidLevel: 0
    * MountPath: /var/log/vertica
    * FsType: ext4

* 1x EBS volume for Node 2:
    * Prefix: adgear_vertica_v0.5
    * NodeId: 2
    * VolumeId: 1
    * VolumeSize: 1
    * RaidLevel: 0
    * MountPath: /var/log/vertica
    * FsType: ext4

Run `kraken` from the EC2 instance (ideally at the user-data phase) to automatically mount the associated EBS volumes with the above properties.
