/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package config

import (
	"bytes"
)

func writeHostStartToCaddyfile(config *Config, buf *bytes.Buffer) error {
	// Add host.
	buf.WriteString(config.Host)
	buf.WriteString(" {\n\n")

	return nil
}

func writeHostEndToCaddyfile(config *Config, buf *bytes.Buffer) error {
	buf.WriteString("\n}\n\n")

	return nil
}
