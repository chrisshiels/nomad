job "time" {
  region = "global"
  datacenters = [ "dc1" ]
  type = "service"

  update {
    stagger = "10s"
    max_parallel = 1
  }

  group "cache" {
    count = 3
    restart {
      attempts = 10
      interval = "5m"
      delay = "25s"
      mode = "delay"
    }

    task "time" {
      driver = "docker"
      config {
        image = "vm1:5000/cs/time:1.0.0"
        port_map {
          http = 7002
        }
      }

      resources {
        cpu    = 500
        memory = 256
        network {
          port "http" {}
        }
      }

      service {
        name = "time"
        tags = [ "urlprefix-time/" ]
        port = "http"
        check {
          type     = "http"
          protocol = "http"
          port     = "http"
          path     = "/status"
          interval = "10s"
          timeout  = "2s"
        }
      }

      kill_timeout = "20s"
    }
  }
}
