package syslog

type Config struct {
	Enabled  bool   `config:"enabled"`
	Facility string `config:"facility"`
	Severity string `config:"severity"`
	Tag      string `config:"tag"`
	Network  string `config:"network"`
	Host     string `config:"host"`
	Port     string `config:"port"`
}
