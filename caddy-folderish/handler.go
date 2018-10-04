/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package caddyfolderish

import (
	"net/http"

	"github.com/mholt/caddy/caddyhttp/httpserver"
)

// FolderishHandler is a handler to return config JSON files.
type FolderishHandler struct {
	path string

	Next httpserver.Handler
}

func (h *FolderishHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	if r.URL.Path == h.path {
		newPath := h.path + "/"
		if q := r.URL.RawQuery; q != "" {
			newPath += "?" + q
		}
		w.Header().Set("Location", newPath)
		w.WriteHeader(http.StatusMovedPermanently)

		return 0, nil
	}

	return h.Next.ServeHTTP(w, r)
}
