package utils

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const configFileName string = "quack_config.yaml"

type Config struct {
	Path           string
	Postgres_url   string
	DBname         string
	ExcludeTables  []string
	ExcludeModels  []string
	MigrationsPath string
}

/*
	Example
project:
  files:
    - file1
    - file2
  folders:
    - folder1
    - folder2
  random1:
  random2:
    - redundant1

  version: 0.1
  database:
~   uri: "postgres://stexp:1!password!2@postgres:5432/stexp"
~   name: "stexp"
    exclude:
+     - "auth_users"
+     - "users"
+     - "goose_migrations"
  models:
~   path: "models"
    exclude:
+     - "Base"
+     - "Users"
+     - "AuthUsers"
  migrations:
~   path: "migrations"

type MyStruct struct {
    Project struct {
        Files   []string `yaml:"files"`
        Folders []string `yaml:"folders"`
    } `yaml:"project"`
}

*/

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

type ConfigYaml struct {
	Version  float32 `yaml:"version"`
	Database struct {
		Uri     StringVal  `yaml:"uri"`
		Name    StringVal  `yaml:"name"`
		Exclude StringList `yaml:"exclude"`
	} `yaml:"database"`
	Models struct {
		Path    StringVal  `yaml:"path"`
		Exclude StringList `yaml:"exclude"`
	} `yaml:"models"`
	Migrations struct {
		Path StringVal `yaml:"path"`
	} `yaml:"migrations"`
}

func findConfigFile() ([]byte, error) {
	cfile, err := os.ReadFile(fmt.Sprintf("./%s", configFileName))
	if err != nil {
		// try parent directory
		cfile, err = os.ReadFile(fmt.Sprintf("../%s", configFileName))
		fmt.Println("Found config in parent directory")
		if err != nil {
			return nil, err
		}
	}
	fmt.Println("Found config file in current directory")
	return cfile, nil
}

func (conf *ConfigYaml) ReadConfig() (*ConfigYaml, error) {
	cfile, err := findConfigFile()
	if err != nil {
		return conf, err
	}
	read_err := yaml.Unmarshal(cfile, conf)
	if read_err != nil {
		return nil, read_err
	}
	return conf, err
}

func (conf *Config) GetConfig() {
	conf.Path = "./models/"
	conf.Postgres_url = "postgres://stexp:1!password!2@postgres:5432/stexp"
	conf.DBname = "stexp"
	conf.ExcludeTables = []string{"auth_users", "users", "goose_migrations"}
	conf.ExcludeModels = []string{"Base", "Users", "AuthUsers"}
	conf.MigrationsPath = "./migrations"
}

func (conf *Config) ParseConsoleArgs() {}
