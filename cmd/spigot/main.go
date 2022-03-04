package main

import (
	"flag"
	"math/rand"
	"time"

	"github.com/elastic/go-ucfg"
	"github.com/elastic/go-ucfg/yaml"
	"github.com/leehinman/spigot/pkg/runner"
)

type Config struct {
	Runners []*ucfg.Config `config:"runners" validate:"required"`
}

type Result struct {
	Done  bool
	Error error
}

func execute_runner(cfg *ucfg.Config, results chan Result) {
	r, err := runner.New(cfg)
	if err != nil {
		results <- Result{Error: err}
		return
	}
	err = r.Execute()
	if err != nil {
		results <- Result{Error: err}
		return
	}
	results <- Result{Done: true}
	return
}

func main() {
	var cfgFile string
	var randomize bool

	flag.StringVar(&cfgFile, "c", "./spigot.yml", "path to configuration file")
	flag.BoolVar(&randomize, "r", false, "seed random number generator with current time")
	flag.Parse()

	c := Config{}
	cfg, err := yaml.NewConfigWithFile(cfgFile, ucfg.PathSep("."))
	if err != nil {
		panic(err)
	}
	err = cfg.Unpack(&c)
	if err != nil {
		panic(err)
	}

	if randomize {
		rand.Seed(time.Now().UnixNano())
	}

	resultCh := make(chan Result)

	for _, rCfg := range c.Runners {
		rCfg := rCfg
		go func() {
			execute_runner(rCfg, resultCh)
		}()
	}

	for i := 0; i < len(c.Runners); i++ {
		r := <-resultCh
		if !r.Done {
			panic(r.Error)
		}
	}
}
