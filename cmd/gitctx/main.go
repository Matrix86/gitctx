package main

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Matrix86/gitctx/internal/core"

	"github.com/jessevdk/go-flags"
)

var argOpts struct {
	//Help   bool   `short:"h" long:"help" description:"Show this menu."`
	Add       bool   `long:"add" description:"Create a new host in the selected config file."`
	Rm        string `long:"rm" description:"Remove an existing host in the selected config file."`
	Edit      string `long:"edit" description:"Edit an existing host in the selected config file."`
	SSHConfig string `short:"s" long:"sshconfig" description:"Set the path of the config (default: ~/.ssh/config)."`
	Hostname  string `long:"hostname" description:"Set the hostname to use for context change (default: github.com)."`
	Config    string `long:"config" description:"Set the path of the gitctx folder (default: ~/.gitctx)."`
}

var (
	currentContexts *core.CurrentContexts
	Config          *core.Configuration
	currentCtxFile  string
	configFilePath  string

	//go:embed completion/gitctx.bash
	completionData []byte
)

func initDefaultConfig() {
	if argOpts.SSHConfig == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			core.Fatal("[!] error: %s\n", err)
		}
		argOpts.SSHConfig = fmt.Sprintf("%s/.ssh/config", home)
	}

	if argOpts.Hostname == "" {
		argOpts.Hostname = "github.com"
	}

	if argOpts.Config == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			core.Fatal("[!] error: %s", err)
		}
		argOpts.Config = fmt.Sprintf("%s/.gitctx", home)
	}
}

func main() {
	args, err := flags.Parse(&argOpts)
	if err != nil {
		switch flagsErr := err.(type) {
		case flags.ErrorType:
			if flagsErr == flags.ErrHelp {
				os.Exit(0)
			}
			os.Exit(1)
		default:
			os.Exit(1)
		}
	}

	initDefaultConfig()

	// Init the configuration folder if it doesn't exist
	configFilePath = strings.Join([]string{argOpts.Config, "config.yml"}, "/")
	if _, err := os.Stat(argOpts.Config); os.IsNotExist(err) {
		os.MkdirAll(argOpts.Config, os.ModePerm)
		// creating empty config file
		err = core.CreateEmptyConfig(configFilePath)
		if err != nil {
			core.Fatal("[!] error: %s", err)
		}

		// Completion file
		completionFile := strings.Join([]string{argOpts.Config, "gitctx.bash"}, "/")
		err = core.CreateFile(completionFile, completionData)
		if err != nil {
			core.Fatal("[!] error: %s", err)
		}
		core.Info("[@] add '. %s' to ~/.bashrc to enable the shell completion.", completionFile)
	}

	// Reading configuration's file
	Config, err = core.LoadConfiguration(configFilePath)
	if err != nil {
		core.Fatal("[!] error: loading configuration: %s", err)
	}

	currentCtxFile = strings.Join([]string{argOpts.Config, "context"}, "/")
	if _, err := os.Stat(currentCtxFile); errors.Is(err, os.ErrNotExist) {
		err = core.CreateEmptyFile(currentCtxFile)
		if err != nil {
			core.Fatal("[!] error: %s", err)
		}
	}

	currentContexts, err = core.LoadFromCurrentFile(currentCtxFile)
	if err != nil {
		core.Fatal("[!] error: loading current contexts: %s", err)
	}

	err = checkSSHConfig()
	if err != nil {
		core.Fatal("[!] error: %s", err)
	}

	if argOpts.Rm != "" {
		err := removeContext()
		if err != nil {
			core.Fatal("[!] error: %s", err)
		}
		return
	}

	if argOpts.Add {
		err := addContext()
		if err != nil {
			core.Fatal("[!] error: %s", err)
		}
		return
	}

	if argOpts.Edit != "" {
		err := editContext(argOpts.Edit)
		if err != nil {
			core.Fatal("[!] error: %s", err)
		}
		return
	}

	if len(args) == 0 {
		// listing current contexts
		err := listContexts()
		if err != nil {
			core.Fatal("[!] error: %s", err)
		}

	} else if len(args) == 1 {
		re := regexp.MustCompile(`([^=]+)`)
		matches := re.FindAllString(args[0], -1)
		if len(matches) == 2 {
			// renaming a context: gitctx oldname=newname
			if err = renameContext(matches[0], matches[1]); err != nil {
				core.Fatal("[!] error: %s", err)
			}
		} else {
			// changing the current context: gitctx ctxname
			err := setContext(args[0])
			if err != nil {
				core.Fatal("[!] error: %s", err)
			}
		}
	} else {
		core.Fatal("[!] error: command not supported")
	}
}
