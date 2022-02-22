package asa

type config struct {
	Type             string `config:"type"`
	IncludeTimestamp bool   `config:"include_timestamp"`
}

func defaultConfig() config {
	return config{
		Type: "cisco:asa",
	}
}
