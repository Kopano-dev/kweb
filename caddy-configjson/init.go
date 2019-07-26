/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package caddyconfigjson

import (
	"github.com/caddyserver/caddy"
	"github.com/caddyserver/caddy/caddyhttp/httpserver"
)

func init() {
	caddy.RegisterPlugin("configjson", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})

	httpserver.RegisterDevDirective("configjson", "")
}

func setup(c *caddy.Controller) error {
	for c.Next() {
		// First param is the url prefix.
		if !c.NextArg() {
			return c.ArgErr()
		}
		url := c.Val()

		// Second parm is optional and a path where to find config json files.
		var path string
		if c.NextArg() {
			path = c.Val()
		}

		// Inject our middle ware.
		cfg := httpserver.GetConfig(c)
		mid := func(next httpserver.Handler) httpserver.Handler {
			return NewConfigJSONHandler(
				url,
				cfg.Root,
				path,
				next,
			)
		}
		cfg.AddMiddleware(mid)
	}

	return nil
}
