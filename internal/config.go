package utils

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const configFileName string = "quack_config.yaml"

// string and arrays types for yaml config
type StringVal string
type StringList []string

func (s StringVal) String() string {
	return string(s)
}

func (s *StringVal) Set(value string) error {
	*s = StringVal(value)
	return nil
}

func (s StringList) String() string {
	return strings.Join(s, ",")
}

func (s *StringList) Set(value string) error {
	*s = strings.Split(value, ",")
	return nil
}

// main config struct
type ConfigYaml struct {
	Version  float32 `yaml:"version"`
	Database struct {
		Uri     StringVal  `yaml:"uri"`
		Name    StringVal  `yaml:"name"`
		Exclude StringList `yaml:"exclude"`
		Type    string
	} `yaml:"database"`
	Models struct {
		Path    StringVal  `yaml:"path"`
		Exclude StringList `yaml:"exclude"`
	} `yaml:"models"`
	Migrations struct {
		Path StringVal `yaml:"path"`
	} `yaml:"migrations"`
}

func (conf *ConfigYaml) ReadConfig() error {
	cfile, err := findConfigFile()
	if err != nil {
		return err
	}
	read_err := yaml.Unmarshal(cfile, conf)
	if read_err != nil {
		return nil
	}
	return err
}

func findConfigFile() ([]byte, error) {
	cfile, err := os.ReadFile(fmt.Sprintf("./%s", configFileName))
	if err != nil {
		// try parent directory
		cfile, err = os.ReadFile(fmt.Sprintf("../%s", configFileName))
		if err != nil {
			return nil, err
		}
		fmt.Println("Found config in parent directory")
	}
	fmt.Println("Found config file in current directory")
	return cfile, nil
}
