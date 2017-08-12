[![Build Status](https://api.travis-ci.org/sevagh/goat.svg?branch=master)](https://travis-ci.org/sevagh/goat) [![ReportCard](http://goreportcard.com/badge/sevagh/goat)](http://goreportcard.com/report/sevagh/goat) [![GitHub tag](https://img.shields.io/github/tag/sevagh/goat.svg)](https://github.com/sevagh/goat/releases)

*VERY EARLY ALPHA - USE AT YOUR OWN RISK*

# goat :goat:

### Attach EBS volumes and ENIs to running EC2 instances

`goat` is a Go application which runs from inside the EC2 instance (it's necessary for the instance to have an IAM Role with full EC2 access).

By setting your tags correctly, `goat` can discover and attach EBS volumes and ENIs.

Furthermore, for EBS volumes, it can perform additional actions such as RAID (with mdadm), mkfs, and mount EBS volumes to the EC2 instance where it's running.

### Motivation

The Terraform resource `aws_volume_attachment` isn't handled well when destroying a stack. See [here](https://github.com/hashicorp/terraform/issues/9000) for some discussion on the matter. We initially wrote instance-specific user-data shell scripts with hardcoded values (e.g. `mkfs.ext4 /dev/xvdb`, `mount /dev/xvdb /var/kafka_data`). With `goat` we can avoid needing to pass parameters or hardcoding values. All the required information comes from the EC2 instance and EBS volume tags.

### Subcommands

`goat` for now supports the subcommands `goat ebs` for EBS volumes and `goat eni` for ENIs.

Docs:

* [ebs](./commands/ebs/README.md)
* [eni](./commands/eni/README.md)

### RPM-based install

Goat is systemd-based and has been developed for CentOS. Install the rpm from the releases page:

```
$ sudo yum install -y https://github.com/sevagh/goat/releases/download/0.4.0/goat-0.4.0-1.fc25.x86_64.rpm
$ sudo systemctl enable goat@ebs
$ sudo systemctl start goat@ebs
$ ...
$ journalctl -u goat@ebs
```

#### Examples

[Link to the example Terraform HCL scripts](https://github.com/sevagh/goat-example).
