package core

import (
	"fmt"
	"os"
	"strings"
)

type CurrentContexts map[string]*Context

type Context struct {
	Name         string
	Host         string
	User         string
	IdentityFile string
	Line         int
}

func (ctx *Context) String() string {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("Host %s\n", ctx.Name))
	buf.WriteString(fmt.Sprintf("  HostName %s\n", ctx.Host))
	if ctx.User != "" {
		buf.WriteString(fmt.Sprintf("  User %s\n", ctx.User))
	}
	buf.WriteString(fmt.Sprintf("  IdentityFile %s\n\n", ctx.IdentityFile))

	return buf.String()
}

func (c *CurrentContexts) WriteOnFile(filename string) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("opening the file %s: %s", filename, err)
	}
	defer f.Close()

	f.Truncate(0)

	for _, ctx := range *c {
		_, err = f.WriteString(ctx.String())
		if err != nil {
			return fmt.Errorf("writing the new context: %s", err)
		}
	}

	return nil
}

func (c CurrentContexts) AddContext(ctx *Context) {
	c[ctx.Host] = ctx
}

func (c CurrentContexts) RemoveContext(ctx *Context) {
	delete(c, ctx.Host)
}

func LoadFromCurrentFile(filename string) (*CurrentContexts, error) {
	contexts, err := GetConfigHosts(filename)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %s", filename, err)
	}

	currentCtxs := CurrentContexts{}
	for _, ctx := range contexts {
		currentCtxs[ctx.Host] = &ctx
	}

	return &currentCtxs, nil
}
