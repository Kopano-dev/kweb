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

	// Plugins.
	_ "github.com/captncraig/caddy-realip"
	_ "github.com/miekg/caddy-prometheus"
	_ "github.com/pyed/ipfilter"
	_ "github.com/xuqingfeng/caddy-rate-limit"
)

func init() {
	caddymain.EnableTelemetry = false
}

func commandCaddy() *cobra.Command {
	caddyCmd := &cobra.Command{
		Use:   "caddy [...args]",
		Short: "Start caddy and listen for requests",
		Run: func(cmd *cobra.Command, args []string) {
			if err := caddy(cmd, args); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		},
		DisableFlagParsing: true,
	}

	return caddyCmd
}

func caddy(cmd *cobra.Command, args []string) error {
	// Reset args, since caddymain has its own parsing.
	subArgs := append(os.Args[:1], args...)
	os.Args = subArgs

	caddymain.Run()

	return nil
}
