/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package caddystaticpwa

import (
	"net"
	"strings"

	"github.com/caddyserver/caddy"
	"github.com/caddyserver/caddy/caddyhttp/httpserver"
)

func init() {
	caddy.RegisterPlugin("staticpwa", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})

	httpserver.RegisterDevDirective("staticpwa", "")
}

func setup(c *caddy.Controller) error {
	for c.Next() {
		// First param is the url prefix.
		if !c.NextArg() {
			return c.ArgErr()
		}
		url := c.Val()

		// Second param is mandatory and the path to the static pwa folder.
		if !c.NextArg() {
			return c.ArgErr()
		}
		path := c.Val()

		// Third parm is optional and defines the name of the pwa.
		name := ""
		if c.NextArg() {
			name = c.Val()
		}
		if name == "" {
			name = strings.Replace(strings.TrimPrefix(url, "/"), "/", "-", -1)
		}

		cfg := httpserver.GetConfig(c)

		host := cfg.Host()
		ip := net.ParseIP(host)
		if ip != nil {
			// Not an IP address - assume it is a hostname.
			host = ""
		} else {
			port := cfg.Port()
			if port != "443" && port != "80" {
				// Add port when not standard.
				host += ":" + port
			}
		}

		// Inject our middle ware.
		mid := func(next httpserver.Handler) httpserver.Handler {
			return NewStaticPWAHandler(
				host,
				cfg.Root,
				name,
				url,
				path,
				next,
			)
		}
		cfg.AddMiddleware(mid)
	}

	return nil
}
