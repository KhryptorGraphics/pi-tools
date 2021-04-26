; BIND data file for homelab zone, when served from a local LAN IP
$TTL  1m
@   IN  SOA localhost. matt.mattmoriarity.com. (
                  1
                 1m     ; Refresh
                 1h             ; Retry
                 1w             ; Expire
                 1h )   ; Negative Cache TTL
@   IN  NS  localhost.
;{{ range nodes }}
{{ .Node }} IN  A {{ .Address }}{{ end }}
*.home.mattmoriarity.com.  IN  CNAME ingress-http.service.consul.
unifi       IN  A 10.0.0.1
nas         IN  A 10.0.0.10
mars        IN  A 10.0.0.50
localhost   IN  A 127.0.0.1
