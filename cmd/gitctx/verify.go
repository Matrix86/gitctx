package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Matrix86/gitctx/internal/core"
)

func askYesNo(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/N]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			return false
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" || response == "" {
			return false
		}
	}
}

// if there are more configurations for the hostname we need to move them to a secondary config
func checkSSHConfig() error {
	hosts, err := core.GetConfigHosts(argOpts.SSHConfig)
	if err != nil {
		return err
	}

	err = nil
	found := []string{}
	for _, host := range hosts {
		if host.Host == argOpts.Hostname {
			found = append(found, fmt.Sprintf("%d", host.Line))
		}
	}
	if len(found) > 0 {
		if askYesNo(fmt.Sprintf("found configurations for Host '%s' on lines %s\nDo you want to move them to gitctx configuration?", argOpts.Hostname, strings.Join(found, ", "))) {
			backupFile, err := doBackup(argOpts.SSHConfig)
			if err != nil {
				return err
			}

			fmt.Printf("Created a backup for %s to %s.\n", argOpts.SSHConfig, backupFile)

			host2Remove := []string{}
			for _, host := range hosts {
				if host.Host == argOpts.Hostname {
					host2Remove = append(host2Remove, host.Name)
					Config.AddHost(host.Name, host.Host, host.User, host.IdentityFile)
				}
			}

			if err := core.DeleteConfigHosts(argOpts.SSHConfig, host2Remove); err != nil {
				return fmt.Errorf("removing hosts from %s: %s", argOpts.SSHConfig, err)
			}

			if err = Config.WriteConfiguration(configFilePath); err != nil {
				return fmt.Errorf("can't write config file %s: %s", configFilePath, err)
			}

			fmt.Printf("Configuration imported on %s.\n", configFilePath)

		} else {
			return fmt.Errorf("stopping")
		}
	}

	err = core.AddInclude(argOpts.SSHConfig, currentCtxFile)
	if err != nil {
		return err
	}
	return nil
}

func doBackup(filename string) (string, error) {
	sourceFileStat, err := os.Stat(filename)
	if err != nil {
		return "", err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return "", fmt.Errorf("%s is not a regular file", filename)
	}

	src, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer src.Close()

	fname := filepath.Base(filename)
	fpath := filepath.Dir(filename)
	dest, err := ioutil.TempFile(fpath, fname)
	if err != nil {
		return "", err
	}
	defer dest.Close()
	io.Copy(dest, src)

	return dest.Name(), nil
}
