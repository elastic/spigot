package file

type Config struct {
	Enabled   bool   `config:"enabled"`
	Directory string `config:"directory"`
	Pattern   string `config:"pattern"`
	Delimiter string `config:"delimiter"`
}
