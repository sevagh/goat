#!/bin/bash

set -o errtrace
set -o nounset
set -o pipefail
set -o errexit

yum install -y wget mdadm
wget https://github.com/sevagh/kraken/releases/download/0.1.0/kraken
chmod +x kraken
./kraken >/var/log/kraken.log 2>&1

exit 0
