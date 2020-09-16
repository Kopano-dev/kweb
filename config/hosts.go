/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package config

import (
	"io"
)

func writeHostsToCaddyfile(config *Config, w io.Writer) error {
	// Add hosts.
	_, err := w.Write(config.HostsD)

	return err
}
