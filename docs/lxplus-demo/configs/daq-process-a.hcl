// get a json version of this by running
// curl -X POST --data @<(jq --arg s "$(<daq-process-a.hcl)" '.JobHCL=$s' -n) http://localhost:4646/v1/jobs/parse | jq '{"meta":{"kind": "nomad-job", "name": .Name},"spec": .} | walk( if type == "object" then with_entries(select(.value != null)) else . end)'
job "daq-process-a" {
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
        set_contains = "42"
      }
      constraint {
        attribute = "${meta.cvmfs}"
        value     = "/cvmfs"
      }
    }

    service {
      name = "daq-process-a"
      port = "api"
    }
  }
}