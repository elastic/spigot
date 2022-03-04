package file

type config struct {
	Type      string `config:"type"`
	Directory string `config:"directory"`
	Pattern   string `config:"pattern"`
	Delimiter string `config:"delimiter"`
}

func defaultConfig() config {
	return config{
		Type:      "file",
		Delimiter: "\n",
	}
}
