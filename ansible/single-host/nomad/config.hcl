datacenter = "dune-rc"

plugin "raw_exec" {
  config {
    enabled = true
    no_cgroups = true
  }
}

client {
  enabled = true
  node_class = "readout_unit"
  meta {
    cvmfs = "/cvmfs"
    module-a-tpc-apa = "42,43"
    module-a-pu-apas = "42-51,52-61"
    ccm-rc-controllers-installed = "daq-application-manager"
  }
}