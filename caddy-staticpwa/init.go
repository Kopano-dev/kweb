/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package caddystaticpwa

import (
	"path/filepath"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
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

		cfg := httpserver.GetConfig(c)

		if !filepath.IsAbs(path) {
			// Relative paths are relative to the configured web root.
			path = filepath.Join(cfg.Root, path)
		}

		// Inject our middle ware.
		mid := func(next httpserver.Handler) httpserver.Handler {
			return NewStaticPWAHandler(
				url,
				path,
				next,
			)
		}
		cfg.AddMiddleware(mid)
	}

	return nil
}
