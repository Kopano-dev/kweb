/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package caddyfolderish

import (
	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

func init() {
	caddy.RegisterPlugin("folderish", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})

	httpserver.RegisterDevDirective("folderish", "redir")
}

func setup(c *caddy.Controller) error {
	for c.Next() {
		// First param is the path which is folderish.
		if !c.NextArg() {
			return c.ArgErr()
		}
		path := c.Val()

		// Inject our middle ware.
		cfg := httpserver.GetConfig(c)
		mid := func(next httpserver.Handler) httpserver.Handler {
			return &FolderishHandler{
				path: path,
				Next: next,
			}
		}
		cfg.AddMiddleware(mid)
	}

	return nil
}
