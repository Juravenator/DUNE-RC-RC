import logging
import requests
import json
import socket
import threading
import time
from jsonschema import validate
from http.server import BaseHTTPRequestHandler, HTTPServer
from .runConsulTemplate import runConsulTemplate
from .getConfigTemplate import getConfigTemplate
log = logging.getLogger('manage')

schemafile = open('daq-application.schema.json')
schema = json.load(schemafile)

async def manage(daqAppKey, doc):
  if "status" not in doc:
    doc['status'] = {}

  # does the daq app exist?
  daqservice = getDAQServiceAddr(doc['spec']['daq-service'])
  doc['status']['daqserviceexists'] = daqservice is not None
  if daqservice is not None:
    log.info("service is at %s", daqservice)

  # does the config key exist?
  configkey = doc['spec']['configkey']
  configTmpl = getConfigTemplate(configkey)
  configkeyexists = configTmpl is not None
  doc['status']['configkeyexists'] = configkeyexists
  if not configkeyexists:
    log.error("specified config key does not exist: %s", configkey)
    return None

  # is the config valid?
  rendered, err = runConsulTemplate(configTmpl)
  doc['status']['configrendered'] = err is None
  if err is not None:
    log.error("template parsing failed: %s", err)
  else:
    log.info("template rendered")

  # do we have the basics to do anything at all?
  if not doc['status']['configkeyexists'] or not doc['status']['configrendered'] or not doc['status']['daqserviceexists']:
    log.error("daq app %s is not healthy enough to receive commands", daqAppKey)
    if not doc['status']['daqserviceexists']:
      # try again later when daq app might be online
      return 5
    else:
      # invalid config, wait for config update
      return None

  if "enabled" in doc["spec"] and not doc["spec"]["enabled"]:
    log.warn("daq app %s is not in autonomous mode", daqAppKey)
    return None

  # are we waiting for a command reply?
  if "lastcommandsent" in doc['status'] and "commandsucceeded" not in doc['status'] and "commandpostfailed" not in doc['status']:
    log.info("waiting for last command to reply...")
    # there is a separate thread that will update our status, no need for a re-trigger
    return None

  # did we receive a reply but did it fail?
  if "commandsucceeded" in doc['status'] and not doc['status']['commandsucceeded']:
    # retry if it timed out
    if "commandtimedout" in doc['status'] and doc['status']['commandtimedout']:
      doc['status'].pop('lastcommandsent', None)
    else:
      # we're stuck
      log.error("daq app stuck %s", daqAppKey)
      return None

  # are we already in desired state?
  if "commandsucceeded" in doc['status'] and doc['status']['commandsucceeded']:
    lc = doc['status']['lastcommandsent']
    ds = doc['spec']['desired-state']
    amdone = False
    if ds == "running" and lc == "start":
      amdone = True
    if ds == "configured" and lc == "conf":
      amdone = True
    if ds == "init" and lc == "init":
      amdone = True
    if amdone:
      log.info("daq application already in desired state '%s'", ds)
      return None

  # determine command to send
  if "lastcommandsent" not in doc['status']:
    nextcommand = getNextCommandName(None)
  else:
    nextcommand = getNextCommandName(doc['status']['lastcommandsent'])

  # send command
  payload = next((cmddata for cmddata in rendered if cmddata["id"] == nextcommand), None)
  if payload is None:
    log.error("cannot determine payload for command %s", nextcommand)
    return None
  
  try:
    p = get_free_tcp_port()
    log.info("setting up response listener on %d", p)
    ResponseListener(p, daqAppKey)
    log.info("sending command %s to %s", nextcommand, daqservice)
    doc['status']['lastcommandsent'] = nextcommand
    doc['status'].pop('commandsucceeded', None)
    resp = requests.post("http://" + daqservice + "/command", json=payload, headers = {"X-Answer-Port": str(p)})
    if not (200 <= resp.status_code <= 299):
      log.error("DAQ App command post failed with code %s", resp.status_code)
      doc['status']['commandpostfailed'] = True
      return 5
    else:
      doc['status'].pop('commandpostfailed', None)

  except Exception:
    log.exception("next command could not be sent")
    doc['status']['commandpostfailed'] = True
    return 5

  return None

