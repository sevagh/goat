#!/bin/bash

set -o errtrace
set -o nounset
set -o pipefail
set -o errexit

yum install -y wget mdadm
yum install -y https://github.com/sevagh/goat/releases/download/0.4.0/goat-0.4.0-1.fc25.x86_64.rpm
systemctl enable goat@ebs
systemctl enable goat@eni
systemctl start goat@ebs
systemctl start goat@eni

exit 0
