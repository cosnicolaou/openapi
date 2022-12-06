// Copyright 2022 Cosmos Nicolaou. All rights reserved.
// Use of this source code is governed by the Apache-2.0
// license that can be found in the LICENSE file.

package transforms

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the loaded transformer configuration.
type Config struct {
	Transforms []yaml.Node `yaml:"transforms"`
	Names      []string
}

func (c Config) Configure(t T) error {
	if t == nil {
		return fmt.Errorf("transformer not specified")
	}
	for i, n := range c.Names {
		if n == t.Name() {
			return t.Configure(c.Transforms[i])
		}
	}
	return fmt.Errorf("transformer is not configured: %v", t.Name())
}

// LoadConfigFile loads the transform configuration from the
// specified YAML file.
func LoadConfigFile(filename string) (Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}
	return ParseConfig(data)
}

// ParseConfig parses the supplied YAML data to create an instance
// of Config.
func ParseConfig(data []byte) (Config, error) {
	cfg := Config{}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	cfg.Names = make([]string, len(cfg.Transforms))
	for i, c := range cfg.Transforms {
		n, err := nameFromNode(c)
		if err != nil {
			return Config{}, err
		}
		cfg.Names[i] = n
	}
	return cfg, nil
}

func nameFromNode(node yaml.Node) (string, error) {
	tmp := struct {
		Name string
	}{}
	err := node.Decode(&tmp)
	return tmp.Name, err
}
