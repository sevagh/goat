data "template_file" "bootstrap" {
  count    = "${var.servers}"
  template = "${file("bootstrap.tpl")}"
}

resource "aws_instance" "instance" {
  count                       = "${var.servers}"
  availability_zone           = "us-east-1a"
  ami                         = "ami-6d1c2007"
  instance_type               = "t2.micro"
  key_name                    = "${var.keypair_name}"
  associate_public_ip_address = true
  user_data                   = "${element(data.template_file.bootstrap.*.rendered, count.index)}"
  iam_instance_profile        = "${aws_iam_instance_profile.iam_profile.id}"

  tags {
    Name = "${var.prefix}-${count.index}"
    "KRKN-IN:Prefix" = "${var.prefix}"
    "KRKN-IN:NodeId" = "${count.index}"
  }
}

output "instance_public_ip" {
  value = "${list(aws_instance.instance.*.public_ip)}"
}

output "instance_name_tag" {
  value = "${list(aws_instance.instance.*.tags)}"
}
