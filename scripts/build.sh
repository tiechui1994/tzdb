#!/usr/bin/env bash

set -e
version="$1"

dir="$(pwd)"
source="$dir/source"
build="$dir/build"
install="$dir/zone"

if [[ -z ${GITHUB_WORKSPACE} ]]; then
    GITHUB_WORKSPACE=$(pwd)
fi


## download
rm -rf ${source} && mkdir -p ${source} && cd ${source}
prefix="https://data.iana.org/time-zones/releases"
wget "$prefix/tzdata$version.tar.gz"
wget "$prefix/tzcode$version.tar.gz"


## code
rm -rf ${build} && mkdir -p ${build} && cd ${build}
rm -rf tzdata && mkdir -p tzdata
cd tzdata

tar xvf "$source/tzdata$version.tar.gz" -C .
tar xvf "$source/tzcode$version.tar.gz" -C .

chmod -Rf a+rX,u+w,g-w,o-w .

make VERSION="$version" "tzdata$version-rearguard.tar.gz"

tar zxf "tzdata$version-rearguard.tar.gz"

rm tzdata.zi
make VERSION="$version" DATAFORM=rearguard tzdata.zi


## build
cd ${build}/tzdata

FILES="africa antarctica asia australasia europe northamerica southamerica pacificnew etcetera backward"
mkdir -p zoneinfo/posix && mkdir -p zoneinfo/right
zic -y ./yearistype -d zoneinfo -L /dev/null -p America/New_York ${FILES}
zic -y ./yearistype -d zoneinfo/posix -L /dev/null  ${FILES}
zic -y ./yearistype -d zoneinfo/right -L leapseconds  ${FILES}


## install
cd ${build}/tzdata
install -d ${install}
cp -prd zoneinfo ${install}
install -p -m 644 zone.tab zone1970.tab iso3166.tab leapseconds tzdata.zi ${install}/zoneinfo
cd ${install}/zoneinfo && zip -r -0 zoneinfo.zip *
mv zoneinfo.zip ${GITHUB_WORKSPACE}/${version}_zoneinfo.zip
