import json
import os
import re
import requests

endpoint = "https://api.github.com"


def get_release_last_version() -> str:
    u = endpoint + "/repos/tiechui1994/tzdb/releases/latest"
    header = {
        'Accept': 'application/vnd.github.v3+json'
    }
    print(u)
    response = requests.request("GET", url=u, headers=header, timeout=120)
    print(response.status_code)
    if response.status_code == 200:
        result = json.JSONDecoder().decode(str(response.content, 'utf-8'))
        return result.get('tag_name')

    return ''


def get_tz_last_version() -> str:
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
    version = get_tz_last_version()
    release = get_release_last_version()
    content = ''
    if version != release:
        out = os.popen('bash build.sh %s' % (version,))
        print(str(out.read()))
        content = version

    with open("/tmp/version", mode='w+') as fd:
        fd.write(content)
