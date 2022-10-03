package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

var argOpts struct {
	//Help   bool   `short:"h" long:"help" description:"Show this menu."`
	Add    string `long:"add" description:"Create a new host in the selected config file."`
	Rm     string `long:"rm" description:"Remove an existing host in the selected config file."`
	Config string `short:"c" long:"config" description:"Set the path of the config (default: ~/.ssh/config)."`
}

func main() {
	args, err := flags.Parse(&argOpts)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", argOpts.Config)
	if argOpts.Config == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			panic(fmt.Sprintf("[!] error: %s\n", err))
		}
		argOpts.Config = fmt.Sprintf("%s/.ssh/config", home)
	}

	if len(args) == 0 {
		// listing current github contexts
		listContexts()
	} else if len(args) == 1 {
		// changing the current github context
		setContext(args[0])
	}

	fmt.Printf("%#v %v\n", argOpts, args)
}
