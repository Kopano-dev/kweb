/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package kweb

import (
	"bytes"
	"fmt"
)

var caddyfile = []byte(`
errors stderr
log stdout

gzip

header / Server kweb

limits {
	header 1MB
	body   50MB
}

# Config
configjson /api/config/v1/kopano/

# Konnect
proxy /upstreams/konnect/ {
	without /upstreams/konnect/
	upstream 127.0.0.1:8777
	policy least_conn
	health_check /health-check
	fail_timeout 10s
	try_duration 30s
	keepalive 100
	transparent
	header_downstream Feature-Policy "midi 'none'"
	header_downstream X-Frame-Options "sameorigin"
}
ratelimit * 100 200 minute {
	/upstreams/konnect/v1/
	/signin/v1/identifier/_/
	whitelist 127.0.0.1/8
}
rewrite /.well-known/openid-configuration {
	to /upstreams/konnect/{path}
}
rewrite /konnect/v1/ {
	to /upstreams/konnect/{path}
}
rewrite /signin/v1/ {
	to /upstreams/konnect/{path}
}
redir /signin /signin/v1/identifier

# Kapi
proxy /upstreams/kapi/ {
	without /upstreams/kapi/
	upstream 127.0.0.1:8039
	policy least_conn
	health_check /health-check
	fail_timeout 10s
	try_duration 30s
	keepalive 100
	transparent
	websocket
}
ratelimit * 100 200 minute {
	/upstreams/kapi/api/
	whitelist 127.0.0.1/8
}
rewrite /api/gc/v1/ {
	to /upstreams/kapi/{path}
}
rewrite /api/pubs/v1/ {
	to /upstream/kapi/{path}
}

# Kwmserver
proxy /upstreams/kwmserver/ {
	without /upstreams/kwmserver/
	upstream 127.0.0.1:8778
	policy least_conn
	health_check /health-check
	fail_timeout 10s
	try_duration 30s
	keepalive 100
	transparent
	websocket
}
ratelimit * 100 200 minute {
	/upstreams/kwmserver/
	whitelist 127.0.0.1/8
}
rewrite /api/v1/websocket {
	to /upstreams/kwmserver/{path}
}
rewrite /api/v1/rtm.connect {
	to /upstreams/kwmserver/{path}
}
rewrite /api/v1/rtm.turn {
	to /upstreams/kwmserver/{path}
}

# Known Kopano static progressive webapps
staticpwa /meet ./meet-webapp
staticpwa /calendar ./calendar-webapp
staticpwa /mail ./mail-webapp
staticpwa /contacts ./contacts-webapp
`)

// Config bundles a bunch of configuration settings.
type Config struct {
	Host  string
	Email string

	TLSEnable         bool
	TLSAlwaysSelfSign bool
	TLSCertBundle     string
	TLSPrivateKey     string
	TLSKeyType        string
	TLSProtocols      string
	TLSMustStaple     bool
}

// Caddyfile returns a functional caddy file representing our config.
func Caddyfile(config *Config) ([]byte, error) {
	var buf = &bytes.Buffer{}

	// Add host.
	buf.WriteString(config.Host)
	buf.WriteString("\n\n")

	// TLS setup
	if !config.TLSEnable {
		buf.WriteString("tls off\n\n")
	} else if (config.Email != "" || (config.TLSCertBundle != "" && config.TLSPrivateKey != "")) && !config.TLSAlwaysSelfSign {
		buf.WriteString("tls")
		if config.TLSCertBundle != "" && config.TLSPrivateKey != "" {
			buf.WriteString(" ")
			buf.WriteString(config.TLSCertBundle)
			buf.WriteString(" ")
			buf.WriteString(config.TLSPrivateKey)
		}
		buf.WriteString(" {\n")
		if config.TLSProtocols != "" {
			buf.WriteString("protocols ")
			buf.WriteString(config.TLSProtocols)
			buf.WriteString("\n")
		}
		if config.TLSMustStaple {
			buf.WriteString("must-staple\n")
		}
		if config.TLSKeyType != "" {
			buf.WriteString("key_type ")
			buf.WriteString(config.TLSKeyType)
			buf.WriteString("\n")
		}
		buf.WriteString("}\n\n")
	} else {
		buf.WriteString("tls self_signed {\n")
		if config.TLSKeyType != "" {
			buf.WriteString("key_type ")
			buf.WriteString(config.TLSKeyType)
			buf.WriteString("\n")
		}
		buf.WriteString("}\n\n")
	}

	// Add routes.
	buf.Write(caddyfile)

	// Debug`
	fmt.Printf("--- cfg start ---\n%s\n-- cfg end --\n", buf.String())

	// Return created config.
	return buf.Bytes(), nil
}
