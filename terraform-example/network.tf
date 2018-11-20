resource "aws_vpc" "goat_vpc" {
  cidr_block           = "10.0.0.0/16"
  instance_tenancy     = "default"
  enable_dns_hostnames = "true"

  tags {
    Name = "goat-vpc"
  }
}

resource "aws_internet_gateway" "goat_gateway" {
  vpc_id = "${aws_vpc.goat_vpc.id}"

  tags {
    Name = "goat-gateway"
  }
}

resource "aws_route_table" "goat_route_table" {
  vpc_id = "${aws_vpc.goat_vpc.id}"

  tags {
    Name = "goat-route-table"
  }
}

resource "aws_route" "goat_route_table_conf" {
  route_table_id         = "${aws_route_table.goat_route_table.id}"
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = "${aws_internet_gateway.goat_gateway.id}"
}

resource "aws_network_interface" "goat_eni" {
  count = "${var.servers}"

  subnet_id       = "${aws_subnet.goat_subnet.id}"
  security_groups = ["${aws_security_group.goat_sg.id}"]

  tags {
    Name             = "${var.prefix}-network-interface-${count.index}"
    "GOAT-IN:Prefix" = "${var.prefix}"
    "GOAT-IN:NodeId" = "${count.index}"
  }
}

resource "aws_eip" "goat_eip" {
  vpc               = true
  count             = "${var.servers}"
  network_interface = "${element(aws_network_interface.goat_eni.*.id, count.index)}"
}

output "aws_eip_public_ip" {
  value = "${list(aws_eip.goat_eip.*.public_ip)}"
}

output "eni_id" {
  value = "${list(aws_network_interface.goat_eni.*.id)}"
}

output "eni_name_tag" {
  value = "${list(aws_network_interface.goat_eni.*.tags)}"
}
