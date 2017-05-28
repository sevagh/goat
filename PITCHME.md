---
# kraken

https://github.com/sevagh/kraken

Automatically attach and mount EBS volumes in EC2 instances based on tags

+++

# About me

Sevag, http://sevag.xyz

SRE (devops?) at AdGear (now part of Samsung). I do some AWS things (automation, etc.)

+++

# Motivations

+++

Tired of re-writing same code in user-data:

```
mkdir -p /my/mount/point
mkfs.ext4 /dev/xvdb
mount /dev/xvdb /my/mount/point
```

+++

Using block device mapping or Terraform's aws_volume_attachment too inflexible

Story: 3 Kafka volumes, 3 Kafka instances. AWS marks one of them for shutdown. Using Terraform `aws_volume_attachment`, no easy solution. Using Kraken: shut down Instance X, re-create it, it picks up where the old one left off, using the same disk.

+++
Create a volume-driven cluster

Instances are less important than the disk for Kafka, Vertica, any data store. Therefore, "we need 3 Kafka brokers" -> "we need 3 Kafka volumes", the broker/EC2 instance is disposable.

---

# Example

https://github.com/sevagh/kraken/tree/master/kraken-example

A Terraform recipe for a fully working Kraken example.

```
variable "servers" {
  default = "3"
}

variable "prefix" {
  default = "kraken-ex-v0.1"
}

```

+++

Terraform apply on the example repo:

```
$ git clone https://github.com/sevagh/kraken
$ cd kraken/kraken-example
$ terraform apply
var.aws_access_key
  Enter a value: ***

var.aws_secret_key
  Enter a value: ***

var.keypair_name
  Enter a value: ***
```
+++

Kraken EC2 tags

```
resource "aws_instance" "instance" {
  count = "${var.count}"

  tags {
    "KRKN-IN:Prefix" = "${var.prefix}"
    "KRKN-IN:NodeId" = "${count.index}"
  }
}

```

+++

EC2 user-data

```
yum install -y wget mdadm
wget https://github.com/sevagh/kraken/releases/download/0.1.0/kraken
chmod +x kraken
./kraken >/var/log/kraken.log 2>&1
```

+++

EBS volumes - RAID case

```
resource "aws_ebs_volume" "data_disk" {
  count = "${var.servers * 3}"

  tags {
    Name = "${var.prefix}-data-disk-${count.index}"
    "KRKN-IN:Prefix" = "${var.prefix}"
    "KRKN-IN:NodeId" = "${count.index / 3}"
    "KRKN-IN:VolumeName" = "data"
    "KRKN-IN:VolumeSize" = "3"
    "KRKN-IN:MountPath" = "/krkn_data"
    "KRKN-IN:RaidLevel" = "0"
    "KRKN-IN:FsType" = "ext4"
  }
}
```
+++

EBS volumes - non-RAID/single disk case

```
resource "aws_ebs_volume" "log_disk" {
  count = "${var.servers}"

  tags {
    Name = "${var.prefix}-log-disk-${count.index}"
    "KRKN-IN:Prefix" = "${var.prefix}"
    "KRKN-IN:NodeId" = "${count.index}"
    "KRKN-IN:VolumeName" = "log"
    "KRKN-IN:VolumeSize" = "1"
    "KRKN-IN:MountPath" = "/krkn_log"
    "KRKN-IN:RaidLevel" = "0" #ignored since volumesize == 1
    "KRKN-IN:FsType" = "ext4"
  }
}
```

---

# Terraform outputs

EC2

```
    map[
	Name:kraken-ex-v0.1-0
	KRKN-IN:Prefix:kraken-ex-v0.1
	KRKN-IN:NodeId:0
    ] (x3 for NodeId 1 and 2)
```

+++

EBS RAID

