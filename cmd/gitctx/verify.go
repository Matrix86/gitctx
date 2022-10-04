package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Matrix86/gitctx/internal/sshconfig"
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
	hosts, err := sshconfig.GetConfigHosts(argOpts.SSHConfig)
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
			// TODO: read and move hosts to config
			err = doBackup(argOpts.SSHConfig)
			if err != nil {
				return err
			}

		} else {
			return fmt.Errorf("stopping")
		}
	}

	err = sshconfig.AddInclude(argOpts.SSHConfig, currentCtxFile)
	if err != nil {
		return err
	}
	return nil
}

func doBackup(filename string) error {
	sourceFileStat, err := os.Stat(filename)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", filename)
	}

	src, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer src.Close()

	fname := filepath.Base(filename)
	fpath := filepath.Dir(filename)
	dest, err := ioutil.TempFile(fpath, fname)
	if err != nil {
		return err
	}
	defer dest.Close()
	_, err = io.Copy(dest, src)

	return err
}
