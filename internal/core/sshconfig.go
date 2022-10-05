package core

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/kevinburke/ssh_config"
)

func GetConfigHosts(configPath string) ([]Context, error) {
	ctxs := []Context{}
	f, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("can't open %s: %s", configPath, err)
	}
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

func commentLines(filename string, lines map[int]bool) error {
	var buf strings.Builder

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)

	cline := 1
	for fileScanner.Scan() {
		if _, ok := lines[cline]; ok {
			buf.WriteString("# ")
		}
		buf.WriteString(fmt.Sprintf("%s\n", fileScanner.Text()))
		cline++
	}

	err = ioutil.WriteFile(filename, []byte(buf.String()), 0)
	if err != nil {
		return err
	}
	return nil
}

func DeleteConfigHosts(configPath string, hosts []string) error {
	lines := map[int]bool{}
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

		clines := map[int]bool{}

		for _, node := range host.Nodes {
			clines[node.Pos().Line] = true
		}

		for _, h := range hosts {
			if host.Patterns[0].String() == h {
				for l, _ := range clines {
					lines[l] = true
				}
			}
		}
	}

	if err := commentLines(configPath, lines); err != nil {
		return fmt.Errorf("can't comment lines on %s: %s", configPath, err)
	}

	return nil
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

	f, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("opening the file %s: %s", configPath, err)
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("reading the file %s: %s", configPath, err)
	}
	f.Close()

	// Include not found, let's add it
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("Include %s\n\n", include))
	buf.WriteString(string(data))

	f, err = os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("opening the file %s: %s", configPath, err)
	}
	defer f.Close()

	f.Truncate(0)

	_, err = f.WriteString(buf.String())
	if err != nil {
		return fmt.Errorf("writing on %s: %s", configPath, err)
	}

	return nil
}
