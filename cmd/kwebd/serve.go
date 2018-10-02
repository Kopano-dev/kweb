/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package main

import (
	"fmt"
	"os"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddy/caddymain"
	"github.com/spf13/cobra"

	"stash.kopano.io/kgol/kweb"
)

func commandServe() *cobra.Command {
	serveCmd := &cobra.Command{
		Use:   "serve [...args]",
		Short: "Start and listen for requests",
		Run: func(cmd *cobra.Command, args []string) {
			if err := serve(cmd, args); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		},
	}

	serveCmd.Flags().Bool("agree", false, "Agree to the CA's Subscriber Agreement")
	serveCmd.Flags().String("ca", "https://acme-v02.api.letsencrypt.org/directory", "URL to certificate authority's ACME server directory")
	serveCmd.Flags().String("email", "", "ACME CA account email address")
	serveCmd.Flags().Bool("http2", true, "Use HTTP/2")
	serveCmd.Flags().Bool("quic", false, "Use experimental QUIC")
	serveCmd.Flags().String("root", ".", "Path to web root")
	serveCmd.Flags().Bool("validate", false, "Parse and validate the configuration but do not start the server")
	serveCmd.Flags().String("revoke", "", "Hostname for which to revoke the certificate")
	serveCmd.Flags().String("host", "0.0.0.0", "Hostname to serve")
	serveCmd.Flags().String("http-port", "80", "Port to use for HTTP")
	serveCmd.Flags().String("https-port", "443", "Port to use for HTTPS")
	serveCmd.Flags().Bool("tls", true, "Enable TLS on listener")
	serveCmd.Flags().Bool("tls-always-self-sign", false, "Always generate self signed certificate")
	serveCmd.Flags().Bool("tls-must-staple", false, "Enable TLS must staple")
	serveCmd.Flags().String("tls-protocols", "tls1.2 tls1.2", "Min and max TLS protocol")
	serveCmd.Flags().String("tls-cert-file", "", "Path to TLS certificate bundle (concatenation of the server's certificate followed by the CA's certificate chain)")
	serveCmd.Flags().String("tls-key-file", "", "Path to the server's private key file which matches the certificate bundle")

	return serveCmd
}

func serve(cmd *cobra.Command, args []string) error {
	caddyArgs := []string{"-type", "http"}

	if agree, _ := cmd.Flags().GetBool("agree"); agree {
		caddyArgs = append(caddyArgs, "-agree")
	}
	if ca, _ := cmd.Flags().GetString("ca"); ca != "" {
		caddyArgs = append(caddyArgs, "-ca", ca)
	}
	if http2, _ := cmd.Flags().GetBool("http2"); http2 {
		caddyArgs = append(caddyArgs, "-http2")
	}
	if quic, _ := cmd.Flags().GetBool("quic"); quic {
		caddyArgs = append(caddyArgs, "-quic")
	}
	if root, _ := cmd.Flags().GetString("root"); root != "" {
		caddyArgs = append(caddyArgs, "-root", root)
	}
	if validate, _ := cmd.Flags().GetBool("validate"); validate {
		caddyArgs = append(caddyArgs, "-validate")
	}
	if revoke, _ := cmd.Flags().GetString("revoke"); revoke != "" {
		caddyArgs = append(caddyArgs, "-revoke", revoke)
	}

	email, _ := cmd.Flags().GetString("email")
	host, _ := cmd.Flags().GetString("host")
	httpPort, _ := cmd.Flags().GetString("http-port")
	httpsPort, _ := cmd.Flags().GetString("https-port")
	tls, _ := cmd.Flags().GetBool("tls")
	tlsAlwaysSelfSign, _ := cmd.Flags().GetBool("tls-always-self-sign")
	tlsMustStaple, _ := cmd.Flags().GetBool("tls-must-staple")
	tlsProtocols, _ := cmd.Flags().GetString("tls-protocols")
	tlsCertBundle, _ := cmd.Flags().GetString("tls-cert-file")
	tlsPrivateKey, _ := cmd.Flags().GetString("tls-key-file")

	caddyArgs = append(caddyArgs, "-host", host, "-http-port", httpPort, "-https-port", httpsPort)
	if tls {
		caddyArgs = append(caddyArgs, "-port", httpsPort)
	} else {
		caddyArgs = append(caddyArgs, "-port", httpPort)
	}
	if email != "" {
		caddyArgs = append(caddyArgs, "-email", email)
	}

	// Configure underlying caddy.
	config := &kweb.Config{
		Host:  host,
		Email: email,

		TLSEnable:         tls,
		TLSAlwaysSelfSign: tlsAlwaysSelfSign,
		TLSCertBundle:     tlsCertBundle,
		TLSPrivateKey:     tlsPrivateKey,
		TLSMustStaple:     tlsMustStaple,
		TLSProtocols:      tlsProtocols,
	}
	caddy.SetDefaultCaddyfileLoader("default", defaultLoader(config))

	// Reset args, since caddymain has its own parsing.
	subArgs := append(os.Args[:1], caddyArgs...)
	os.Args = subArgs

	caddymain.Run()

	return nil
}

func defaultLoader(config *kweb.Config) caddy.LoaderFunc {
	return caddy.LoaderFunc(func(serverType string) (caddy.Input, error) {
		contents, err := kweb.Caddyfile(config)
		if err != nil {
			return nil, err
		}

		return caddy.CaddyfileInput{
			Contents:       contents,
			ServerTypeName: serverType,
		}, nil
	})
}