def get_free_tcp_port():
  tcp = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
  tcp.bind(('', 0))
  addr, port = tcp.getsockname()
  tcp.close()
  return port

class ResponseHandler(BaseHTTPRequestHandler):
  def __init__(self, socket, b, httpserver):
    self.httpserver = httpserver
    super().__init__(socket, b, httpserver)
  
  def do_POST(self):
    try:
      content_length = int(self.headers['Content-Length'])
      raw_data = self.rfile.read(content_length)
      data = json.loads(raw_data)
    except Exception:
      self.send_error(400, "invalid body")
      return

    log.info("received response '%s' for command '%s'", data['result'], data['command'])
    self.httpserver.response = data['result']

    self.send_response(200)
    self.send_header('Content-type', 'text/plain')
    self.end_headers()
    self.wfile.write("Response received".encode('utf-8'))

class ResponseServer(HTTPServer):
  def __init__(self, a, b):
    self.response = ""
    super().__init__(a, b)

class ResponseListener(object):
  def __init__(self, port, key):
    self.response = ""
    self.key = key
    self.thread = threading.Thread(target=self.run, args=(port,))
    self.thread.daemon = True
    self.thread.start()

  def run(self, port):
    log.info("listening for response on port %d", port)
    httpd = ResponseServer(('', port), ResponseHandler)
    httpd.timeout = 10
    log.info("listening timeout is %ds", httpd.timeout)

    httpd.handle_request()
    self.response = httpd.response

    log.info("response received for command to '%s': '%s'", self.key, self.response)
    doc = getDoc(self.key)
    if doc is None:
      log.error("cannot update status for %s after command response", self.key)
      return

    doc['status']['commandsucceeded'] = self.response == "OK"
    doc['status']['commandtimedout'] = self.response == ""
    postDoc(self.key, serializeDoc(doc))

def getDoc(key):
  resp = requests.get('http://localhost:8500/v1/kv/' + key + "?raw=")
  if resp.status_code != 200:
    log.error("DAQ App fetch failed with code %s", resp.status_code)
    return None
  try:
    doc = resp.json()
    validate(instance=doc, schema=schema)
    return doc
  except Exception:
    log.error("DAQ App configuration is not valid json", exc_info=True)
    return None

def serializeDoc(doc):
  return json.dumps(doc, sort_keys=True)

def getDAQServiceAddr(name):
  resp = requests.get('http://localhost:8500/v1/health/service/' + name + "?dc=dune-rc")
  if resp.status_code != 200:
    log.error("DAQ App service fetch failed with code %s", resp.status_code)
    return None

  service = resp.json()
  if len(service) == 0:
    log.error("DAQ service does not exist: %s", name)
    return None

  return service[0]['Service']['Address'] + ":" + str(service[0]['Service']['Port'])

async def sendCommand(daqservice, payload, port):
  resp = requests.post("http://" + daqservice + "/command", json=payload, headers = {"X-Answer-Port": str(port)})
  if resp.status_code != 200:
    log.error("DAQ App command post failed with code %s", resp.status_code)
    return None, Exception("sending of command failed")
  return 

def getNextCommandName(lastCommand):
  if lastCommand is None or lastCommand == "":
    return "init"
  elif lastCommand == "init":
    return "conf"
  elif lastCommand == "conf":
    return "start"
  else:
    return None

async def manageAndUpdate(daqAppKey):
  log.info("managing %s", daqAppKey)
  doc = getDoc(daqAppKey)
  if doc is None:
    # don't re-loop for a shitty key
    return None

  startDoc = serializeDoc(doc)
  requeueTime = await manage(daqAppKey, doc)
  endDoc = serializeDoc(doc)

  if startDoc != endDoc:
    err = postDoc(daqAppKey, endDoc)
    if err is not None:
      # updating status is important, re-queue with a timeout
      return 10
  return requeueTime

def postDoc(key, doc):
  log.info("pushing new status changes %s -> %s", key, doc)
  resp = requests.put('http://localhost:8500/v1/kv/' + key, data=doc)
  if resp.status_code != 200:
    log.error("DAQ App update failed with code %s", resp.status_code)
    return Exception("app update failed")
  return None
