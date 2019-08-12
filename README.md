# goat

**N.B.: the active maintenance fork of this project is available at [https://github.com/steamhaus/goat](https://github.com/steamhaus/goat)**

### Attach EBS volumes and ENIs to running EC2 instances

`goat` is a Go program which runs from inside the EC2 instance.

By setting your tags correctly, `goat` can discover and attach EBS volumes and ENIs. For EBS volumes, it can perform additional actions such as RAID (with mdadm), mkfs, and mount EBS volumes to the EC2 instance where it's running.

### Install and run

Goat is a Go binary that should be able to run on any Linux instance. In the releases tab you can find a zip of the binary, and a `.deb` and `.rpm` package with systemd support. Goat needs `mdadm` to perform RAID (which is a dependency in the deb and rpm).

To use goat, run it during the launch process of your EC2 instance - you can `systemctl enable goat@TARGET` and `systemctl start goat@TARGET` (where TARGET is one of ebs or eni) in the EC2 user data script. [Full Terraform example here](./terraform-example).

### Usage

In the most basic case, you should run `goat ebs` or `goat eni`.

Full usage:

```
Usage: goat [OPTIONS] ebs|eni

OPTIONS
  -debug
        Interactive debug prompts
  -logLevel string
        Log level (default "info")
  -tagPrefix string
        Prefix for GOAT related tags (default "GOAT-IN")
  -version
        Display version and exit
```

You can set `-tagPrefix` and `-logLevel` with environment variables (which take precedence):

* `GOAT_LOG_LEVEL`
* `GOAT_TAG_PREFIX`

### Tags

These are the tags you need (recall that the `GOAT-IN` prefix is configurable):

| Tag Name             | Description             | Resource type     | Required            | Effect      |
| :------------------- | :---------------------- | :---------------- | :------------------ | :---------- |
| GOAT-IN:Prefix       | Logical app name        | EC2, EBS, ENI     | :heavy_check_mark:  | attach      |
| GOAT-IN:NodeId       | Node id                 | EC2, EBS, ENI     | :heavy_check_mark:  | attach      |
| GOAT-IN:VolumeName   | Distinct volume name    | EBS               |                     |             |
| GOAT-IN:VolumeSize   | # of disks in vol group | EBS               |                     | mdadm       |
| GOAT-IN:RaidLevel    | level of RAID (0 or 1)  | EBS               |                     | mdadm       |
| GOAT-IN:MountPath    | Linux path to mount vol | EBS               |                     | mount       |
| GOAT-IN:FsType       | Linux filesystem type   | EBS               |                     | mkfs        |

If non-required tags are omitted, that step is skipped. The barest case will simply attach the EBS volumes and perform no further actions.

The filesystem and fstab entries are created with the label `GOAT-{VolumeName}` for convenience. Running `goat` multiple times will result in it detecting the existing label it intended to create and not proceeding.

Aside from the `mount` syscall, `goat` shells out to `mdadm`, `blkid`, and `mkfs`. If the mount and RAID steps are performed, the configs will be persisted to `/etc/fstab` and `/etc/mdadm.conf`.

Check the [Terraform example](./terraform-example) for example tag values.

### Permissions

It's necessary for the instance to have an IAM Role with _at least_ access to the EBS and ENI resources that it will be attaching - see [here](./terraform-example/iam_role.tf). Your roles can be even more permissive (i.e. full EC2 access) but that comes with its own risks.

Unfortunately, resource-level permissions are [currently not supported](https://docs.aws.amazon.com/AWSEC2/latest/APIReference/ec2-api-permissions.html#ec2-api-unsupported-resource-permissions) for attaching network interfaces. This means that to use `goat@eni`, your instances must have full permissions for __all__ ENIs.

### Example EBS usecase - attaching old disks to a new instance

The specific use-case that `goat` was developed to solve is the following. Say you have 3 instances with their own respective disks, and you receive a termination notice for instance 1. I want the `goat` workflow to be:

1. Terminate instance 1
2. Create a new instance with the same GOAT tags (to indicate to `goat` that it's the logical equivalent as the machine it is replacing)
   1. **No need to modify or manipulate the EBS volumes or their tags**
4. On boot, everything works magically

After `goat` ran on the first fresh run, the EBS volumes got the correct filesystems, labels, and in the case of RAID, `mdadm` metadata on them.

The event flow on a re-created instance is:

1. Get EC2 metadata on the running instance, create an EC2 client, and search EBS volumes
2. Attach the volumes it needs based on their tags
3. Discover that `/dev/disk/by-label` already contains the correct disks
    1. From `mdadm` magic, after the EBS attachment the RAID array is already detected correctly
4. Proceed to perform the `fstab` and `mount` phases - skip `mdadm`, `mkfs`

**CAVEAT**: the mdadm metadata will have the hostname of the previous EC2 instance:

```
[centos@ip-172-31-29-69 ~]$ sudo mdadm --detail --scan --verbose
ARRAY /dev/md127 level=raid0 num-devices=3 metadata=1.2 name="ip-172-31-25-105:'GOAT-data'" UUID=2d08b310:fd13bd21:bc2417a4:56a1ec57
   devices=/dev/xvdb,/dev/xvdc,/dev/xvdd
[centos@ip-172-31-29-69 ~]$
```

To avoid this, define a good/persistent hostname for the EC2 instance, that you will then re-apply to any instance taking ownership of the previous instance's disks.

### ENI notes

As mentioned, the only action `goat` will do for ENIs is attaching them. You can try to use [ec2-net-utils](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-eni.html#ec2-net-utils), a tool available on Amazon Linux AMIs, or [this port to CentOS/systemd](https://github.com/etuttle/ec2-utils), to configure an ENI after `goat` attaches it.

ENI attachments take a parameter called [DeviceIndex](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-network-interface-attachment.html). Goat isn't smart, and always starts from DeviceIndex `1`. This means that your EC2 instance should have no attached ENIs to use `goat`. If it does, they should be the ones that `goat` was going to attach anyway, not external ENIs that have no `goat` tags.

### Build and develop

The deb, rpm, and zip are generated from a multi-stage [Dockerfile.build](./Dockerfile.build). Invoke it with `make docker-build`. If you have docker locally, you can use the following command in order to quickly get a development env ready: `make dev-env`.
