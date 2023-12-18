package core

import (
	"fmt"
	"os/exec"
)

func ExecCommand(cmdName string, params []string) error {
	cmd := exec.Command(cmdName, params...)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("can't exec %s: %s", cmdName, err)
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if exiterr.Success() {
				return nil
			} else {
				return fmt.Errorf("%s exited with: %s", cmdName, exiterr.Error())
			}
		} else {
			return fmt.Errorf("during cmd waiting for %s: %s", cmdName, err)
		}
	}
	return nil
}
