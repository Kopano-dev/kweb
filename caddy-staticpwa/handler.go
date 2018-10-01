/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package caddystaticpwa

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/mholt/caddy/caddyhttp/httpserver"

	"stash.kopano.io/kgol/kweb/nonce"
)

var nonceMarker = []byte("__CSP_NONCE__")

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
		upath = "/index.html"
		r.URL.Path = upath
	}

	name := path.Clean(upath)

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

	// Handle headers.
	headers := w.Header()
	headers.Set("X-Content-Type-Options", "nosniff")
	headers.Set("X-XSS-Protection", "1; mode=block")
	headers.Set("X-Frame-Options", "sameorigin")
	headers.Set("Referrer-Policy", "no-referrer")
	headers.Set("Feature-Policy", "midi 'none'")

	switch name {
	case "/index.html":
		index, err := ioutil.ReadAll(f)
		if err != nil {
			http.Error(w, "500 failed to load app", http.StatusInternalServerError)
			return
		}

		// Nonce.
		n := nonce.New()
		// NOTE(longsleep): This is not particularly efficient.
		content := bytes.Replace(index, nonceMarker, n, 1)
		sendSize := int64(len(content))

		// CSP and no caching.
		headers.Set("Content-Security-Policy", fmt.Sprintf("default-src 'self'; style-src 'self' 'nonce-%s'; base-uri 'none'", string(n)))
		headers.Set("Cache-Control", "private, max-age=0")

		// Directly return data from replaced content.
		headers.Set("Content-Type", "text/html; charset=utf-8")
		headers.Set("Accept-Ranges", "none")
		headers.Set("Content-Length", strconv.FormatInt(sendSize, 10))
		w.WriteHeader(http.StatusOK)
		if r.Method != "HEAD" {
			w.Write(content)
		}
		return

	case "/service-worker.js":
		// No caching.
		headers.Set("Cache-Control", "public, max-age=0")
		headers.Set("Content-Type", "application/javascript")

	default:
		if strings.HasPrefix(name, "/static/") {
			// Long term caching for static resources.
			headers.Set("Cache-Control", "public, max-age=31536000")
			headers.Set("Content-Security-Policy", "default-src 'self'")
		}
	}

	// Serve content.
	http.ServeContent(w, r, d.Name(), d.ModTime(), f)
}