```
    [
        map[
	    Name:kraken-ex-v0.1-data-disk-0
	    KRKN-IN:Prefix:kraken-ex-v0.1
	    KRKN-IN:NodeId:0
	    KRKN-IN:VolumeName:data
	    KRKN-IN:VolumeSize:3
	    KRKN-IN:RaidLevel:0
	    KRKN-IN:FsType:ext4
	    KRKN-IN:MountPath:/krkn_data
	] (x3 for 3 disks)
    ] (x3 for NodeId 1 and 2)
```

+++

EBS single

```
    map[
	Name:kraken-ex-v0.1-log-disk-0
	KRKN-IN:Prefix:kraken-ex-v0.1
	KRKN-IN:NodeId:0
	KRKN-IN:VolumeName:log
	KRKN-IN:VolumeSize:1
	KRKN-IN:RaidLevel:0
	KRKN-IN:FsType:ext4
	KRKN-IN:MountPath:/krkn_log
    ] (x3 for NodeId 1 and 2)
```

---

# Kraken output


```
[centos@ip-172-31-18-112 ~]$ cat /var/log/kraken.log
time="2017-05-26T17:46:37Z" level=info msg="
#####################
# WELCOME TO KRAKEN #
#####################
"
time="2017-05-26T17:46:37Z" level=info msg="
```

+++

Phase 1: EC2 metadata

```
##########################
# 1: COLLECTING EC2 INFO #
##########################
"
msg="Establishing metadata client"
msg="Retrieved metadata successfully" instance_id=i-0fe6ad25448fc8b95
msg="Using metadata to initialize EC2 SDK client" az=us-east-1a instance_id=i-0fe6ad25448fc8b95 region=us-east-1
```

+++

Phase 2: EBS scanning with tags

```
##########################
# 2: COLLECTING EBS INFO #
##########################
"
msg="Searching for EBS volumes with previously established EC2 client"
msg="Classifying EBS volumes based on tags"
```

+++

Phase 3: EBS attaching for log disk

```
#########################
# 3: ATTACHING EBS VOLS #
#########################
"
msg="Volume is unattached, picking drive name" vol_id=vol-01d2bf59335d5da42 vol_name=log
msg="Executing AWS SDK attach command" vol_id=vol-01d2bf59335d5da42 vol_name=log
msg="{
  AttachTime: 2017-05-26 17:46:38.683 +0000 UTC,
  Device: "/dev/xvdb",
  InstanceId: "i-0fe6ad25448fc8b95",
  State: "attaching",
  VolumeId: "vol-01d2bf59335d5da42"
}" vol_id=vol-01d2bf59335d5da42 vol_name=log
```

+++

Phase 3: EBS attaching for data disks

```
msg="Volume is unattached, picking drive name" vol_id=vol-06276f78a4769bb0f vol_name=data
msg="Executing AWS SDK attach command" vol_id=vol-06276f78a4769bb0f vol_name=data
msg="{
  AttachTime: 2017-05-26 17:46:40.995 +0000 UTC,
  Device: "/dev/xvdc",
  InstanceId: "i-0fe6ad25448fc8b95",
  State: "attaching",
  VolumeId: "vol-06276f78a4769bb0f"
}" vol_id=vol-06276f78a4769bb0f vol_name=data
```
(x3 for the other 2 disks, `/dev/xvdd` and `/dev/xvde`)

+++

Phase 4: Mounting single log disk

```
#############################
# 4: MOUNTING ATTACHED VOLS #
#############################
"
msg="Single drive, no RAID" vol_name=log vols=[
	{vol-01d2bf59335d5da42 log 0 1 /dev/xvdb /krkn_log ext4}
    ]
msg="Checking for existing filesystem"
msg="Checking if something already mounted at /krkn_log"
msg="Appending fstab entry"
msg="Now mounting"
```

+++

Phase 4: Mounting RAID data disks

