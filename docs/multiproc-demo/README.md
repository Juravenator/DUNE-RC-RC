# Running the march demo on LXPlus

## Step 0 - go to LXPlus
```bash
$ ssh lxplus #duh
$ git clone https://github.com/Juravenator/DUNE-RC-RC.git
$ cd DUNE-RC-RC
$ git checkout march-rc0
```

# Quick version

```bash
$ hacks/lxplus/runInMem.sh
$ export PATH="$(pwd)/cli/build:$PATH"
$ docs/multiproc-demo/build-daq.sh
$ run-control apply docs/multiproc-demo/configs/*.json
$ run-control get all
$ run-control daq command init all
$ run-control daq command conf all
$ run-control daq command --run-number 42 start all
$ run-control daq command resume lala
$ run-control daq command stop all
$ hacks/lxplus/stopInMem.sh
```
