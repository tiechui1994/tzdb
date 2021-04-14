#!/usr/bin/env bash

url=$(curl \
  -H "Accept: application/vnd.github.v3+json" \
  https://api.github.com/repos/tiechui1994/tzdb/releases/latest -s | \
  grep -E -o '"browser_download_url": ".*?[^"]' | \
  grep -E -o 'https.*')

echo "url: $url"

wget --continue --quiet ${url} -O /tmp/zoneinfo.zip

rm -rf /tmp/zoneinfo && mkdir /tmp/zoneinfo \
unzip -d /tmp/zoneinfo /tmp/zoneinfo.zip

mv /tmp/zoneinfo /usr/share/zoneinfo