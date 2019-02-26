/*
 * SPDX-License-Identifier: AGPL-3.0-or-later
 * Copyright 2018 Kopano and its licensors
 */

package main

import (
	"fmt"
	"os"

	"github.com/mholt/caddy/caddy/caddymain"
	"github.com/spf13/cobra"

	// 3rd party plugins.
	_ "github.com/captncraig/caddy-realip"
	_ "github.com/pyed/ipfilter"
	_ "github.com/xuqingfeng/caddy-rate-limit"
	_ "stash.kopano.io/kgol/caddy-prometheus"

	// Our plugins.
	_ "stash.kopano.io/kgol/kweb/caddy-alias"
	_ "stash.kopano.io/kgol/kweb/caddy-configjson"
	_ "stash.kopano.io/kgol/kweb/caddy-folderish"
	_ "stash.kopano.io/kgol/kweb/caddy-staticpwa"
)

func init() {
	caddymain.EnableTelemetry = false
}

func commandCaddy() *cobra.Command {
	caddyCmd := &cobra.Command{
		Use:   "caddy [...args]",
		Short: "Start like caddy and listen for requests",
		Run: func(cmd *cobra.Command, args []string) {
			if err := caddyCommandHandler(cmd, args); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		},
		DisableFlagParsing: true,
	}

	return caddyCmd
}

func caddyCommandHandler(cmd *cobra.Command, args []string) error {
	// Reset args, since caddymain has its own parsing.
	subArgs := append(os.Args[:1], args...)
	os.Args = subArgs

	caddymain.Run()

	return nil
}
