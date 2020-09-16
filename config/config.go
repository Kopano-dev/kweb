/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package config

import (
	"strings"
)

// Config bundles a bunch of configuration settings.
type Config struct {
	Root string

	Bind  string
	Host  string
	Email string

	RequestLog string

	TLSEnable         bool
	TLSAlwaysSelfSign bool
	TLSCertBundle     string
	TLSPrivateKey     string
	TLSKeyType        string
	TLSProtocols      string
	TLSMustStaple     bool

	HSTS                   string
	ReverseProxyLegacyHTTP string
	DefaultRedirect        string

	HTTPPortString  string
	HTTPSPortString string

	HostsD    []byte
	SnippetsD []byte
	DefaultsD []byte
	ExtraD    []byte
}

// Hosts returns all hosts of the associcated Config.Host field as string slice.
func (c *Config) Hosts() []string {
	return strings.Split(c.Host, " ")
}
