#!/usr/bin/env bash

# Versions to use.
CODE="$1"
DATA="$1"

set -e
if [[ -z ${GITHUB_WORKSPACE} ]]; then
    GITHUB_WORKSPACE=$(pwd)
fi

rm -rf work && mkdir work && cd work && mkdir zoneinfo
curl -L -O https://www.iana.org/time-zones/repository/releases/tzcode${CODE}.tar.gz
curl -L -O https://www.iana.org/time-zones/repository/releases/tzdata${DATA}.tar.gz
tar xzf tzcode${CODE}.tar.gz
tar xzf tzdata${DATA}.tar.gz

make CFLAGS=-DSTD_INSPIRED AWK=awk TZDIR=zoneinfo posix_only

cd zoneinfo && zip -0 -r zoneinfo.zip * && mv zoneinfo.zip ${GITHUB_WORKSPACE}/${CODE}_zoneinfo.zip
cd ${GITHUB_WORKSPACE}

echo
if [[ "$1" = "-work" ]]; then
	echo Left workspace behind in work/.
else
	rm -rf work
fi

echo "New time zone files in ${GITHUB_WORKSPACE}/${CODE}_zoneinfo.zip"