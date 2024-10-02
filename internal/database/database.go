package database

import "github.com/spoonboy-io/dujour/internal"

// DbConfig represents the YAML config which is provided in a '*.dbyaml' file
type DbConfig struct {
	Description string `yaml:"description"`
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	//Type        string `yaml:"type"` // only MYSQL supported ATM
	SQLQuery string `yaml:"sqlQuery"`
	Active   bool   `yaml:"active"`
}

// CheckConfig does some checks on the supplied configuration, values and validity
func CheckConfig(cfg *DbConfig) error {
	if cfg.User == "" {
		return internal.ERR_NO_USER
	}

	if cfg.Port == "" {
		return internal.ERR_NO_PORT
	}

	if cfg.Host == "" {
		return internal.ERR_NO_HOST
	}

	if cfg.Password == "" {
		return internal.ERR_NO_PASSWORD
	}

	if cfg.SQLQuery == "" {
		return internal.ERR_NO_QUERY
	}

	return nil
}
