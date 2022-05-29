// nolint
package templates 

const (
DYNAMIC_TOML = `[tls.options]

[tls.options.default]
minVersion = "VersionTLS12"
sniStrict = true

[[tls.certificates]]
certFile = "/var/certs/cert.pem"
keyFile = "/var/certs/key.pem"
`
TRAEFIK_TOML = `[log]
level = "DEBUG"

[providers]
[providers.docker]
endpoint = "tcp://wpe_dockerproxy:2375"
exposedByDefault = false
network = "traefik-proxy"
[providers.file]
filename = "/etc/traefik/dynamic.toml"

[api]
dashboard = true
debug = true
insecure = true

[entryPoints]
[entryPoints.web]
address = ":80"
[entryPoints.web-secure]
address = ":443"
[entryPoints.mysql]
address = ":3306"
`
)