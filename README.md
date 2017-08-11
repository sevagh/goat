# example

Terraform recipe for a minimally viable example of the [goat](https://github.com/sevagh/goat) EBS/EC2 tag-based mounting system.

### Iterating

To iterate with this Terraform recipe, it's helpful to export the 3 required variables:

```
$ export TF_VAR_aws_access_key=xxxx
$ export TF_VAR_aws_secret_key=xxxx
$ export TF_VAR_keypair_name=mykeypair
```

### Goat for EBS

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

Result:

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
