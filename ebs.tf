resource "aws_ebs_volume" "data_disk" {
  count = "${var.servers * 2}"

  availability_zone = "us-east-1a"
  size              = "100"
  type              = "standard"

  tags {
    Name                 = "${var.prefix}-data-disk-${count.index}"
    "GOAT-IN:Prefix"     = "${var.prefix}"
    "GOAT-IN:NodeId"     = "${count.index / 2}"
    "GOAT-IN:VolumeName" = "data"
    "GOAT-IN:VolumeSize" = "2"
    "GOAT-IN:MountPath"  = "/krkn_data"
    "GOAT-IN:RaidLevel"  = "0"
    "GOAT-IN:FsType"     = "ext4"
  }
}

resource "aws_ebs_volume" "log_disk" {
  count = "${var.servers}"

  availability_zone = "us-east-1a"
  size              = "20"
  type              = "standard"

  tags {
    Name                 = "${var.prefix}-log-disk-${count.index}"
    "GOAT-IN:Prefix"     = "${var.prefix}"
    "GOAT-IN:NodeId"     = "${count.index}"
    "GOAT-IN:VolumeName" = "log"
    "GOAT-IN:VolumeSize" = "1"
    "GOAT-IN:MountPath"  = "/krkn_log"
    "GOAT-IN:RaidLevel"  = "0"                                     #ignored since volumesize == 1
    "GOAT-IN:FsType"     = "ext4"
  }
}

output "disk_id" {
  value = "${concat(list(aws_ebs_volume.data_disk.*.id), list(aws_ebs_volume.log_disk.*.id))}"
}

output "disk_name_tag" {
  value = "${concat(list(aws_ebs_volume.data_disk.*.tags), list(aws_ebs_volume.log_disk.*.tags))}"
}
