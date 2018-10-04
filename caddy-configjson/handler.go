/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package caddyconfigjson

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/mholt/caddy/caddyhttp/httpserver"
)

// VersionStamp is the version added to config JSON kweb.v key.
const VersionStamp = 20181004

const grapiPrefix = "/api/gc/v1"

var defaultJSON = []byte(fmt.Sprintf(`{
  "apiPrefix": "%s",
  "kweb": {
    "v": %d
  }
}`, grapiPrefix, VersionStamp))

// ConfigJSONHandler is a handler to return config JSON files.
type ConfigJSONHandler struct {
	url     string
	handler http.Handler
	fs      http.Dir

	Next httpserver.Handler
}

// NewConfigJSONHandler creates a new ConfigJSONHandler with the provided options.
func NewConfigJSONHandler(url, path string, next httpserver.Handler) *ConfigJSONHandler {
	h := &ConfigJSONHandler{
		url: url,

		Next: next,
	}
	if path != "" {
		h.fs = http.Dir(path)
	}
	h.handler = http.HandlerFunc(h.handle)

	return h
}

func (h *ConfigJSONHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	if strings.HasPrefix(r.RequestURI, h.url) {

		http.StripPrefix(h.url, h.handler).ServeHTTP(w, r)
		return 0, nil
	}

	return h.Next.ServeHTTP(w, r)
}

func (h *ConfigJSONHandler) handle(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path

	// Limit scope.
	if !strings.HasSuffix(upath, "/config.json") {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}
	parts := strings.SplitN(upath, "/", 3)
	if len(parts) > 2 {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}

	// Return empty JSON document.
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.WriteHeader(http.StatusOK)
	w.Write(defaultJSON)
}
