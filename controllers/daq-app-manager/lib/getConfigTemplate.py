import requests

import logging
log = logging.getLogger('getConfigTemplate')

def getConfigTemplate(key):
  resp = requests.get('http://localhost:8500/v1/kv/' + key + '?raw=')
  if resp.status_code != 200:
    return None

  return resp.json()["spec"]["template"]