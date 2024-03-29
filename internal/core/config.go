package core

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"gopkg.in/yaml.v3"
)

type Host struct {
	Hostname     string `yaml:"hostname"`
	User         string `yaml:"user"`
	IdentityFile string `yaml:"identity_file"`
}

type GitSettings struct {
	Name          string `yaml:"name"`
	Email         string `yaml:"email"`
	SigningKey    string `yaml:"signingkey"`
	CommitGPGSign string `yaml:"commit_gpgsign"`
}

// Configuration contains all the configs read by yaml file
type Configuration struct {
	DefaultHostname string                       `yaml:"default_hostname"`
	Hosts           map[string]Host              `yaml:"hosts"`
	GitSettings     map[string]map[string]string `yaml:"git_settings"`
}

func (c *Configuration) AddHost(name, hostname, user, IdentityFile string) {
	c.Hosts[name] = Host{
		Hostname:     hostname,
		User:         user,
		IdentityFile: IdentityFile,
	}
}

func (c *Configuration) WriteConfiguration(filename string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, data, 0)
	if err != nil {
		return err
	}
	return nil
}

// LoadConfiguration create a Configuration struct from a filename
func LoadConfiguration(path string) (*Configuration, error) {
	configuration := &Configuration{}

	if err := cleanenv.ReadConfig(path, configuration); err != nil {
		return configuration, err
	}

	if configuration.Hosts == nil {
		configuration.Hosts = map[string]Host{}
	}
	return configuration, nil
}

func CreateEmptyConfig(filename string) error {
	content := []byte("default_hostname: \"github.com\"\nhosts:\n")
	return os.WriteFile(filename, content, 0644)
}

func CreateEmptyFile(filename string) error {
	return os.WriteFile(filename, []byte{}, 0644)
}

func CreateFile(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0644)
}
