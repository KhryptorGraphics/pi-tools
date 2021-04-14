job "otel-collector" {
  datacenters = ["dc1"]

  type     = "system"
  priority = 70

  group "otel-collector" {
    network {
      port "healthcheck" {
        to = 13133
      }
      port "jaeger_thrift" {
        static = 14268
        to     = 14268
      }
      port "zipkin" {
        static = 9411
        to     = 9411
      }
      port "otlp_grpc" {
        static = 4317
        to     = 4317
      }
      port "otlp_grpc_2" {
        static = 55680
        to     = 55680
      }
      port "otlp_http" {
        static = 55681
        to     = 55681
      }
      port "metrics" {
        to = 8888
      }
    }

    service {
      name = "otel-collector"
      port = "otlp_grpc"

      tags = ["grpc"]

      meta {
        metrics_path = "/metrics"
        metrics_port = "${NOMAD_HOST_PORT_metrics}"
      }

      check {
        type                   = "http"
        path                   = "/"
        port                   = "healthcheck"
        timeout                = "3s"
        interval               = "15s"
        success_before_passing = 3
      }
    }

    task "otel-collector" {
      driver = "docker"

      config {
        image   = "mmoriarity/opentelemetry-collector"
        command = "/otelcol"
        args    = [
          "--config",
          "${NOMAD_SECRETS_DIR}/config.yaml",
          "--mem-ballast-size-mib",
          "50",
        ]
        ports   = [
          "jaeger_thrift",
          "zipkin",
          "otlp_grpc",
          "otlp_grpc_2",
          "otlp_http",
          "healthcheck",
          "metrics",
        ]
      }

      resources {
        cpu    = 100
        memory = 150
      }

      vault {
        policies    = ["otel-collector"]
        change_mode = "noop"
      }

      template {
        data        = file("otel-collector/otel-collector-config.yaml")
        destination = "secrets/config.yaml"
        change_mode = "restart"
      }
    }
  }
}
