/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package config

import (
	"bytes"
	"fmt"
)

const defaultsRedirect = `
redir 301 {
	if {path} is /
	/ %s
}`

func writeDefaultsToCaddyfile(config *Config, buf *bytes.Buffer) error {
	// Add bind.
	if config.DefaultRedirect != "" {
		buf.WriteString(fmt.Sprintf(defaultsRedirect, config.DefaultRedirect))
		buf.WriteString("\n")
	}

	return nil
}
