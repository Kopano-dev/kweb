/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package caddystaticpwa

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"stash.kopano.io/kgol/kweb/nonce"
)

func (h *StaticPWAHandler) handleIndex(w http.ResponseWriter, r *http.Request, f io.ReadSeeker) {
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

	// Compute host.
	var secureHost string
	if h.host == "" {
		// No host set - this is potentically insecure if no other validation
		// of the incoming host header takes places.
		secureHost += r.Host
	} else {
		secureHost += h.host
	}

	// CSP and no caching.
	headers := w.Header()
	if h.IndexCSPTemplate != "" {
		r := strings.NewReplacer("%nonce", string(n), "%securehost", secureHost)
		headers.Set("Content-Security-Policy", r.Replace(h.IndexCSPTemplate))
	}
	headers.Set("Cache-Control", "private, max-age=0")

	// Directly return data from replaced content.
	headers.Set("Content-Type", "text/html; charset=utf-8")
	headers.Set("Accept-Ranges", "none")
	headers.Set("Content-Length", strconv.FormatInt(sendSize, 10))
	w.WriteHeader(http.StatusOK)
	if r.Method != http.MethodHead {
		w.Write(content)
	}
}
