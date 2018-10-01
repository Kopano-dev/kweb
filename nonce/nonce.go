/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package nonce

import (
	"crypto/rand"
	"encoding/base64"
	"sync"
)

const size = 16
const total = size * size

var pool = &sync.Pool{
	New: func() interface{} {
		return &generator{
			buf: make([]byte, total),
		}
	},
}

// New returns a new nonce base64 standard encoded value suitable to
// be used as CSP nonce value (see https://www.w3.org/TR/CSP2/#base64_value).
func New() []byte {
	g := pool.Get().(*generator)
	nonce := g.nonce()
	pool.Put(g)
	return nonce
}

// A generator reads bulk random data to effenciently return nonces.
type generator struct {
	count uint64
	buf   []byte
}

func (np *generator) nonce() []byte {
	pos := np.count % size
	np.count++
	if pos == 0 {
		_, err := rand.Read(np.buf)
		if err != nil {
			// NOTE(longsleep): Always panic on rand fails.
			panic(err)
		}
	}
	s := pos * size
	e := (pos + 1) * size

	// Encode Base64 standard - see https://www.w3.org/TR/CSP2/#base64_value
	nonce := make([]byte, base64.StdEncoding.EncodedLen(size))
	base64.StdEncoding.Encode(nonce, np.buf[s:e])
	return nonce
}
