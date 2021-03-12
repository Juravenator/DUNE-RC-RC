To have stuff work on your local machine, and not have to set up a custom host. You could opt to rsync the DUNE space in cvmfs locally.

Warning: This will cost ~20GiB of your disk and network for a first sync.

```bash
sudo mkdir -p /cvmfs/dune.opensciencegrid.org/dunedaq/DUNE
sudo chown -R `whoami`:`whoami` /cvmfs
rsync -a --info=progress2 lxplus:/cvmfs/dune.opensciencegrid.org/dunedaq/DUNE/products /cvmfs/dune.opensciencegrid.org/dunedaq/DUNE
rsync -a --info=progress2 lxplus:/cvmfs/dune.opensciencegrid.org/dunedaq/DUNE/products_dev /cvmfs/dune.opensciencegrid.org/dunedaq/DUNE
rsync -a --info=progress2 lxplus:/cvmfs/dune.opensciencegrid.org/dunedaq/DUNE/releases /cvmfs/dune.opensciencegrid.org/dunedaq/DUNE
rsync -a --info=progress2 lxplus:/cvmfs/dune.opensciencegrid.org/dunedaq/DUNE/pypi-repo /cvmfs/dune.opensciencegrid.org/dunedaq/DUNE
```