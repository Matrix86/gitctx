package main

import (
	"fmt"

	"github.com/Matrix86/gitctx/internal/sshconfig"
	"github.com/fatih/color"
	"github.com/kevinburke/ssh_config"
)

type Context struct {
	Name    string
	Host    string
	Current bool
	Node    *ssh_config.Host
}

func listContexts() error {
	contexts, err := sshconfig.GetConfigHosts(currentCtxFile)
	if err != nil {
		return fmt.Errorf("reading %s: %s", currentCtxFile, err)
	}
	currentIdentity := ""
	if len(contexts) > 0 {
		// the context file should contain a single config
		currentIdentity = contexts[0].IdentityFile
	}

	for name, ctx := range Config.Hosts {
		if ctx.Hostname != argOpts.Hostname {
			continue
		}
		if ctx.IdentityFile == currentIdentity {
			lineColor := color.New()
			lineColor.Add(color.Bold)
			lineColor.Add(color.FgYellow)
			lineColor.Add(color.BgWhite)
			name = lineColor.Sprintf(name)
		}

		fmt.Printf("%s\n", name)
	}

	return nil
}
