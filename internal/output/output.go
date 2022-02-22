package output

import (
	"fmt"

	"github.com/elastic/go-ucfg"
	"github.com/leehinman/spigot/internal/output/file"
	"github.com/leehinman/spigot/internal/output/s3"
	"github.com/leehinman/spigot/internal/output/syslog"
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
	switch c.Type {
	case "file":
		return file.New(cfg)
	case "s3":
		return s3.New(cfg)
	case "syslog":
		return syslog.New(cfg)
	default:
		return nil, fmt.Errorf("Unknown output: %s", c.Type)
	}
}
