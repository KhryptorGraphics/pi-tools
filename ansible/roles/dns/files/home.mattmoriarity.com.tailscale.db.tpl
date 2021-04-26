; BIND data file for homelab zone, when served from a Tailscale IP
$TTL  1m
@   IN  SOA localhost. matt.mattmoriarity.com. (
                  1
                 1m     ; Refresh
                 1h             ; Retry
                 1w             ; Expire
                 1h )   ; Negative Cache TTL
@   IN  NS  localhost.
;{{ range nodes }}
{{ .Node }} IN  A {{ .Meta.tailscale_ip }}{{ end }}
;{{ range service "ingress-http" }}
*.home.mattmoriarity.com.  0  IN  A {{ .NodeMeta.tailscale_ip }}{{ end }}
unifi       IN  A 10.0.0.1
nas             IN  A 10.0.0.10
mars        IN  A 100.117.39.47
localhost   IN  A 127.0.0.1
