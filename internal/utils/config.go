package utils

type Config struct {
	Path           string
	Postgres_url   string
	DBname         string
	ExcludeTables  []string
	ExcludeModels  []string
	MigrationsPath string
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
