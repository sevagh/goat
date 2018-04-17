[![ReportCard](http://goreportcard.com/badge/sevagh/goat)](http://goreportcard.com/report/sevagh/goat) [![GitHub tag](https://img.shields.io/github/tag/sevagh/goat.svg)](https://github.com/sevagh/goat/releases)

# goat :goat:

### Attach EBS volumes and ENIs to running EC2 instances

`goat` is a Go application which runs from inside the EC2 instance.

By setting your tags correctly, `goat` can discover and attach EBS volumes and ENIs.

Furthermore, for EBS volumes, it can perform additional actions such as RAID (with mdadm), mkfs, and mount EBS volumes to the EC2 instance where it's running.

The `goat` package consists of the subcommands [goat ebs](./docs/EBS.md) and [goat eni](./docs/ENI.md).

### Permission model

It's necessary for the instance to have an IAM Role with _at least_ access to the EBS and ENI resources that it will be attaching - see [here](./docs/hcl-example/iam_role.tf). Your roles can be even more permissive (i.e. full EC2 access) but that comes with its own risks.

Unfortunately, resource-level permissions are [currently not supported](https://docs.aws.amazon.com/AWSEC2/latest/APIReference/ec2-api-permissions.html#ec2-api-unsupported-resource-permissions) for attaching network interfaces. This means that to use `goat@eni`, your instances must have full permissions for __all__ ENIs.

### RPM-based install

Goat is systemd-based and has been developed for CentOS. Install the rpm from the releases page:

```
$ sudo yum install -y https://github.com/sevagh/goat/releases/download/0.4.1/goat-ebs-0.4.2-1.fc25.x86_64.rpm
$ sudo systemctl enable goat-ebs
$ sudo systemctl start goat-ebs
```

### Additional dependencies for ENI

Goat by itself is sufficient for the EBS feature, but needs help for setting up an ENI.

Refer to [this](./commands/eni#setting-up-the-eni---ec2-net-utils) document. It refers to a [port of ec2-net-utils](https://github.com/sevagh/ec2-utils/releases) from the Amazon Linux AMI to CentOS/systemd.

A fully working chunk of `ec2 user-data` with `goat` looks like [this](./hcl-example/blob/master/bootstrap.tpl#L8):

```
$ sudo yum install -y https://github.com/sevagh/goat/releases/download/0.4.0/goat-eni-0.4.2-1.fc25.x86_64.rpm
$ sudo yum install -y https://github.com/sevagh/ec2-utils/releases/download/v0.5.3/ec2-net-utils-0.5-2.fc25.noarch.rpm
$ sudo systemctl enable elastic-network-interfaces
$ sudo systemctl start elastic-network-interfaces
$ sudo systemctl enable goat-eni
$ sudo systemctl start goat-eni
```

### Examples

[Link to the example Terraform HCL scripts](./docs/hcl-example).
