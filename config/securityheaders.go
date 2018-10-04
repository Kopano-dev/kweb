/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package config

import (
	"bytes"
)

func writeSecurityHeadersToCaddyfile(config *Config, buf *bytes.Buffer) error {
	buf.WriteString("header / {\n")

	// Bunch of security headers.
	buf.WriteString("X-Content-Type-Options \"nosniff\"\n")
	buf.WriteString("X-XSS-Protection \"1; mode=block\"\n")
	buf.WriteString("X-Frame-Options \"sameorigin\"\n")
	buf.WriteString("Feature-Policy \"midi 'none'\"\n")

	// Add hsts when set and we have a host configured.
	if config.HSTS != "" && config.Host != "" {
		buf.WriteString("Strict-Transport-Security \"")
		buf.WriteString(config.HSTS)
		buf.WriteString("\"\n")
	}

	buf.WriteString("}\n\n")
	return nil
}
