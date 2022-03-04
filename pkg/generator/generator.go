package generator

import (
	"github.com/elastic/go-ucfg"
)

type Generator interface {
	Next() ([]byte, error)
}

type config struct {
	Type string `config:"type" validate:"required"`
}

func New(cfg *ucfg.Config) (Generator, error) {
	c := config{}
	err := cfg.Unpack(&c)
	if err != nil {
		return nil, err
	}
	factory, err := GetFactory(c.Type)
	if err != nil {
		return nil, err
	}
	return factory(cfg)
}