```
msg="Creating RAID array" vol_name=data vols=[
	{vol-06276f78a4769bb0f data 0 3 /dev/xvdc /krkn_data ext4}
	{vol-0ae6037094a5b8e91 data 0 3 /dev/xvdd /krkn_data ext4}
	{vol-0383c46db8b2f23df data 0 3 /dev/xvde /krkn_data ext4}
    ]
msg="Creating RAID drive: mdadm [--create /dev/md0 --level=0 --name='KRKN-data' --raid-devices=3 /dev/xvdc /dev/xvdd /dev/xvde]"
```
(same as before - check fs, mkfs, mount)

---

# Kraken results

disk by-label:

```
[centos@ip-172-31-18-112 ~]$ ls /dev/disk/by-label/
KRKN-data  KRKN-log
```

+++

# Fstab

```
[centos@ip-172-31-18-112 ~]$ tail -n2 /etc/fstab
LABEL=KRKN-log /krkn_log ext4 defaults 0 1
LABEL=KRKN-data /krkn_data ext4 defaults 0 1
```

+++

# Mount

```
[centos@ip-172-31-18-112 ~]$ mount | tail -n2
/dev/xvdb on /krkn_log type ext4 (rw,relatime,seclabel,data=ordered)
/dev/md0 on /krkn_data type ext4 (rw,relatime,seclabel,stripe=384,data=ordered)
```

+++

# mdadm.conf

```
[centos@ip-172-31-18-112 ~]$ sudo cat /etc/mdadm.conf
ARRAY /dev/md0 level=raid0 num-devices=3 metadata=1.2 name="ip-172-31-18-112:'KRKN-data'" UUID=4a739fd5:5f392dd6:1b198f88:5cd3f50b
   devices=/dev/xvdc,/dev/xvdd,/dev/xvde
```

---

# Re-run behavior

Kraken behaves well when run again, based on labels:

```
Searching for EBS volumes with previously established EC2 client
Classifying EBS volumes based on tags
Label already exists in /dev/disk/by-label    vol_name=log
Label already exists in /dev/disk/by-label    vol_name=data
```

+++

# Dry-run behavior

No real attaching/mount/mkfs/mdadm, but it does real scanning

+++

# Fail behavior

Fail fast at every step. Don't proceed if anything is already mounted, if disk is busy, if wrong filesystem detected. Won't run `mkfs` on a filesystem that exists i.e. no chance of data loss.

+++

# Reboot behavior

Because of `mdadm.conf` and `fstab`, settings persist across reboots cleanly.

---

# Re-creating an instance

Now for the real test... destroy NodeId 1 and see what happens when we re-create it.

+++

Don't delete EBS on termination:

```
Are you sure you want to terminate these instances?
i-00151abc2eec8b58a (kraken-ex-v0.1-1, ec2-52-90-250-217.compute-1.amazonaws.com)
Clean up associated resources
 The following volumes are not set to delete on termination: [
    vol-0b44998438e1b16ba,
    vol-04a31ee820273a7aa,
    vol-090ca93e7f1ad430a,
    vol-0a63c461848cf8179,
    vol-06e1e969beb67742d
] 
```

+++

Terraform discovering that it needs to re-create the instance:

```
$ terraform plan
+ aws_instance.instance.1
    tags.%:                       "3"
    tags.KRKN-IN:NodeId:          "1"
    tags.KRKN-IN:Prefix:          "kraken-ex-v0.1"
    tags.Name:                    "kraken-ex-v0.1-1"
```

+++

After creation, logging on and checking:

```
level=fatal msg="Error when executing mdadm command: Exit Status: 2" drives=[{vol-04a31ee820273a7aa data 0 3 /dev/xvdc /krkn_data ext4} {vol-090ca93e7f1ad430a data 0 3 /dev/xvdd /krkn_data ext4} {vol-0a63c461848cf8179 data 0 3 /dev/xvde /krkn_data ext4}] vol_name=data
```

So close...

+++

Needs some work!

---

# Final thoughts & conclusion

* Why Go? Good AWS SDK, statically compiled binary

Thanks!
