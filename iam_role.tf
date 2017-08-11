resource "aws_iam_role" "iam_role" {
  name = "${var.prefix}_iam_role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_instance_profile" "iam_profile" {
  name = "${var.prefix}_instance_profile"
  role = "${aws_iam_role.iam_role.name}"
}

resource "aws_iam_role_policy" "iam_role_policy" {
  name = "${var.prefix}_iam_role_policy"
  role = "${aws_iam_role.iam_role.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Action": [
      "ec2:*"
    ],
    "Resource": "*"
  }]
}
EOF
}
