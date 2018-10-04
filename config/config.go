/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package config

// Config bundles a bunch of configuration settings.
type Config struct {
	Root string

	Bind  string
	Host  string
	Email string

	TLSEnable         bool
	TLSAlwaysSelfSign bool
	TLSCertBundle     string
	TLSPrivateKey     string
	TLSKeyType        string
	TLSProtocols      string
	TLSMustStaple     bool

	ReverseProxyLegacyHTTP string
	DefaultRedirect        string

	Extra []byte
}
