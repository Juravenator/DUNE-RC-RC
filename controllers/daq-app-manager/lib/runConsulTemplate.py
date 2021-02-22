import logging
import subprocess
import sys
from subprocess import Popen, PIPE, STDOUT, TimeoutExpired
log = logging.getLogger('runConsulTemplate')

def runConsulTemplate(input):
  p = Popen(["./consul-filler/build/consul-filler"], stdout=PIPE, stdin=PIPE, stderr=sys.stderr)
  try:
    stdout_data = p.communicate(input=str.encode(input), timeout=15)[0]
    if p.returncode != 0:
      return None, Exception("template parsing failed")
    return stdout_data, None
  except TimeoutExpired:
    p.kill()
    p.communicate()

  return None, Exception("something went wrong")

