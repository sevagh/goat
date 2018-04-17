#!/bin/bash

set -o errtrace
set -o nounset
set -o pipefail
set -o errexit

yum install -y wget mdadm
yum install -y https://github.com/sevagh/goat/releases/download/0.4.0/goat-0.4.0-1.fc25.x86_64.rpm
yum install -y https://github.com/sevagh/ec2-utils/releases/download/v0.5.3/ec2-net-utils-0.5-2.fc25.noarch.rpm
systemctl enable elastic-network-interfaces
systemctl start elastic-network-interfaces
systemctl enable goat@ebs
systemctl enable goat@eni
systemctl start goat@ebs
systemctl start goat@eni

exit 0
