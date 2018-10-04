/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package caddystaticpwa

import (
	"net/http"
	"path"

	"strings"

	"github.com/mholt/caddy/caddyhttp/httpserver"
)

var nonceMarker = []byte("__CSP_NONCE__")

const (
	indexPath         = "/index.html"
	webworkerPath     = "/service-worker.js"
	assetManifestPath = "/asset-manifest.json"
)

// StaticPWAHandler is a handler for static progressive webapps.
type StaticPWAHandler struct {
	appURL  string
	handler http.Handler
	fs      http.Dir

	Next httpserver.Handler
}

// NewStaticPWAHandler creates a new StaticPWAHandler with the provided options.
func NewStaticPWAHandler(appURL, path string, next httpserver.Handler) *StaticPWAHandler {
	h := &StaticPWAHandler{
		appURL: appURL,
		fs:     http.Dir(path),

		Next: next,
	}
	h.handler = http.HandlerFunc(h.handle)

	return h
}

func (h *StaticPWAHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	if strings.HasPrefix(r.RequestURI, h.appURL) {
		http.StripPrefix(h.appURL, h.handler).ServeHTTP(w, r)
		return 0, nil
	}

	return h.Next.ServeHTTP(w, r)
}

func (h *StaticPWAHandler) handle(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}
	if upath == "/" {
		url := r.RequestURI
		// Ensure we are called as path.
		if !strings.HasSuffix(url, "/") {
			localRedirect(w, r, path.Base(url)+"/")
			return
		}
		// If called as path, always serve index.html directly.
		upath = indexPath
		r.URL.Path = upath
	}

	// Handle headers.
	headers := w.Header()
	headers.Set("X-Content-Type-Options", "nosniff")
	headers.Set("X-XSS-Protection", "1; mode=block")
	headers.Set("X-Frame-Options", "sameorigin")
	headers.Set("Referrer-Policy", "no-referrer")
	headers.Set("Feature-Policy", "midi 'none'")

	name := path.Clean(upath)

	// Routes.
	switch name {
	case indexPath:
		// pass

	case webworkerPath:
		fallthrough
	case assetManifestPath:
		// No caching.
		headers.Set("Cache-Control", "public, max-age=0")
		headers.Set("Content-Type", "application/javascript")

	default:
		if strings.HasPrefix(name, "/static/") {
			// Long term caching for static resources.
			headers.Set("Cache-Control", "public, max-age=31536000")
			headers.Set("Content-Security-Policy", "default-src 'self'")

		} else {
			// Handle rest with index (it is propably client side URL routing).
			name = "/index.html"
		}
	}

	// Open file.
	f, err := h.fs.Open(name)
	if err != nil {
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		return
	}
	defer f.Close()

	d, err := f.Stat()
	if err != nil {
		msg, code := toHTTPError(err)
		http.Error(w, msg, code)
		return
	}

	if d.IsDir() {
		// No directories.
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}

	// Handle content.
	switch name {
	case indexPath:
		handleIndex(w, r, f)
	default:
		http.ServeContent(w, r, d.Name(), d.ModTime(), f)
	}
}
