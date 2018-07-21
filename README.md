[![ReportCard](http://goreportcard.com/badge/sevagh/goat)](http://goreportcard.com/report/sevagh/goat) [![GitHub tag](https://img.shields.io/github/tag/sevagh/goat.svg)](https://github.com/sevagh/goat/releases)

# goat :goat:

### Attach EBS volumes and ENIs to running EC2 instances

`goat` is a Go application which runs from inside the EC2 instance.

By setting your tags correctly, `goat` can discover and attach EBS volumes and ENIs. Furthermore, for EBS volumes, it can perform additional actions such as RAID (with mdadm), mkfs, and mount EBS volumes to the EC2 instance where it's running.

The `goat` package consists of the subcommands [goat ebs](./docs/EBS.md) and [goat eni](./docs/ENI.md).

### Permission model

It's necessary for the instance to have an IAM Role with _at least_ access to the EBS and ENI resources that it will be attaching - see [here](./docs/hcl-example/iam_role.tf). Your roles can be even more permissive (i.e. full EC2 access) but that comes with its own risks.

Unfortunately, resource-level permissions are [currently not supported](https://docs.aws.amazon.com/AWSEC2/latest/APIReference/ec2-api-permissions.html#ec2-api-unsupported-resource-permissions) for attaching network interfaces. This means that to use `goat@eni`, your instances must have full permissions for __all__ ENIs.

### RPM-based install

Goat is systemd-based and has been developed for CentOS. Install the rpm from the releases page:

```
$ sudo yum install -y https://github.com/sevagh/goat/releases/download/0.6.0/goat-0.6.0-1.fc27.x86_64.rpm
$ sudo systemctl enable goat@ebs
$ sudo systemctl start goat@ebs
```

### Additional dependencies for ENI

Goat by itself is sufficient for the EBS feature, but needs help for setting up an ENI. Refer to [this](./docs/ENI.md#setting-up-the-eni---ec2-net-utils) document.

### Hack

If you have docker locally, you can use the following command in order to quickly get a development env ready: `make dev-env`.

### Examples

[Link to the example Terraform HCL scripts](./docs/hcl-example).
