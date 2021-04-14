#!/usr/bin/env bash

# Versions to use.
TAG="$1"

set -e
if [[ -z ${GITHUB_WORKSPACE} ]]; then
    GITHUB_WORKSPACE=$(pwd)
fi

rm -rf work && mkdir work && cd work
git clone --depth 1 --branch ${TAG} git@github.com:eggert/tz.git
cd tz && make TOPDIR=data install
cd data/etc/zoneinfo && zip -0 -r zoneinfo.zip * && mv zoneinfo.zip ${GITHUB_WORKSPACE}/${TAG}_zoneinfo.zip
cd ${GITHUB_WORKSPACE}

echo
if [[ "$1" = "-work" ]]; then
	echo Left workspace behind in work/.
else
	rm -rf work
fi

echo "New time zone files in ${GITHUB_WORKSPACE}/${TAG}_zoneinfo.zip"