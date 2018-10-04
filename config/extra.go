/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package config

import (
	"bytes"
)

func writeExtraToCaddyfile(config *Config, buf *bytes.Buffer) error {
	// Add extra.
	if config.Extra != nil {
		buf.WriteString("\n# Custom extra configuration\n")
		buf.Write(config.Extra)
		buf.WriteString("\n")
	}

	return nil
}
