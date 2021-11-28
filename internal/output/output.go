package output

type Output interface {
	Write(p []byte) (n int, err error)
	Close() error
}

type OutputConfig struct {
	Type    string `config:"type" validate:"required"`
	Enabled bool   `config:"enabled" validate:"required"`
}
