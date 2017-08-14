# example

Terraform recipe for a minimally viable example of the [goat](https://github.com/sevagh/goat) EBS/EC2 tag-based mounting system.

### Iterating

To iterate with this Terraform recipe, it's helpful to export the 3 required variables:

```
$ export TF_VAR_aws_access_key=xxxx
$ export TF_VAR_aws_secret_key=xxxx
$ export TF_VAR_keypair_name=mykeypair
```

### Results for this repo

#### EBS

```
[centos@ip-10-0-1-105 ~]$ df -h
Filesystem      Size  Used Avail Use% Mounted on
/dev/xvda1      8.0G 1009M  7.1G  13% /
devtmpfs        478M     0  478M   0% /dev
tmpfs           496M     0  496M   0% /dev/shm
tmpfs           496M   13M  484M   3% /run
tmpfs           496M     0  496M   0% /sys/fs/cgroup
tmpfs           100M     0  100M   0% /run/user/1000
/dev/md0        197G   61M  187G   1% /goat_data
/dev/xvdd        20G   45M   19G   1% /goat_log
[centos@ip-10-0-1-105 ~]$ ls /dev/disk/by-label/
GOAT-data  GOAT-log
```

#### ENI

```
[centos@ip-10-0-1-105 ~]$ ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host
       valid_lft forever preferred_lft forever
2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 9001 qdisc pfifo_fast state UP qlen 1000
    link/ether 0e:e8:f6:66:e3:4a brd ff:ff:ff:ff:ff:ff
    inet 10.0.1.105/24 brd 10.0.1.255 scope global dynamic eth0
       valid_lft 3489sec preferred_lft 3489sec
    inet6 fe80::ce8:f6ff:fe66:e34a/64 scope link
       valid_lft forever preferred_lft forever
3: eth1: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 9001 qdisc pfifo_fast state UP qlen 1000
    link/ether 0e:ce:a4:f0:be:f0 brd ff:ff:ff:ff:ff:ff
    inet 10.0.1.13/24 brd 10.0.1.255 scope global dynamic eth1
       valid_lft 3589sec preferred_lft 3589sec
    inet6 fe80::cce:a4ff:fef0:bef0/64 scope link
       valid_lft forever preferred_lft forever
```
