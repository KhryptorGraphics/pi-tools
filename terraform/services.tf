resource "consul_config_entry" "global_proxy_defaults" {
  kind = "proxy-defaults"
  name = "global"

  config_json = jsonencode({
    Config = {
      envoy_prometheus_bind_addr = "0.0.0.0:9102"
    }
  })
}

resource "consul_config_entry" "vault_proxy_defaults" {
  kind = "service-defaults"
  name = "vault-proxy"

  config_json = jsonencode({
    Protocol = "http"
  })
}

resource "consul_config_entry" "vault_proxy_intentions" {
  kind = "service-intentions"
  name = "vault-proxy"

  config_json = jsonencode({
    Sources = [
      {
        Name        = "ingress-http"
        Precedence  = 9
        Type        = "consul"
        Permissions = [
          {
            Action = "allow"
            HTTP   = {
              PathPrefix = "/"
            }
          },
        ]
      },
      {
        Action     = "deny"
        Name       = "*"
        Precedence = 8
        Type       = "consul"
      },
    ]
  })
}
