package database

// DbConfig represents the YAML config which is provided in a '*.dbyaml' file
type DbConfig struct {
	Description string `yaml:"description"`
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	Type        string `yaml:"type"`
	SQLQuery    string `yaml:"sqlQuery"`
	Active      bool   `yaml:"active"`
}
