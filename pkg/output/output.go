package output

import (
	"github.com/elastic/go-ucfg"
)

type Output interface {
	Write(p []byte) (n int, err error)
	Close() error
}

type config struct {
	Type string `config:"type" validate:"required"`
}

func New(cfg *ucfg.Config) (Output, error) {
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
