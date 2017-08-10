#!/bin/bash

set -o errtrace
set -o nounset
set -o pipefail
set -o errexit

yum install -y wget mdadm
wget https://github.com/sevagh/goat/releases/download/0.2.0/goat
chmod +x goat
./goat >/var/log/goat.log 2>&1

exit 0
