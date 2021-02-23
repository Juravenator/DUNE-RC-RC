job "daq-application-manager" {
  datacenters = ["dune-rc"]

  type = "service"

  update {
    max_parallel = 1
  }

  group "daq-application-manager" {

    network {
      port "api" {}
    }

    task "daq-application-manager" {
      driver = "raw_exec"
      config {
        command = "/bin/bash"
        args    = ["/opt/dune/ccm-rc/controllers/daq-app-manager/start.sh"]
      }
      constraint {
        attribute    = "${meta.ccm-rc-controllers-installed}"
        set_contains = "daq-application-manager"
      }
    }

    service {
      name = "daq-application-manager-api"
      port = "api"
    }
  }
}