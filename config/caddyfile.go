/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package config

import (
	"bytes"
)

// Caddyfile returns a functional caddy file representing our config.
func Caddyfile(config *Config) ([]byte, error) {
	var buf = &bytes.Buffer{}

	// Add base.
	err := writeBaseToCaddyfile(config, buf)
	if err != nil {
		return nil, err
	}

	// Return created config.
	return buf.Bytes(), nil
}
