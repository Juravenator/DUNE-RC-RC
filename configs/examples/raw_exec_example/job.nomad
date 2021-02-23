# you can make me into a JSON using
# https://www.hcl2json.com/

job "example-job" {
  datacenters = ["dune-rc"]

  type = "service"

  update {
    // our app cannot run multiple versions of itself side by side
    // when deploying a new version, this will cause the old version
    // to be destroyed before the new is started
    max_parallel = 1
  }

  // tasks are grouped, tasks in a group are co-located on the same compute unit
  group "my-app" {

    network {
      // our app will need some (dynamically chosen) port to listen on
      port "api" {}
    }

    task "my-app" {
      driver = "raw_exec"
      config {
        command = "/usr/bin/python3"
        args    = ["-m", "http.server", "${NOMAD_PORT_api}"]
      }
      // resources {
      //   cpu = 500 # MHz
      //   memory = 128 # MB
      // }
      // constraint {
      //   // force to run on specific node
      //   attribute = "${node.unique.hostname}"
      //   // operator  = "="
      //   // nomad server members
      //   value     = "f4c4bc7602e4"
      // }
    }



    // we can register this app under a name, so we don't have to manually
    // figure out host:port in other apps
    service {
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