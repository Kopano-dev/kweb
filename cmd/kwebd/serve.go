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
		DisableFlagParsing: true,
	}

	return serveCmd
}

func serve(cmd *cobra.Command, args []string) error {
	caddy.SetDefaultCaddyfileLoader("default", caddy.LoaderFunc(defaultLoader))

	// Reset args, since caddymain has its own parsing.
	subArgs := append(os.Args[:1], args...)
	os.Args = subArgs

	caddymain.Run()

	return nil
}

func defaultLoader(serverType string) (caddy.Input, error) {
	contents, err := kweb.Caddyfile()
	if err != nil {
		return nil, err
	}

	return caddy.CaddyfileInput{
		Contents:       contents,
		ServerTypeName: serverType,
	}, nil
}
