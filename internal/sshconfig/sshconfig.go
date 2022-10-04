package sshconfig

import (
	"fmt"
	"os"
	"strings"

	"github.com/kevinburke/ssh_config"
)

type Context struct {
	Name         string
	Host         string
	User         string
	IdentityFile string
	Line         int
}

func GetConfigHosts(configPath string) ([]Context, error) {
	ctxs := []Context{}
	f, _ := os.Open(configPath)
	defer f.Close()
	cfg, _ := ssh_config.Decode(f)
	for _, host := range cfg.Hosts {
		// ignore patterns that contains multiple patterns or wildcards
		if len(host.Patterns) > 1 {
			continue
		}
		if strings.ContainsAny(host.Patterns[0].String(), "!* ") {
			continue
		}

		ctx := Context{}

		for _, node := range host.Nodes {
			switch t := node.(type) {
			case *ssh_config.Empty:
				continue

			case *ssh_config.KV:
				// "keys are case insensitive" per the spec
				lkey := strings.ToLower(t.Key)
				if lkey == "hostname" {
					ctx.Name = strings.ToLower(host.Patterns[0].String())
					ctx.Host = t.Value
					ctx.Line = node.Pos().Line
				} else if lkey == "user" {
					ctx.User = t.Value
				} else if lkey == "identityfile" {
					ctx.IdentityFile = t.Value
				}

			default:
				continue
			}
		}

		if ctx.Name != "" {
			ctxs = append(ctxs, ctx)
		}
	}

	return ctxs, nil
}

func AddInclude(configPath string, include string) error {
	f, _ := os.Open(configPath)
	defer f.Close()
	cfg, _ := ssh_config.Decode(f)
	for _, host := range cfg.Hosts {
		for _, node := range host.Nodes {
			switch t := node.(type) {
			case *ssh_config.Include:
				if strings.Contains(t.String(), include) {
					return nil
				}

			default:
				continue
			}
		}
	}

	// Include not found, let's add it
	f, err := os.OpenFile(configPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("opening the file %s: %s", configPath, err)
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("Include %s\n", include))
	if err != nil {
		return fmt.Errorf("adding the Include directive: %s", err)
	}

	return nil
}

func (ctx *Context) WriteOnFile(filename string) error {
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("Host %s\n", ctx.Name))
	buf.WriteString(fmt.Sprintf("HostName %s\n", ctx.Host))
	if ctx.User != "" {
		buf.WriteString(fmt.Sprintf("User %s\n", ctx.User))
	}
	buf.WriteString(fmt.Sprintf("IdentityFile %s\n", ctx.IdentityFile))

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("opening the file %s: %s", filename, err)
	}
	defer f.Close()

	f.Truncate(0)

	_, err = f.WriteString(buf.String())
	if err != nil {
		return fmt.Errorf("adding the Include directive: %s", err)
	}
	return nil
}
