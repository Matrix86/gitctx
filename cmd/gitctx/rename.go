package main

import "fmt"

func renameContext(from string, to string) error {
	if _, ok := Config.Hosts[from]; !ok {
		return fmt.Errorf("context %s not found", from)
	}
	if _, ok := Config.Hosts[to]; ok {
		return fmt.Errorf("context %s already exists", to)
	}
	Config.Hosts[to] = Config.Hosts[from]
	delete(Config.Hosts, from)

	if old, ok := Config.GitSettings[from]; ok {
		Config.GitSettings[to] = old
		delete(Config.GitSettings, from)
	}

	// set one of the remaining contexts
	for n, ctx := range *currentContexts {
		if ctx.Name == from {
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
