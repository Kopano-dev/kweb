/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package config

import (
	"bytes"
	"fmt"
)

var legacyReverseProxyWebapp = `
proxy /webapp/ %s {
	fail_timeout 10s
	try_duration 30s
	transparent
	keepalive 100
}
folderish /webapp
`

var legacyReverseProxyZPush = `
proxy /Microsoft-Server-ActiveSync %s {
	transparent
	keepalive 0
	timeout 3540s
}
proxy /AutoDiscover/AutoDiscover.xml %s {
	transparent
	keepalive 0
	fail_timeout 10s
	try_duration 30s
}
`

func writeLegacyToCaddyfile(config *Config, buf *bytes.Buffer) error {
	if config.ReverseProxyLegacyHTTP != "" {
		buf.WriteString(fmt.Sprintf(legacyReverseProxyWebapp, config.ReverseProxyLegacyHTTP))
		buf.WriteString(fmt.Sprintf(legacyReverseProxyZPush, config.ReverseProxyLegacyHTTP, config.ReverseProxyLegacyHTTP))
	}

	return nil
}
