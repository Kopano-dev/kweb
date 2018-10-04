/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package config

import (
	"bytes"
	"fmt"
)

// Caddyfile returns a functional caddy file representing our config.
func Caddyfile(config *Config) ([]byte, error) {
	var buf = &bytes.Buffer{}
	var err error

	// Add host.
	err = writeHostToCaddyfile(config, buf)
	if err != nil {
		return nil, err
	}
	// Add bind.
	err = writeBindToCaddyfile(config, buf)
	if err != nil {
		return nil, err
	}
	// Add TLS.
	err = writeTLSToCaddyfile(config, buf)
	if err != nil {
		return nil, err
	}
	// Add base.
	err = writeBaseToCaddyfile(config, buf)
	if err != nil {
		return nil, err
	}
	// Add legacy.
	err = writeLegacyToCaddyfile(config, buf)
	if err != nil {
		return nil, err
	}
	// Add defaults.
	err = writeDefaultsToCaddyfile(config, buf)
	if err != nil {
		return nil, err
	}

	// Debug`
	fmt.Printf("--- cfg start ---\n%s\n-- cfg end --\n", buf.String())

	// Return created config.
	return buf.Bytes(), nil
}
