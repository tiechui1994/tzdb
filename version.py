import json
import re
import os
import requests

endpoint = "https://api.github.com"


def set_env(key, value: str):
    u = endpoint + "/repos/tiechui1994/tzdata/environments/%s" % (key,)
    header = {
        'Accept': 'application/vnd.github.v3+json'
    }
    body = {
        'wait_timer': 0,
    }
    response = requests.request("PUT", url=u, headers=header, data=body)
    if response.status_code == 200:
        result = json.JSONDecoder().decode(str(response.content, 'utf-8'))
        if len(result) > 0:
            return result[0].get('tag_name')

    return ''


def get_release_last_version() -> str:
    u = endpoint + "/repos/tiechui1994/tzdata/releases"
    header = {
        'Accept': 'application/vnd.github.v3+json'
    }
    params = {
        'per_page': 100,
        'page': 0,
    }
    response = requests.request("GET", url=u, headers=header, params=params)
    if response.status_code == 200:
        result = json.JSONDecoder().decode(str(response.content, 'utf-8'))
        if len(result) > 0:
            return result[0].get('tag_name')

    return ''


def get_tzdata_last_version() -> str:
    u = "https://www.iana.org/time-zones"
    header = {
        'Accept': 'text/html'
    }

    response = requests.request("GET", url=u, headers=header)
    if response.status_code == 200:
        ver = re.compile(r'<span id="version">(.*?)</span>', flags=re.MULTILINE)
        versons = ver.findall(str(response.content))
        if len(versons) > 0:
            return versons[0]
    return ''


if __name__ == '__main__':
    version = get_tzdata_last_version()
    release = get_release_last_version()

    if version == release:
        out = os.popen('bash build.sh %s' % (version,))
        print(str(out.read()))
        with open("/tmp/version", mode='w+') as fd:
            fd.write(version)
