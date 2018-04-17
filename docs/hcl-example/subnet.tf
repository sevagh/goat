resource "aws_subnet" "goat_subnet" {
  vpc_id                  = "${aws_vpc.goat_vpc.id}"
  cidr_block              = "10.0.1.0/24"
  availability_zone       = "${var.az}"
  map_public_ip_on_launch = "true"

  tags {
    Name = "goat-subnet"
  }
}

resource "aws_route_table_association" "goat_route_table_association" {
  subnet_id      = "${aws_subnet.goat_subnet.id}"
  route_table_id = "${aws_route_table.goat_route_table.id}"
}
