package generator

import (
	"fmt"

	"github.com/elastic/go-ucfg"
	"github.com/leehinman/spigot/internal/generator/aws/vpcflow"
	"github.com/leehinman/spigot/internal/generator/cisco/asa"
	"github.com/leehinman/spigot/internal/generator/fortinet/firewall"
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
	switch c.Type {
	case "aws:vpcflow":
		return vpcflow.New(cfg)
	case "cisco:asa":
		return asa.New(cfg)
	case "fortinet:firewall":
		return firewall.New(cfg)
	default:
		return nil, fmt.Errorf("Unknown generator: %s", c.Type)
	}
}
