/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package config

import (
	"io"
)

func writeBaseToCaddyfile(config *Config, w io.Writer) error {
	// Add base.
	_, err := w.Write(config.Base)

	return err
}
