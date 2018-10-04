/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package config

import (
	"bytes"
)

func writeBindToCaddyfile(config *Config, buf *bytes.Buffer) error {
	// Add bind.
	if config.Bind != "" {
		buf.WriteString("bind ")
		buf.WriteString(config.Bind)
		buf.WriteString("\n\n")
	}

	return nil
}
