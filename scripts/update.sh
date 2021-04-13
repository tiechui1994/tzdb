#!/usr/bin/env bash

function json() {
  local json=$1
  local key=$2

  if [[ -z "$3" ]]; then
    local num=1
  else
    local num=$3
  fi

  local value=$(echo "${json}" | awk -F "[,:}]" '{
    for(i=1;i<=NF;i++) {
      if($i~/'${key}'\042/) {
        print $(i+1)
      }
    }
  }' | tr -d '"' | sed -n ${num}p)

  echo ${value}
}

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