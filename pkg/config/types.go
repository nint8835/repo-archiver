package config

type Config struct {
	Accounts map[string]Account `yaml:"accounts"`
}

type Account struct {
	Name      string `yaml:"name"`
	IsArchive bool   `yaml:"is_archive"`
}
