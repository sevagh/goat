resource "aws_ebs_volume" "data_disk" {
  count = "${var.servers * 3}"

  availability_zone = "us-east-1a"
  size              = "100"
  type              = "standard"

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

resource "aws_ebs_volume" "log_disk" {
  count = "${var.servers}"

  availability_zone = "us-east-1a"
  size              = "20"
  type              = "standard"

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

output "disk_id" {
  value = "${concat(list(aws_ebs_volume.data_disk.*.id), list(aws_ebs_volume.log_disk.*.id))}"
}

output "disk_name_tag" {
  value = "${concat(list(aws_ebs_volume.data_disk.*.tags), list(aws_ebs_volume.log_disk.*.tags))}"
}
