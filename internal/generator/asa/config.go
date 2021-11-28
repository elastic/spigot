package asa

type config struct {
	Type             string `config:"type"`
	Enabled          bool   `config:"enabled"`
	IncludeTimestamp bool   `config:"include_timestamp"`
}

func defaultConfig() config {
	return config{
		Type: "asa",
	}
}
