module stash.kopano.io/kgol/kweb

go 1.13

require (
	github.com/caddyserver/caddy v1.0.5
	github.com/captncraig/caddy-realip v0.0.0-20190710144553-6df827e22ab8
	github.com/cenkalti/backoff/v4 v4.0.2 // indirect
	github.com/go-acme/lego/v3 v3.7.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/klauspost/cpuid v1.2.5 // indirect
	github.com/mholt/certmagic v0.8.3
	github.com/miekg/dns v1.1.29 // indirect
	github.com/oschwald/maxminddb-golang v1.3.1 // indirect
	github.com/prometheus/procfs v0.1.3 // indirect
	github.com/pyed/ipfilter v1.1.4
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/xuqingfeng/caddy-rate-limit v1.6.7
	stash.kopano.io/kgol/caddy-prometheus v0.0.0-20190726090614-6055bc7a4bdf
)

replace github.com/caddyserver/caddy => github.com/longsleep/caddy v1.0.6-0.20211008153515-2e836ef87ed2

replace github.com/lucas-clemente/quic-go => github.com/lucas-clemente/quic-go v0.21.1
