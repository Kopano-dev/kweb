/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/caddyserver/caddy"
	"github.com/caddyserver/caddy/caddy/caddymain"
	"github.com/mholt/certmagic"
	"github.com/spf13/cobra"

	"stash.kopano.io/kgol/kweb/config"
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
	serveCmd.Flags().String("ca", certmagic.Default.CA, "URL to certificate authority's ACME server directory")
	serveCmd.Flags().String("catimeout", certmagic.HTTPTimeout.String(), "Default ACME CA HTTP timeout")
	serveCmd.Flags().String("email", "", "ACME CA account email address")
	serveCmd.Flags().Bool("http2", true, "Use HTTP/2")
	serveCmd.Flags().Bool("quic", false, "Use experimental QUIC")
	serveCmd.Flags().String("root", ".", "Path to web root")
	serveCmd.Flags().Bool("validate", false, "Parse and validate the configuration but do not start the server")
	serveCmd.Flags().String("revoke", "", "Hostname for which to revoke the certificate")
	serveCmd.Flags().String("default-sni", certmagic.Default.DefaultServerName, "If a ClientHello ServerName is empty, use this ServerName to choose a TLS certificate")
	serveCmd.Flags().String("request-log", "", "Log destination for request logging")
	serveCmd.Flags().String("host", "*", "Hostname to serve (use \"*\" to serve all hostnames)")
	serveCmd.Flags().String("http-port", "80", "Port to use for HTTP")
	serveCmd.Flags().String("https-port", "443", "Port to use for HTTPS")
	serveCmd.Flags().String("bind", "", "IP to bind listener to (default \"0.0.0.0\")")
	serveCmd.Flags().Bool("tls", true, "Enable TLS on listener")
	serveCmd.Flags().Bool("tls-always-self-sign", false, "Always generate self signed certificate")
	serveCmd.Flags().Bool("tls-must-staple", false, "Enable TLS must staple")
	serveCmd.Flags().String("tls-protocols", "tls1.2 tls1.3", "Min and max TLS protocol")
	serveCmd.Flags().String("tls-cert-file", "", "Path to TLS certificate bundle (concatenation of the server's certificate followed by the CA's certificate chain)")
	serveCmd.Flags().String("tls-key-file", "", "Path to the server's private key file which matches the certificate bundle")
	serveCmd.Flags().String("hsts", "max-age=31536000;", "HTTP Strict Transport Security (default enabled when --host is given unless explicitly set to empty)")
	serveCmd.Flags().String("reverse-proxy-legacy-http", "", "URL to reverse proxy requests for Webapp and Z-Push")
	serveCmd.Flags().String("default-redirect", "", "URL to redirect to when no other path is given (/)")
	serveCmd.Flags().String("extra", "", "Extra configuration file (append at the end)")

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
	if catimeout, _ := cmd.Flags().GetString("catimeout"); catimeout != "" {
		caddyArgs = append(caddyArgs, "-catimeout", catimeout)
	}
	if http2, _ := cmd.Flags().GetBool("http2"); http2 {
		caddyArgs = append(caddyArgs, "-http2")
	}
	if quic, _ := cmd.Flags().GetBool("quic"); quic {
		caddyArgs = append(caddyArgs, "-quic")
	}
	if validate, _ := cmd.Flags().GetBool("validate"); validate {
		caddyArgs = append(caddyArgs, "-validate")
	}
	if revoke, _ := cmd.Flags().GetString("revoke"); revoke != "" {
		caddyArgs = append(caddyArgs, "-revoke", revoke)
	}
	if defaultSNI, _ := cmd.Flags().GetString("default-sni"); defaultSNI != "" {
		caddyArgs = append(caddyArgs, "-default-sni", defaultSNI)
	}

	root, _ := cmd.Flags().GetString("root")
	email, _ := cmd.Flags().GetString("email")
	host, _ := cmd.Flags().GetString("host")
	httpPort, _ := cmd.Flags().GetString("http-port")
	httpsPort, _ := cmd.Flags().GetString("https-port")
	bind, _ := cmd.Flags().GetString("bind")
	requestLog, _ := cmd.Flags().GetString("request-log")
	tls, _ := cmd.Flags().GetBool("tls")
	tlsAlwaysSelfSign, _ := cmd.Flags().GetBool("tls-always-self-sign")
	tlsMustStaple, _ := cmd.Flags().GetBool("tls-must-staple")
	tlsProtocols, _ := cmd.Flags().GetString("tls-protocols")
	tlsCertBundle, _ := cmd.Flags().GetString("tls-cert-file")
	tlsPrivateKey, _ := cmd.Flags().GetString("tls-key-file")
	hsts, _ := cmd.Flags().GetString("hsts")

	caddyArgs = append(caddyArgs, "-root", root, "-host", host, "-http-port", httpPort, "-https-port", httpsPort)
	if tls {
		caddyArgs = append(caddyArgs, "-port", httpsPort)
	} else {
		caddyArgs = append(caddyArgs, "-port", httpPort)
	}
	if email != "" {
		caddyArgs = append(caddyArgs, "-email", email)
	}

	reverseProxyLegacyHTTP, _ := cmd.Flags().GetString("reverse-proxy-legacy-http")
	defaultRedirect, _ := cmd.Flags().GetString("default-redirect")
	extraFile, _ := cmd.Flags().GetString("extra")

	// Configure underlying caddy.
	cfg := &config.Config{
		Root: root,

		Bind:  bind,
		Host:  host,
		Email: email,

		RequestLog: requestLog,

		TLSEnable:         tls,
		TLSAlwaysSelfSign: tlsAlwaysSelfSign,
		TLSCertBundle:     tlsCertBundle,
		TLSPrivateKey:     tlsPrivateKey,
		TLSMustStaple:     tlsMustStaple,
		TLSProtocols:      tlsProtocols,
		HSTS:              hsts,

		ReverseProxyLegacyHTTP: reverseProxyLegacyHTTP,
		DefaultRedirect:        defaultRedirect,
	}

	// Conditionals.
	if extraFile != "" {
		extra, err := ioutil.ReadFile(path.Clean(extraFile))
		if err != nil {
			return fmt.Errorf("failed to read extra file: %v", err)
		}
		cfg.Extra = extra
	}

	// Setup caddy.
	setupAssetsPath()

	// Set as default loader.
	caddy.SetDefaultCaddyfileLoader("default", defaultLoader(cfg))

	// Reset args, since caddymain has its own parsing.
	subArgs := append(os.Args[:1], caddyArgs...)
	os.Args = subArgs

	caddymain.Run()

	return nil
}

func defaultLoader(cfg *config.Config) caddy.LoaderFunc {
	return caddy.LoaderFunc(func(serverType string) (caddy.Input, error) {
		contents, err := config.Caddyfile(cfg)
		if err != nil {
			return nil, err
		}

		return caddy.CaddyfileInput{
			Contents:       contents,
			ServerTypeName: serverType,
		}, nil
	})
}

func setupAssetsPath() string {
	ap := os.Getenv("KOPANO_KWEB_ASSETS_PATH")
	if ap == "" {
		home := os.Getenv("HOME")
		ap = filepath.Join(home, ".kweb")
	}

	os.Setenv("CADDYPATH", ap)
	return ap
}
