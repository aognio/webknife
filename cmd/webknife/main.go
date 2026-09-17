// Package main is the composition root for webknife.
//
// It wires the CLI parser to the application layer, translating between
// the two Config types. This is the only place where adapters meet.
package main

import (
	"os"

	"github.com/aognio/webknife/internal/adapters/cli"
	"github.com/aognio/webknife/internal/application"
)

func main() {
	cmd, cliCfg, err := cli.Parse(os.Args, os.Stdout, os.Stderr)
	if err != nil {
		os.Exit(1)
	}
	if cmd == "" {
		return
	}

	appCfg := &application.Config{
		Command:               cmd,
		ListenAddr:            cliCfg.ListenAddr,
		Root:                  cliCfg.Root,
		Upstream:              cliCfg.Upstream,
		Auth:                  cliCfg.Auth,
		TLSCert:               cliCfg.TLSCert,
		TLSKey:                cliCfg.TLSKey,
		LogFormat:             cliCfg.LogFormat,
		Verbose:               cliCfg.Verbose,
		EchoMaxBody:           cliCfg.EchoMaxBody,
		RespondStatusCode:     cliCfg.RespondStatusCode,
		RespondBody:           cliCfg.RespondBody,
		RespondBodyFile:       cliCfg.RespondBodyFile,
		RespondContentType:    cliCfg.RespondContentType,
		RedirectTarget:        cliCfg.RedirectTarget,
		RedirectStatusCode:    cliCfg.RedirectStatusCode,
		RedirectPreservePath:  cliCfg.RedirectPreservePath,
		SetRequestHeaders:     cliCfg.SetRequestHeaders,
		AddRequestHeaders:     cliCfg.AddRequestHeaders,
		RemoveRequestHeaders:  cliCfg.RemoveRequestHeaders,
		SetResponseHeaders:    cliCfg.SetResponseHeaders,
		AddResponseHeaders:    cliCfg.AddResponseHeaders,
		RemoveResponseHeaders: cliCfg.RemoveResponseHeaders,
	}

	if err := application.Run(appCfg); err != nil {
		os.Exit(1)
	}
}
