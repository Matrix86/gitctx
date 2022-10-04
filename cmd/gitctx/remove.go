package main

import "fmt"

func removeContext() error {
	if _, ok := Config.Hosts[argOpts.Rm]; ok {
		delete(Config.Hosts, argOpts.Rm)
	} else {
		return fmt.Errorf("context %s not found", argOpts.Rm)
	}

	// set one of the remaining contexts
	for n := range Config.Hosts {
		err := setContext(n)
		if err != nil {
			return fmt.Errorf("updating context %s: %s", n, err)
		}
		break
	}

	// write updated configuration
	err := Config.WriteConfiguration(configFilePath)
	if err != nil {
		return fmt.Errorf("writing updated configuration: %s", err)
	}
	return nil
}
