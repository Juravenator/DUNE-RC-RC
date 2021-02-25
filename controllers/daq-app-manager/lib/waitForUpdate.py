import requests
import time

import logging
log = logging.getLogger('waitForUpdate')

def waitForUpdate(lastChangeID, rerunTimeOverride):
  log.info("wating for updates since %s", lastChangeID)
  url = 'http://localhost:8500/v1/kv/daq-applications/?keys=&separator=/'
  if lastChangeID is not None and rerunTimeOverride is None:
    url += '&index=' + str(lastChangeID)
  if rerunTimeOverride is not None:
    time.sleep(rerunTimeOverride)
#   log.info("url is %s", url)
  resp = requests.get(url)
  if resp.status_code != 200:
      log.error("update fetch failed with code %s", resp.status_code)
      # This means something went wrong.
      return None, None, Exception('GET daq_applications failed with code {}'.format(resp.status_code))
  log.info("raft update received", extra={'resp': resp, 'headers': resp.headers, 'body': resp.json()})
  return [key for key in resp.json() if not key.endswith("/")], resp.headers['X-Consul-Index'], None