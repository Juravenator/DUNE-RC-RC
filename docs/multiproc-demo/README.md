# Running the march demo on LXPlus

## LXPlus version (single host)
```bash
$ ssh lxplus #duh
$ git clone https://github.com/Juravenator/DUNE-RC-RC.git
$ cd DUNE-RC-RC
$ git checkout march-rc0
```

```bash
$ hacks/lxplus/runInMem.sh
$ export PATH="$(pwd)/cli/build:$PATH"
$ docs/multiproc-demo/build-daq.sh
$ run-control apply docs/multiproc-demo/configs/*.json
$ run-control get all
$ run-control daq command init all
$ run-control daq command conf all
$ run-control daq command --run-number 42 start all
$ run-control daq command resume daq-app-multiproc-emu
$ run-control daq command stop all
$ hacks/lxplus/stopInMem.sh
```

## Docker version (multi host)
This will not work on LXPlus, you don't have access to docker.  
You could make it work on your local machine if you [accept a really ugly hack](../../hacks/rsync-cvmfs.md)

```bash
$ git clone https://github.com/Juravenator/DUNE-RC-RC.git
$ cd DUNE-RC-RC
$ git checkout march-rc0
```

Run the docker multi-host setup
```bash
$ make docker.images
$ make docker.start
# in another shell
$ make docker.ansible
```

Setup build environments.  
To save you time, I'll let you cheat and only run these long steps on the two instances we will run daq apps on today.
```bash
$ docker exec -it docker_dune-rc-march-ru-1_1 bash
$ /dune-rc/docs/multiproc-demo/build-daq.sh
$ exit
$ docker exec -it docker_dune-rc-march-ru-2_1 bash
$ /dune-rc/docs/multiproc-demo/build-daq.sh
$ exit
```

```bash
$ export PATH="$(pwd)/cli/build:$PATH"
$ run-control get all
$ run-control daq command init all
$ run-control daq command conf all
$ run-control daq command --run-number 42 start all
$ run-control daq command resume daq-app-multiproc-emu
$ run-control daq command stop all
```