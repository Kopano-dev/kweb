/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package config

import (
	"bytes"
)

func writeHostToCaddyfile(config *Config, buf *bytes.Buffer) error {
	// Add host.
	buf.WriteString(config.Host)
	buf.WriteString("\n\n")

	return nil
}
