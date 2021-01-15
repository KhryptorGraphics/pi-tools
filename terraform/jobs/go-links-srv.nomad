job "go-links" {
  datacenters = [
    "dc1"
  ]

  type = "service"

  group "go-links" {
    count = 2

    network {
      mode = "bridge"
      port "http" {
        to = 4240
      }
    }

    service {
      name = "go-links"
      port = 4240

      check {
        type = "http"
        port = "http"
        path = "/healthz"
        timeout = "3s"
        interval = "15s"
        success_before_passing = 3
      }

      connect {
        sidecar_service {
          proxy {
            upstreams {
              destination_name = "postgresql"
              local_bind_port = 5432
            }
          }
        }
      }
    }

    service {
      name = "go-links-metrics"
      port = "http"

      meta {
        metrics_path = "/metrics"
      }
    }

    service {
      name = "go-links-grpc"
      port = 4241

//      check {
//        type = "grpc"
//        timeout = "3s"
//        interval = "15s"
//        success_before_passing = 3
//      }
      connect {
        sidecar_service {}
      }
    }

    task "go-links-srv" {
      driver = "docker"

      config {
        image = "mmoriarity/go-links-srv@__DIGEST__"
        command = "/go-links"
        args = [
          "-db",
          "dbname=golinks host=localhost sslmode=disable",
        ]

        logging {
          type = "journald"
          config {
            tag = "go-links-srv"
          }
        }
      }

      resources {
        cpu = 50
        memory = 50
      }

      vault {
        policies = ["go-links"]
      }

      template {
        data = <<EOF
{{ with secret "database/creds/go-links" }}
PGUSER="{{ .Data.username }}"
PGPASSWORD={{ .Data.password | toJSON }}
{{ end }}
EOF
        destination = "secrets/db.env"
        env = true
        change_mode = "restart"
      }
    }
  }
}
