/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package caddystaticpwa

import (
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mholt/caddy/caddyhttp/httpserver"
)

var nonceMarker = []byte("__CSP_NONCE__")

const (
	indexPath         = "/index.html"
	webworkerPath     = "/service-worker.js"
	assetManifestPath = "/asset-manifest.json"
	clientRoute       = "/r"

	clientRoutePrefix      = clientRoute + "/"
	precacheManifestPrefix = "/precache-manifest."
	staticPrefix           = "/static/"
)

// StaticPWAHandler is a handler for static progressive webapps.
type StaticPWAHandler struct {
	host    string
	appURL  string
	handler http.Handler
	fs      http.Dir

	Next httpserver.Handler
}

// NewStaticPWAHandler creates a new StaticPWAHandler with the provided options.
func NewStaticPWAHandler(host, root, name, appURL, path string, next httpserver.Handler) *StaticPWAHandler {
	if !filepath.IsAbs(path) {
		// Use a pre configured pattern or are relative to the configured web
		// root.
		pattern := os.Getenv("STATICPWA_PATH_PATTERN")
		if pattern != "" {
			patternPath := strings.Replace(pattern, "%name", name, -1)
			patternPath = strings.Replace(patternPath, "%path", path, -1)
			path = filepath.Clean(patternPath)
		} else {
			path = filepath.Join(root, path)
		}
	}

	h := &StaticPWAHandler{
		host:   host,
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
		// FIXME(longsleep): Inefficient, avoid additional parse. Rather
		// do StropPrefix some place else or so.
		ourl, _ := url.ParseRequestURI(r.RequestURI)
		// Ensure we are called as path.
		if !strings.HasSuffix(ourl.Path, "/") {
			localRedirect(w, r, path.Base(ourl.Path)+"/")
			return
		}
		// If called as path, always serve index.html directly.
		upath = indexPath
		r.URL.Path = upath
	}

	// Handle headers.
	headers := w.Header()

	// Never send a referrer for pwas.
	headers.Set("Referrer-Policy", "no-referrer")

	// What to open.
	name := path.Clean(upath)

	// Fastpath routes.
	switch {
	case name == clientRoute:
		fallthrough
	case strings.HasPrefix(name, clientRoutePrefix):
		// We know that this does not exist, so open up index directly.
		name = indexPath
	}

	// Open file.
	f, err := h.fs.Open(name)

	// Header routes..
	switch {
	case name == indexPath:
		// pass

	case name == webworkerPath:
		fallthrough
	case name == assetManifestPath:
		// No caching.
		headers.Set("Cache-Control", "public, max-age=0")
		headers.Set("Content-Type", "application/javascript")

	case strings.HasPrefix(name, precacheManifestPrefix):
		fallthrough
	case strings.HasPrefix(name, staticPrefix):
		// Long term caching for static resources.
		headers.Set("Cache-Control", "public, max-age=31536000")
		if strings.HasSuffix(name, ".svg") {
			headers.Set("Content-Security-Policy", staticSVGCSP)
		} else {
			headers.Set("Content-Security-Policy", staticDefaultCSP)
		}

	default:
		if err == nil {
			// Other top level file, no special headers.
		} else if os.IsNotExist(err) {
			// Handle rest with index (it is propably client side URL routing).
			f, err = h.fs.Open(indexPath)
		}
	}

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
		h.handleIndex(w, r, f)
	default:
		http.ServeContent(w, r, d.Name(), d.ModTime(), f)
	}
}
