#!/usr/bin/env bash

# Versions to use.
CODE=2021a
DATA=2021a

set -e
rm -rf work && mkdir work && cd work && mkdir zoneinfo
curl -L -O https://www.iana.org/time-zones/repository/releases/tzcode${CODE}.tar.gz
curl -L -O https://www.iana.org/time-zones/repository/releases/tzdata${DATA}.tar.gz
tar xzf tzcode${CODE}.tar.gz
tar xzf tzdata${DATA}.tar.gz

make CFLAGS=-DSTD_INSPIRED AWK=awk TZDIR=zoneinfo posix_only

cd zoneinfo && zip -0 -r zoneinfo.zip *
cd ../../

echo
if [[ "$1" = "-work" ]]; then
	echo Left workspace behind in work/.
else
	rm -rf work
fi

echo New time zone files in zoneinfo.zip.