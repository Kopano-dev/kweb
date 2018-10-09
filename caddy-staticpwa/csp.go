/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package caddystaticpwa

const indexCSPTemplate = "default-src 'self'; " +
	"style-src 'self' 'nonce-%s'; " + // Random nonce.
	"base-uri 'none'; " +
	"object-src 'none'; " + // Disabled for security - no crap in our house.
	"connect-src 'self' %s; " + // Additional connect urls.
	"img-src 'self' data:; " + // NOTE(longsleep): We need data image URLs for now.
	"frame-ancestors 'self'" // NOTE(longsleep): Better than X-Frame-Options since this goes up to the top frame.

const staticDefaultCSP = "default-src 'self'; " +
	"img-src 'self' data:; " +
	"object-src 'none'"

const staticSVGCSP = staticDefaultCSP + "; " +
	"style-src 'self' 'unsafe-inline'" // Allow inline style for svg images.
