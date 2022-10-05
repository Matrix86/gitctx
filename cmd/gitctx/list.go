package main

import (
	"fmt"

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
	currentIdentity := ""
	for _, ctx := range *currentContexts {
		if ctx.Host == argOpts.Hostname {
			currentIdentity = ctx.IdentityFile
		}
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
