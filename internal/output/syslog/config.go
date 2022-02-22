package syslog

type config struct {
	Type     string `config:"type"`
	Facility string `config:"facility"`
	Severity string `config:"severity"`
	Tag      string `config:"tag"`
	Network  string `config:"network"`
	Host     string `config:"host"`
	Port     string `config:"port"`
}

func defaultConfig() config {
	return config{
		Type: "syslog",
	}
}
