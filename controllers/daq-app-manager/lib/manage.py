import logging
import requests
import json
from jsonschema import validate
from .runConsulTemplate import runConsulTemplate
from .getConfigTemplate import getConfigTemplate
log = logging.getLogger('manage')

schemafile = open('daq-application.schema.json')
schema = json.load(schemafile)

def serializeDoc(doc):
  return json.dumps(doc, sort_keys=True)

async def manage(daqAppKey):
  log.info("managing %s", daqAppKey)
  resp = requests.get('http://localhost:8500/v1/kv/' + daqAppKey + "?raw=")
  if resp.status_code != 200:
      log.error("DAQ App fetch failed with code %s", resp.status_code)
      # don't re-loop for a shitty key
      return None

  try:
    doc = resp.json()
    validate(instance=doc, schema=schema)
    startDoc = serializeDoc(doc)
  except Exception:
    log.error("DAQ App configuration is not valid json", exc_info=True)
    # don't re-loop for a shitty key
    return None

  if "status" not in doc:
    doc['status'] = {}
  # starter pack of status flags
  doc['status']['configkeyexists'] = False
  doc['status']['configrendered'] = False

  configkey = doc['spec']['configkey']
  configTmpl = getConfigTemplate(configkey)
  configkeyexists = configTmpl is not None
  doc['status']['configkeyexists'] = configkeyexists
  if not configkeyexists:
    log.error("specified config key does not exist: %s", configkey)
    return None
  else:
    rendered, err = runConsulTemplate(configTmpl)
    if err is not None:
      log.error("template parsing failed: %s", err)
    else:
      log.info("template rendered")
      doc['status']['configrendered'] = True

  endDoc = serializeDoc(doc)
  if startDoc != endDoc:
    log.info("pushing new status changes", doc['status'])
    resp = requests.put('http://localhost:8500/v1/kv/' + daqAppKey, data=endDoc)
    if resp.status_code != 200:
        log.error("DAQ App update failed with code %s", resp.status_code)
        # updating status is important, re-queue with a timeout
        return 10

  log.info("content of daq key %s: %s", daqAppKey, doc['meta']['name'])
  return None