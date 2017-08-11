data "template_file" "bootstrap" {
  count    = "${var.servers}"
  template = "${file("bootstrap.tpl")}"
}

resource "aws_instance" "instance" {
  count                  = "${var.servers}"
  availability_zone      = "us-east-1a"
  ami                    = "ami-6d1c2007"
  instance_type          = "t2.micro"
  key_name               = "${var.keypair_name}"
  user_data              = "${element(data.template_file.bootstrap.*.rendered, count.index)}"
  iam_instance_profile   = "${aws_iam_instance_profile.iam_profile.id}"
  subnet_id              = "${aws_subnet.goat_subnet.id}"
  vpc_security_group_ids = ["${aws_security_group.goat_sg.id}"]

  tags {
    Name             = "${var.prefix}-${count.index}"
    "GOAT-IN:Prefix" = "${var.prefix}"
    "GOAT-IN:NodeId" = "${count.index}"
  }
}

output "instance_public_ip" {
  value = "${list(aws_instance.instance.*.public_ip)}"
}

output "instance_name_tag" {
  value = "${list(aws_instance.instance.*.tags)}"
}
