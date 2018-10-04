/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package config

import (
	"bytes"
)

func writeSecurityHeadersToCaddyfile(config *Config, buf *bytes.Buffer) error {
	// Add hsts when set and we have a host configured.
	if config.HSTS != "" && config.Host != "" {
		buf.WriteString("header / Strict-Transport-Security \"")
		buf.WriteString(config.HSTS)
		buf.WriteString("\"\n\n")
	}

	return nil
}
