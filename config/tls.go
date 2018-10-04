/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package config

import (
	"bytes"
)

func writeTLSToCaddyfile(config *Config, buf *bytes.Buffer) error {
	// TLS setup
	if !config.TLSEnable {
		buf.WriteString("tls off\n\n")
	} else if (config.Email != "" || (config.TLSCertBundle != "" && config.TLSPrivateKey != "")) && !config.TLSAlwaysSelfSign {
		buf.WriteString("tls")
		if config.TLSCertBundle != "" && config.TLSPrivateKey != "" {
			buf.WriteString(" ")
			buf.WriteString(config.TLSCertBundle)
			buf.WriteString(" ")
			buf.WriteString(config.TLSPrivateKey)
		}
		buf.WriteString(" {\n")
		if config.TLSProtocols != "" {
			buf.WriteString("protocols ")
			buf.WriteString(config.TLSProtocols)
			buf.WriteString("\n")
		}
		if config.TLSMustStaple {
			buf.WriteString("must-staple\n")
		}
		if config.TLSKeyType != "" {
			buf.WriteString("key_type ")
			buf.WriteString(config.TLSKeyType)
			buf.WriteString("\n")
		}
		buf.WriteString("}\n\n")
	} else {
		buf.WriteString("tls self_signed {\n")
		if config.TLSKeyType != "" {
			buf.WriteString("key_type ")
			buf.WriteString(config.TLSKeyType)
			buf.WriteString("\n")
		}
		buf.WriteString("}\n\n")
	}

	return nil
}
