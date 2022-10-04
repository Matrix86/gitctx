package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Matrix86/gitctx/internal/core"

	"github.com/jessevdk/go-flags"
)

var argOpts struct {
	//Help   bool   `short:"h" long:"help" description:"Show this menu."`
	Add       bool   `long:"add" description:"Create a new host in the selected config file."`
	Rm        string `long:"rm" description:"Remove an existing host in the selected config file."`
	SSHConfig string `short:"s" long:"sshconfig" description:"Set the path of the config (default: ~/.ssh/config)."`
	Hostname  string `long:"hostname" description:"Set the hostname to use for context change (default: github.com)."`
	Config    string `long:"config" description:"Set the path of the gitctx folder (default: ~/.gitctx)."`
}

var (
	Config         *core.Configuration
	currentCtxFile string
	configFilePath string
)

func main() {
	args, err := flags.Parse(&argOpts)
	if err != nil {
		panic(err)
	}

	if argOpts.SSHConfig == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(fmt.Sprintf("[!] error: %s\n", err))
		}
		argOpts.SSHConfig = fmt.Sprintf("%s/.ssh/config", home)
	}

	if argOpts.Hostname == "" {
		argOpts.Hostname = "github.com"
	}

	if argOpts.Config == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(fmt.Sprintf("[!] error: %s\n", err))
		}
		argOpts.Config = fmt.Sprintf("%s/.gitctx", home)
	}

	configFilePath = strings.Join([]string{argOpts.Config, "config.yml"}, "/")
	if _, err := os.Stat(argOpts.Config); os.IsNotExist(err) {
		os.MkdirAll(argOpts.Config, os.ModePerm)
		// creating empty config file
		err = core.CreateEmptyConfig(configFilePath)
		if err != nil {
			panic(fmt.Sprintf("[!] error: %s\n", err))
		}
	}

	Config, err = core.LoadConfiguration(configFilePath)
	if err != nil {
		panic(fmt.Sprintf("[!] error: loading configuration: %s\n", err))
	}

	currentCtxFile = strings.Join([]string{argOpts.Config, "context"}, "/")

	err = checkSSHConfig()
	if err != nil {
		panic(fmt.Sprintf("[!] error: %s\n", err))
	}

	if argOpts.Rm != "" {
		err := removeContext()
		if err != nil {
			panic(fmt.Sprintf("[!] error: %s\n", err))
		}
		return
	}

	if len(args) == 0 {
		// listing current github contexts
		err := listContexts()
		if err != nil {
			panic(fmt.Sprintf("[!] error: %s\n", err))
		}

	} else if len(args) == 1 {
		// changing the current github context
		err := setContext(args[0])
		if err != nil {
			panic(fmt.Sprintf("[!] error: %s\n", err))
		}
	}

	//fmt.Printf("%#v %v\n", argOpts, args)
}
