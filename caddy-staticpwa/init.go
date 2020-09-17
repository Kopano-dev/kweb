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
		args := c.RemainingArgs()

		if len(args) < 2 || len(args) > 3 {
			return c.ArgErr()
		}
		// First param is the url prefix.
		url := args[0]
		// Second param is mandatory and the path to the static pwa folder.
		path := args[1]
		// Third parm is optional and defines the name of the pwa.
		name := ""
		if len(args) > 2 {
			name = args[2]
		}
		if name == "" {
			name = strings.Replace(strings.TrimPrefix(url, "/"), "/", "-", -1)
		}

		var indexCSPTemplate string
		var staticDefaultCSP string
		var staticSVGCSP string

		for c.NextBlock() {
			switch c.Val() {
			case "csp_index":
				if !c.NextArg() {
					return c.ArgErr()
				}
				indexCSPTemplate = c.Val()
			case "csp_default":
				if !c.NextArg() {
					return c.ArgErr()
				}
				staticDefaultCSP = c.Val()
			case "csp_svg":
				if !c.NextArg() {
					return c.ArgErr()
				}
				staticSVGCSP = c.Val()
			}
		}

		cfg := httpserver.GetConfig(c)

		host := cfg.Host()
		ip := net.ParseIP(host)
		if ip != nil || host == "*" {
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
			h := NewStaticPWAHandler(
				host,
				cfg.Root,
				name,
				url,
				path,
				next,
			)
			h.IndexCSPTemplate = indexCSPTemplate
			h.StaticDefaultCSP = staticDefaultCSP
			h.StaticSVGCSP = staticSVGCSP

			return h
		}
		cfg.AddMiddleware(mid)
	}

	return nil
}
