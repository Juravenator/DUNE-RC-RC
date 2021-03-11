// get a json version of this by running
// curl -X POST --data @<(jq --arg s "$(<daq-app-multiproc.hcl)" '.JobHCL=$s' -n) http://localhost:4646/v1/jobs/parse | jq '{"meta":{"kind": "nomad-job", "name": .Name},"spec": .} | walk( if type == "object" then with_entries(select(.value != null)) else . end)'
job "daq-app-multiproc" {
  datacenters = ["dune-rc"]

  type = "service"

  update {
    max_parallel = 1
  }

  group "df" {
    network {
      port "triggerDecision" {}
      port "commandFacility" {}
    }
    service {
      name = "daq-app-multiproc-td"
      port = "triggerDecision"
    }
    service {
      name = "daq-app-multiproc-df-cf"
      port = "CommandFacility"
    }
    restart {
      attempts = 1
    }
    task "daq-application" {
      driver = "raw_exec"
      config {
        command = "bash"
        args    = ["/tmp/dune-rc-hacks/multi-proc-daq-2.4/start.sh", "--commandFacility", "rest://localhost:${NOMAD_PORT_CommandFacility}", "--name", "${NOMAD_JOB_NAME}_${NOMAD_GROUP_NAME}"]
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

  }

  group "trgemu" {

    network {
      port "TriggerDecisionToken" {}
      port "TimeSync" {}
      port "CommandFacility" {}
    }

    service {
      name = "daq-app-multiproc-trgemu-cf"
      port = "CommandFacility"
    }
    service {
      name = "daq-app-multiproc-tdt"
      port = "TriggerDecisionToken"
    }
    service {
      name = "daq-app-multiproc-ts"
      port = "TimeSync"
    }

    restart {
      attempts = 1
    }
    task "daq-application" {
      driver = "raw_exec"
      config {
        command = "bash"
        args    = ["/tmp/dune-rc-hacks/multi-proc-daq-2.4/start.sh", "--commandFacility", "rest://localhost:${NOMAD_PORT_CommandFacility}", "--name", "${NOMAD_JOB_NAME}_${NOMAD_GROUP_NAME}"]
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
  }
}