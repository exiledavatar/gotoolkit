package client

import (
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Client is an attempt to establish a standard way of configuring
// sql clients, and using struct tags and methods when interacting with them.
// This is very much a work in progress... as I'm simply trying to reduce the
// amount of redundant code...
type Client[T any] struct {
	Config    Config
	Templator Templator
	Conn      *T
}

type Config struct {
	Connection ConnectionConfig
	Template   TemplatorConfig
}

func NewConfig() Config {
	return Config{
		Connection: NewConnectionConfig(),
		Template:   NewTemplatorConfig(),
	}
}

type ConnectionConfig struct {
	Name               string
	Type               string
	Driver             string
	DataSourceName     string
	Host               string
	Port               int
	Protocol           string
	Options            map[string]string
	ConnectionString   string
	Database           string
	Schema             string
	ExpandEnvVars      bool
	ExpandFileContents bool
	Username           string
	Password           string
}

func NewConnectionConfig(cfg ...ConnectionConfig) ConnectionConfig {
	cc := &ConnectionConfig{}
	cc.Merge(cfg...)
	return *cc
}

func (cc *ConnectionConfig) Merge(cfg ...ConnectionConfig) *ConnectionConfig {

	for _, cf := range cfg {
		if cf.Name != "" {
			cc.Name = cf.Name
		}
		if cf.Type != "" {
			cc.Type = cf.Type
		}
		if cf.Driver != "" {
			cc.Driver = cf.Driver
		}
		if cf.DataSourceName != "" {
			cc.DataSourceName = cf.DataSourceName
		}
		if cf.Host != "" {
			cc.Host = cf.Host
		}
		if cf.Port != 0 {
			cc.Port = cf.Port
		}
		if cf.Protocol != "" {
			cc.Protocol = cf.Protocol
		}
		if cf.Options != nil {
			cc.Options = cf.Options
		}
		if cf.ConnectionString != "" {
			cc.ConnectionString = cf.ConnectionString
		}
		if cf.Database != "" {
			cc.Database = cf.Database
		}
		if cf.Schema != "" {
			cc.Schema = cf.Schema
		}
		cc.ExpandEnvVars = cf.ExpandEnvVars
		cc.ExpandFileContents = cf.ExpandFileContents
		if cf.Username != "" {
			cc.Username = cf.Username
		}
		if cf.Password != "" {
			cc.Password = cf.Password
		}
	}

	return cc
}

func LoadAndParseConfig(filename string) (*Config, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	b = ExpandEnvVars(b)
	cfg := Config{}
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// ExpandEnvVars substitutes environment variables of the form ${ENV_VAR_NAME}
// if you have characters that need to be escaped, they should be surrounded in
// quotes in the source string.
func ExpandEnvVars[T []byte | string](value T) T {
	s := string(value)

	re := regexp.MustCompile(`\$\{.+\}`)

	envvars := map[string]string{}
	for _, m := range re.FindAllString(s, -1) {
		mre := regexp.MustCompile(`[${}]`)
		mtrimmed := mre.ReplaceAllString(m, "")
		// fmt.Printf("%s:\t%s\n", mtrimmed, os.Getenv(mtrimmed))
		envvars[m] = os.Getenv(mtrimmed)
	}

	for k, v := range envvars {
		s = strings.ReplaceAll(s, k, v)
	}
	return T(s)
}

func ConfigToYAML(cfg Config) (string, error) {
	b, err := yaml.Marshal(cfg)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
