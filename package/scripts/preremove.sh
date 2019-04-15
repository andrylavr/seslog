#!/bin/sh
set -e
if [ "$1" = remove ]; then
    /bin/systemctl stop    seslog-server
    /bin/systemctl disable seslog-server
    /bin/systemctl daemon-reload
fi