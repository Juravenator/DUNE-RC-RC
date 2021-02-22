#!/usr/bin/env python3

# import os
# os.chdir(os.path.dirname(os.path.abspath(__file__)))

import asyncio
import logging
import time
from lib import waitForUpdate, manage
logging.basicConfig(level=logging.INFO)
log = logging.getLogger('main')

async def main():
  lastChangeID = 0
  while True:
    daqAppKeys, changeID, err = waitForUpdate(lastChangeID)
    if err is not None:
      time.sleep(5)
      continue
    # waitForUpdate uses long polling
    # if the ID is not changed, it means the 
    # poll simply timed out
    if changeID == lastChangeID:
      log.info("ignoring this update (%s)", changeID)
      continue
    log.info("received new change %s: %s", changeID, daqAppKeys)

    requeueAfter = None
    for key in daqAppKeys:
      requestedRedos = await asyncio.gather(*[manage(daqAppKey) for daqAppKey in daqAppKeys])
      for r in requestedRedos:
        if r is not None:
          if requeueAfter is None or r < requeueAfter:
            requeueAfter = r

    if requeueAfter is None:
      lastChangeID = changeID
      log.info("all daq apps are up to date (%s)", changeID)
    else:
      log.info("not all daq apps ready, rerunningin %ds", requeueAfter)

log.info("starting")
loop = asyncio.get_event_loop()
loop.run_until_complete(asyncio.wait([main()]))