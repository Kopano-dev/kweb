/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package config

import (
	"bytes"
)

func writeLogToCaddyfile(config *Config, buf *bytes.Buffer) error {
	// Add log.
	if config.RequestLog != "" {
		buf.WriteString("log / ")
		buf.WriteString(config.RequestLog)
		buf.WriteString(" {combined}\n\n")
	}

	return nil
}
