package s3

type Config struct {
	Enabled   bool   `config:"enabled"`
	Bucket    string `config:"bucket"`
	Region    string `config:"region"`
	Delimiter string `config:"delimiter"`
	Prefix    string `config:"prefix"`
}
