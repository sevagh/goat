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
