job "daq-ru-pu-52-51" {
  datacenters = ["dune-rc"]

  type = "service"

  update {
    max_parallel = 1
  }

  group "daq-ru" {

    network {
      port "api" {}
    }

    task "daq-application" {
      driver = "raw_exec"
      config {
        command = "bash"
        args    = ["/dune-rc/hacks/daq-app-starter.sh", "--commandFacility", "rest://localhost:${NOMAD_PORT_api}"]
      }
      constraint {
        attribute    = "${meta.module-a-pu-apas}"
        set_contains = "42-51"
      }
      constraint {
        attribute = "${meta.cvmfs}"
        value     = "/cvmfs"
      }
    }

    service {
      name = "daq-ru-pu-52-51-api"
      port = "api"
      // check {
      //   type = "http"
      //   path = "/health"
      //   interval = "10s"
      //   timeout = "2s"
      // }
    }
  }
}