package s3

type config struct {
	Type      string `config:"type"`
	Bucket    string `config:"bucket"`
	Region    string `config:"region"`
	Delimiter string `config:"delimiter"`
	Prefix    string `config:"prefix"`
}

func defaultConfig() config {
	return config{
		Type:      "s3",
		Delimiter: "\n",
	}
}
