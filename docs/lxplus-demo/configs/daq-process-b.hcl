// identical to daq-process-b, except '-a'->'-b' and the TCP APA is nr 43
job "daq-process-b" {
  datacenters = ["dune-rc"]

  type = "service"

  update {
    max_parallel = 1
  }

  group "daq-ru" {

    restart {
      attempts = 1
    }

    network {
      port "api" {}
    }

    task "daq-application" {
      driver = "raw_exec"
      config {
        command = "bash"
        args    = ["/tmp/dune-rc-hacks/lxplus-demo/listrev-app-starter.sh", "--commandFacility", "rest://localhost:${NOMAD_PORT_api}", "--name", "${NOMAD_JOB_NAME}"]
      }
      constraint {
        attribute    = "${meta.module-a-tpc-apa}"
        set_contains = "43"
      }
      constraint {
        attribute = "${meta.cvmfs}"
        value     = "/cvmfs"
      }
    }

    service {
      name = "daq-process-b"
      port = "api"
    }
  }
}