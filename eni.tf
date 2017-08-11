resource "aws_network_interface" "goat_eni" {
  count = "${var.servers}"

  availability_zone = "us-east-1a"
  size              = "100"
  type              = "standard"

  tags {
    Name = "${var.prefix}-network-interface-${count.index}"
    "GOAT-IN:Prefix" = "${var.prefix}"
    "GOAT-IN:NodeId" = "${count.index}" 
  }
}

output "eni_id" {
  value = "${list(aws_network_interface.goat_eni.*.id)}"
}

output "eni_name_tag" {
  value = "${list(aws_network_interface.goat_eni.*.tags)}"
}
