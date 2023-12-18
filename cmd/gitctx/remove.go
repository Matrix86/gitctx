package main

import "fmt"

func removeContext() error {
	if _, ok := Config.Hosts[argOpts.Rm]; ok {
		delete(Config.Hosts, argOpts.Rm)
	} else {
		return fmt.Errorf("context %s not found", argOpts.Rm)
	}
	delete(Config.GitSettings, argOpts.Rm)

	// set one of the remaining contexts
	for n, ctx := range *currentContexts {
		if ctx.Name == argOpts.Rm {
			delete(*currentContexts, n)
			currentContexts.WriteOnFile(currentCtxFile)
			break
		}
	}

	// write updated configuration
	err := Config.WriteConfiguration(configFilePath)
	if err != nil {
		return fmt.Errorf("writing updated configuration: %s", err)
	}
	return nil
}
