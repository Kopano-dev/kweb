/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package main

import (
	"github.com/mholt/caddy/caddy/caddymain"

	// Plugins.
	_ "github.com/captncraig/caddy-realip"
	_ "github.com/mastercactapus/caddy-proxyprotocol"
	_ "github.com/miekg/caddy-prometheus"
	_ "github.com/pyed/ipfilter"
	_ "github.com/xuqingfeng/caddy-rate-limit"
)

func init() {
	caddymain.EnableTelemetry = false
}

func main() {
	caddymain.Run()
}
