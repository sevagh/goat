[![ReportCard](http://goreportcard.com/badge/sevagh/goat)](http://goreportcard.com/report/sevagh/goat) [![GitHub tag](https://img.shields.io/github/tag/sevagh/goat.svg)](https://github.com/sevagh/goat/releases)

*VERY EARLY ALPHA - USE AT YOUR OWN RISK*

# goat :goat:

### Attach EBS volumes and ENIs to running EC2 instances

`goat` is a Go application which runs from inside the EC2 instance.

By setting your tags correctly, `goat` can discover and attach EBS volumes and ENIs.

Furthermore, for EBS volumes, it can perform additional actions such as RAID (with mdadm), mkfs, and mount EBS volumes to the EC2 instance where it's running.

### Permission model

It's necessary for the instance to have an IAM Role with _at least_ access to the EBS and ENI resources that it will be attaching - see [here](https://github.com/sevagh/goat-example/blob/master/iam_role.tf). Your roles can be even more permissive (i.e. full EC2 access) but that comes with its own risks.

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

### Additional dependencies for ENI

Goat by itself is sufficient for the EBS feature, but needs help for setting up an ENI.

Refer to [this](./commands/eni#setting-up-the-eni---ec2-net-utils) document. It refers to a [port of ec2-net-utils](https://github.com/sevagh/ec2-utils/releases) from the Amazon Linux AMI to CentOS/systemd.

A fully working chunk of `ec2 user-data` with `goat` looks like [this](https://github.com/sevagh/goat-example/blob/master/bootstrap.tpl#L8):

```
yum install -y wget mdadm
yum install -y https://github.com/sevagh/goat/releases/download/0.4.0/goat-0.4.0-1.fc25.x86_64.rpm
yum install -y https://github.com/sevagh/ec2-utils/releases/download/v0.5.3/ec2-net-utils-0.5-2.fc25.noarch.rpm
systemctl enable elastic-network-interfaces
systemctl start elastic-network-interfaces
systemctl enable goat@ebs
systemctl enable goat@eni
systemctl start goat@ebs
systemctl start goat@eni
```

### Examples

[Link to the example Terraform HCL scripts](https://github.com/sevagh/goat-example).
