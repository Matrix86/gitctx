package main

import (
	"fmt"

	"github.com/Matrix86/gitctx/internal/core"
)

func setContext(ctxName string) error {
	if c, ok := Config.Hosts[ctxName]; !ok {
		return fmt.Errorf("context %s not found in the configuration file", ctxName)
	} else {
		ctx := &core.Context{
			Name:         ctxName,
			Host:         c.Hostname,
			User:         c.User,
			IdentityFile: c.IdentityFile,
			Line:         0,
		}
		(*currentContexts)[c.Hostname] = ctx

		err := currentContexts.WriteOnFile(currentCtxFile)
		if err != nil {
			return fmt.Errorf("writing context: %s", err)
		}
	}
	return nil
}
