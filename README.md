# kraken

### Attach & mount EBS volumes from inside EC2 instances

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

Kraken runs from inside the EC2 instance (it's necessary for the instance to have an IAM Role with full EC2 access). 

By associating `instance tag.Prefix <-> volume tag.Prefix` and `instance tag.NodeId <-> volume tag.VolumeId`, it discovers which EBS volumes it needs to attach and mount.

Usecase:

Use this if you want to provision X instances + volumes with Terraform.

Disks > 1 will be mounted in RAID (coming soon).