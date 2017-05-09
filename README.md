# kraken
### Attach & mount EBS volumes from inside EC2 instances

## What is it

Kraken runs from inside the EC2 instance (it's necessary for the instance to have an IAM Role with full EC2 access). 

By associating `instance tag.Prefix <-> volume tag.Prefix` and `instance tag.NodeId <-> volume tag.VolumeId`, it discovers which EBS volumes it needs to attach and mount.

### How do I use it

Tag your EC2 instances:

```
Prefix: <logical-stack-name>
NodeId: <machine node #>
```

Tag your EBS volumes:

```
Prefix: <logical-stack-name>
VolumeId: <volume # corresponding 1:1 with node #>
DiskId: <disk # within volume group>
```

Run `kraken` from the EC2 instance to automatically mount the associated EBS volumes.

### Usecase

Use this if you want to provision X instances + volumes with Terraform.

Disks > 1 will be mounted in RAID (coming soon).

### Outputs

An early look at running outputs:

Before attach:

```
[dbadmin@ip-172-31-34-108 ~]$ ./kraken
kraken: cli.go:12: RUNNING KRAKEN: Tuesday, 09-May-17 16:03:20 UTC
kraken: aws.go:162: {
  AttachTime: 2017-05-09 16:03:24.393 +0000 UTC,
  Device: "/dev/xvde",
  InstanceId: "i-02802839b3fa11cb2",
  State: "attaching",
  VolumeId: "vol-041081d4f86a36fd2"
}
kraken: aws.go:162: {
  AttachTime: 2017-05-09 16:03:24.742 +0000 UTC,
  Device: "/dev/xvdc",
  InstanceId: "i-02802839b3fa11cb2",
  State: "attaching",
  VolumeId: "vol-0d9189f5a6c8c8a99"
}
kraken: aws.go:162: {
  AttachTime: 2017-05-09 16:03:25.058 +0000 UTC,
  Device: "/dev/xvdb",
  InstanceId: "i-02802839b3fa11cb2",
  State: "attaching",
  VolumeId: "vol-07af31509a1ebe8ab"
}
kraken: aws.go:162: {
  AttachTime: 2017-05-09 16:03:25.456 +0000 UTC,
  Device: "/dev/xvdd",
  InstanceId: "i-02802839b3fa11cb2",
  State: "attaching",
  VolumeId: "vol-001be1be9765a6cd1"
}
kraken: cli.go:18: Attached: [/dev/xvde /dev/xvdc /dev/xvdb /dev/xvdd
```

After attach:

```
[dbadmin@ip-172-31-34-108 ~]$ ./kraken
kraken: cli.go:12: RUNNING KRAKEN: Tuesday, 09-May-17 17:18:26 UTC
2017/05/09 13:18:26 Active attachments on volume vol-041081d4f86a36fd2, investigating...
2017/05/09 13:18:26 Active attachment is on current instance-id, continuing
2017/05/09 13:18:26 Active attachments on volume vol-0d9189f5a6c8c8a99, investigating...
2017/05/09 13:18:26 Active attachment is on current instance-id, continuing
2017/05/09 13:18:26 Active attachments on volume vol-07af31509a1ebe8ab, investigating...
2017/05/09 13:18:26 Active attachment is on current instance-id, continuing
2017/05/09 13:18:26 Active attachments on volume vol-001be1be9765a6cd1, investigating...
2017/05/09 13:18:26 Active attachment is on current instance-id, continuing
kraken: aws.go:168: Nothing to attach, returning existing attached device names
kraken: cli.go:18: Attached: [/dev/xvde /dev/xvdc /dev/xvdb /dev/xvdd]
```
