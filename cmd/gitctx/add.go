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
	gitGlobal, err := askData("Do you want to specify git global settings? [y/N]", "(y|Y|n|N)?")
	if err != nil {
		return err
	}
	if gitGlobal == "y" || gitGlobal == "Y" {
		gitGlobalEmail, err := askData("user.email", "[a-zA-Z0-9_\\.\\-+\\\\/\\~@]+")
		if err != nil {
			return err
		}
		gitGlobalName, err := askData("user.name", "[a-zA-Z0-9_\\.\\-+\\\\/\\~@]+")
		if err != nil {
			return err
		}
		gitGlobalSigningKey, err := askData("user.signingkey", "[a-zA-Z0-9_\\.\\-+\\\\/\\~@]+")
		if err != nil {
			return err
		}
		Config.GitSettings[name] = core.GitSettings{
			Name:       gitGlobalName,
			Email:      gitGlobalEmail,
			SigningKey: gitGlobalSigningKey,
		}
	}

	ctx := core.Host{
		Hostname:     hostname,
		User:         user,
		IdentityFile: identity,
	}
	Config.Hosts[name] = ctx
	Config.WriteConfiguration(configFilePath)
	return nil
}
