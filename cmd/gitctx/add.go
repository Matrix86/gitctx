package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Matrix86/gitctx/internal/core"
)

var alphaNumUnderscorePattern = regexp.MustCompile("[a-zA-Z0-9_.]+")

func askData(data string, pattern string) (string, error) {
	var re *regexp.Regexp
	var err error

	if pattern == "" {
		re = alphaNumUnderscorePattern
	} else {
		re, err = regexp.Compile(pattern)
		if err != nil {
			return "", fmt.Errorf("can't compile pattern '%s': %s", pattern, err)
		}
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s: ", data)

		response, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if re.MatchString(response) {
			return response, nil
		} else {
			fmt.Printf("[!] '%s' uses invalid chars\n", response)
		}
	}
}

func addContext() error {
	name, err := askData("Name", "")
	if err != nil {
		return err
	}

	if _, ok := Config.Hosts[name]; ok {
		return fmt.Errorf("context %s already exists", name)
	}

	hostname, err := askData("Hostname", "")
	if err != nil {
		return err
	}

	user, err := askData("User", "")
	if err != nil {
		return err
	}

	identity, err := askData("IdentityFile", "[a-zA-Z0-9_.-\\/~]+")
	if err != nil {
		return err
	}

	// git global settings
	fmt.Printf("It is possible to specify git configuration to be launched changing manually the file %s\n", configFilePath)

	ctx := core.Host{
		Hostname:     hostname,
		User:         user,
		IdentityFile: identity,
	}
	Config.Hosts[name] = ctx
	Config.WriteConfiguration(configFilePath)
	return nil
}
