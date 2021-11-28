package generator

type Generator interface {
	Next() ([]byte, error)
}

type GeneratorConfig struct {
	Type    string `config:"type" validate:"required"`
	Enabled bool   `config:"enabled" validate:"required"`
}
