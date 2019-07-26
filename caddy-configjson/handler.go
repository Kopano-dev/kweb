/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package caddyconfigjson

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/caddyserver/caddy/caddyhttp/httpserver"
)

const grapiPrefix = "/api/gc/v1"

var defaultJSON = []byte(fmt.Sprintf(`{
  "apiPrefix": "%s"
}`, grapiPrefix))

// ConfigJSONHandler is a handler to return config JSON files.
type ConfigJSONHandler struct {
	url      string
	handler  http.Handler
	fs       *http.Dir
	internal string

	Next httpserver.Handler
}

// NewConfigJSONHandler creates a new ConfigJSONHandler with the provided options.
func NewConfigJSONHandler(url, root, path string, next httpserver.Handler) *ConfigJSONHandler {
	h := &ConfigJSONHandler{
		url: url,

		Next: next,
	}
	if path != "" {
		if !filepath.IsAbs(path) {
			h.internal = "/" + filepath.Clean(path)
			path = filepath.Join(root, path)
		}
		fs := http.Dir(path)
		h.fs = &fs
	}
	h.handler = http.HandlerFunc(h.handle)

	return h
}

func (h *ConfigJSONHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	if strings.HasPrefix(r.RequestURI, h.url) {

		http.StripPrefix(h.url, h.handler).ServeHTTP(w, r)
		return 0, nil
	}

	if h.internal != "" && strings.HasPrefix(r.RequestURI, h.internal) {
		return http.StatusNotFound, nil
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

	// Avoid caching.
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	if h.fs != nil {
		// Support serving existing files from file system. Maps config
		// requests like /:name/config.json to :name.json files found in
		// accociated directory.
		f, err := h.fs.Open(parts[0] + ".json")
		if err == nil {
			d, statErr := f.Stat()
			if statErr == nil {
				if !d.IsDir() {
					http.ServeContent(w, r, "config.json", d.ModTime(), f)
					return
				}
				err = errors.New("config is a directory")
			} else {
				err = statErr
			}
		}
		if !os.IsNotExist(err) {
			// Handle error when file exists but access or read failed.
			if os.IsPermission(err) {
				http.Error(w, "403 Forbidden", http.StatusForbidden)
				return
			}
			// Default.
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	// Return default JSON document.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(defaultJSON)
}
