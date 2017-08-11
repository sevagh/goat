resource "aws_network_interface" "goat_eni" {
  count = "${var.servers}"

  subnet_id = "${aws_subnet.goat_subnet.id}"

  tags {
    Name             = "${var.prefix}-network-interface-${count.index}"
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
