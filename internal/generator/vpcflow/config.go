package vpcflow

type config struct {
	Type    string `config:"type"`
	Enabled bool   `config:"enabled"`
}

func defaultConfig() config {
	return config{
		Type: "vpcflow",
	}
}
