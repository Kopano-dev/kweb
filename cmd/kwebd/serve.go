/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/caddyserver/caddy"
	"github.com/caddyserver/caddy/caddy/caddymain"
	"github.com/mholt/certmagic"
	"github.com/spf13/cobra"

	"stash.kopano.io/kgol/kweb/config"
	"stash.kopano.io/kgol/kweb/version"
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
	serveCmd.Flags().Bool("log-timestamps", true, "Enable timestamps for the process log")
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
	serveCmd.Flags().String("extra", "", "Path to extra configuration file or folder with .cfg files, separate multiple with : (ealier entries have priority)")
	serveCmd.Flags().String("base", "", "Path to Base configuration file or folder with .cfg files, separate multiple with : (ealier entries have priority)")
	serveCmd.Flags().String("snippets", "", "Path to snippets configuration file or folder with .cfg files, separate multiple with : (ealier entries have priority)")

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
	logTimestamps, _ := cmd.Flags().GetBool("log-timestamps")
	tls, _ := cmd.Flags().GetBool("tls")
	tlsAlwaysSelfSign, _ := cmd.Flags().GetBool("tls-always-self-sign")
	tlsMustStaple, _ := cmd.Flags().GetBool("tls-must-staple")
	tlsProtocols, _ := cmd.Flags().GetString("tls-protocols")
	tlsCertBundle, _ := cmd.Flags().GetString("tls-cert-file")
	tlsPrivateKey, _ := cmd.Flags().GetString("tls-key-file")
	hsts, _ := cmd.Flags().GetString("hsts")

	if !logTimestamps {
		// Disable timestamps for logging.
		log.SetFlags(0)
	}

	caddyArgs = append(caddyArgs, "-root", root, "-host", host, "-http-port", httpPort, "-https-port", httpsPort, "-log", "stdout")
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

	snippets, _ := getStringFlagOrEnv(cmd, "snippets", "KOPANO_KWEB_CFG_SNIPPETS_PATH")
	base, _ := getStringFlagOrEnv(cmd, "base", "KOPANO_KWEB_CFG_BASE_PATH")
	extra, _ := getStringFlagOrEnv(cmd, "extra", "KOPANO_KWEB_CFG_EXTRA_PATH")

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

	// Snippets.
	if snippets != "" {
		var b bytes.Buffer
		if err := loadD("snippets", snippets, &b, cfg); err != nil {
			return err
		}
		cfg.Snippets = b.Bytes()
	}

	// Base.
	if base != "" {
		var b bytes.Buffer
		if err := loadD("base", base, &b, cfg); err != nil {
			return err
		}
		cfg.Base = b.Bytes()
	}

	// Extra..
	if extra != "" {
		var b bytes.Buffer
		if err := loadD("extra", extra, &b, cfg); err != nil {
			return err
		}
		cfg.Extra = b.Bytes()
	}

	// Setup caddy.
	setupAssetsPath()

	// Set as default loader.
	caddy.SetDefaultCaddyfileLoader("default", defaultLoader(cfg))

	// Reset args, since caddymain has its own parsing.
	subArgs := append(os.Args[:1], caddyArgs...)
	os.Args = subArgs

	log.Printf("[INFO] Kweb version: %v\n", version.Version)
	caddymain.Run()

	return nil
}

func getStringFlagOrEnv(cmd *cobra.Command, name, env string) (string, error) {
	v, _ := cmd.Flags().GetString(name)
	if v == "" && env != "" {
		v = os.Getenv(env)
	}
	return v, nil
}

func defaultLoader(cfg *config.Config) caddy.LoaderFunc {
	return caddy.LoaderFunc(func(serverType string) (caddy.Input, error) {
		contents, err := config.Caddyfile(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to load generated configuration: %w", err)
		}

		if os.Getenv("KOPANO_KWEB_DUMP_INTERNAL_CADDYFILE") != "" {
			fmt.Println("----- internal configuration start -----")
			fmt.Printf("%s\n", contents)
			fmt.Println("----- internal configuration end -----")
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

// byFileName is a custom sorter to ensure numeric file configuration files
// from multiple directories are sorted with lowest numbers first.
type byFileName []string

func (a byFileName) Len() int      { return len(a) }
func (a byFileName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byFileName) Less(i, j int) bool {
	ifn := filepath.Base(a[i])
	jfn := filepath.Base(a[j])

	var in int64 = -1
	ins := strings.SplitN(ifn, "-", 2)
	if len(ins) > 1 {
		in, _ = strconv.ParseInt(ins[0], 10, 64)
	}

	var jn int64 = -1
	jns := strings.SplitN(jfn, "-", 2)
	if len(jns) > 1 {
		jn, _ = strconv.ParseInt(jns[0], 10, 64)
	}

	if in >= 0 {
		if jn >= 0 {
			if in != jn {
				return in < jn
			}
			ifn = ins[1]
			jfn = jns[1]
		} else {
			return true
		}
	}
	return ifn < jfn
}

func loadD(name, pathString string, b *bytes.Buffer, context *config.Config) error {
	paths := strings.Split(pathString, ":")

	reader := func(fn string) error {
		if t, err := template.ParseFiles(fn); err != nil {
			return fmt.Errorf("failed to read %s file %s: %w", name, fn, err)
		} else {
			b.WriteString(fmt.Sprintf("# <-- %s \n", fn))
			if err := t.Execute(b, context); err != nil {
				return fmt.Errorf("failed to process %s file %s:%w", name, fn, err)
			}
			b.WriteString(fmt.Sprintf("# --> %s end\n\n", fn))
		}
		return nil
	}

	load := []string{}
	seen := make(map[string]bool)

	for _, p := range paths {
		p = path.Clean(p)
		// Extra can either be a file or folder.
		stat, err := os.Stat(p)
		if err != nil {
			return fmt.Errorf("failed to read %s configuration: %w", name, err)
		}
		if stat.IsDir() {
			// Add all files in alphabetical order.
			files, err := ioutil.ReadDir(p)
			if err != nil {
				return fmt.Errorf("failed to read %s configuration directory: %w", name, err)
			}
			for _, f := range files {
				fn := f.Name()
				if filepath.Ext(fn) != ".cfg" {
					continue
				}
				if seen[fn] {
					// Skip files which were already seen.
					continue
				}
				seen[fn] = true
				fn = filepath.Join(p, fn)
				load = append(load, fn)
			}
		} else {
			// Add configured file directly.
			load = append(load, p)
		}
	}

	sort.Sort(byFileName(load))
	for _, f := range load {
		if err := reader(f); err != nil {
			return err
		}
	}

	return nil
}
