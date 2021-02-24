# DUNE-RC-RC

## Quick start

```bash
# build local base images
make docker.images
# start a cluster of said images
make docker.start
```

In another shell, run
```bash
# run ansible
# presently, you might have to run this twice. 'start daq-application-manager' can fail because nomad didn't initialize fast enough
make docker.ansible

# run the daq application example
configs/examples/working-daq-app/add.sh
# check the status of your app
configs/examples/working-daq-app/status.sh
```

You should see your daq app has been instrumented
```json
{
  "meta": {
    "kind": "daq-application",
    "name": "my-first-daq-app",
    "owner": ""
  },
  "spec": {
    "configkey": "/daq-applications/configs/my-first-config",
    "daq-service": "daq-ru-pu-52-51-api",
    "desired-state": "running",
    "enabled": true,
    "run-number": "123"
  },
  "status": {
    "commandpostfailed": false,
    "configkeyexists": true,
    "configrendered": true,
    "daqserviceexists": true
  }
}
```