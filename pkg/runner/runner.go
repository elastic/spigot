package runner

import (
	"time"

	"github.com/elastic/go-ucfg"
	"github.com/leehinman/spigot/pkg/generator"
	_ "github.com/leehinman/spigot/pkg/include"
	"github.com/leehinman/spigot/pkg/output"
)

type Runner struct {
	config    Config
	generator generator.Generator
	output    output.Output
}

type Config struct {
	Generator *ucfg.Config  `config:"generator" validate:"required"`
	Output    *ucfg.Config  `config:"output" validate:"required"`
	Interval  time.Duration `config:"interval"`
	Records   int           `config:"records"`
}

func defaultConfig() Config {
	c := Config{
		Records: 1024,
	}
	return c
}

func New(cfg *ucfg.Config) (Runner, error) {
	r := Runner{}
	c := defaultConfig()
	err := cfg.Unpack(&c)
	if err != nil {
		return r, err
	}

	r.config = c

	o, err := output.New(c.Output)

	if err != nil {
		return r, err
	}

	r.output = o

	g, err := generator.New(c.Generator)
	if err != nil {
		return r, err
	}
	r.generator = g

	return r, nil
}

func (r *Runner) Execute() error {
	var ticker *time.Ticker = nil
	if r.config.Interval > 0 {
		ticker = time.NewTicker(r.config.Interval)
	}

	for ; true; <-ticker.C {
		for i := 0; i < r.config.Records; i++ {
			b, err := r.generator.Next()
			if err != nil {
				return err
			}
			_, err = r.output.Write(b)
			if err != nil {
				return err
			}
		}
		if r.config.Interval == 0 {
			break
		}
	}
	r.output.Close()
	return nil
}
