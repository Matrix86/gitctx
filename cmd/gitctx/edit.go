package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Matrix86/gitctx/internal/core"
)

var alphaNumEmptyUnderscorePattern = regexp.MustCompile("[a-zA-Z0-9_.]*")

func askDataWithDefault(data string, pattern string, def string) (string, error) {
	var re *regexp.Regexp
	var err error

	if pattern == "" {
		re = alphaNumEmptyUnderscorePattern
	} else {
		re, err = regexp.Compile(pattern)
		if err != nil {
			return "", fmt.Errorf("can't compile pattern '%s': %s", pattern, err)
		}
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [%s]: ", data, def)

		response, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if re.MatchString(response) {
			if response == "" {
				return def, nil
			}
			return response, nil
		} else {
			fmt.Printf("[!] '%s' uses invalid chars\n", response)
		}
	}
}

func editContext(ctxName string) error {
	var ctx core.Host
	var ok bool
	if ctx, ok = Config.Hosts[ctxName]; !ok {
		return fmt.Errorf("context %s does not exist", ctxName)
	}

	hostname, err := askDataWithDefault("Hostname", "", ctx.Hostname)
	if err != nil {
		return err
	}

	user, err := askDataWithDefault("User", "", ctx.User)
	if err != nil {
		return err
	}

	identity, err := askDataWithDefault("IdentityFile", "[a-zA-Z0-9_\\.\\-\\/\\~]*", ctx.IdentityFile)
	if err != nil {
		return err
	}

	// git global settings
	gitGlobal, err := askData("Do you want to edit git global settings? [y/N]", "(y|Y|n|N)?")
	if err != nil {
		return err
	}
	if gitGlobal == "y" || gitGlobal == "Y" {
		gitCtx := Config.GitSettings[ctxName]

		gitGlobalEmail, err := askDataWithDefault("user.email", "[a-zA-Z0-9_\\.\\-+\\\\/\\~@]*", gitCtx.Email)
		if err != nil {
			return err
		}
		gitGlobalName, err := askDataWithDefault("user.name", "[a-zA-Z0-9_\\.\\-+\\\\/\\~@]*", gitCtx.Name)
		if err != nil {
			return err
		}
		gitGlobalSigningKey, err := askDataWithDefault("user.signingkey", "[a-zA-Z0-9_\\.\\-+\\\\/\\~@]*", gitCtx.SigningKey)
		if err != nil {
			return err
		}
		Config.GitSettings[ctxName] = core.GitSettings{
			Name:       gitGlobalName,
			Email:      gitGlobalEmail,
			SigningKey: gitGlobalSigningKey,
		}
	}

	newCtx := core.Host{
		Hostname:     hostname,
		User:         user,
		IdentityFile: identity,
	}
	Config.Hosts[ctxName] = newCtx
	Config.WriteConfiguration(configFilePath)
	return nil
}
