package main

import (
	"github.com/elastic/go-ucfg"
	"github.com/elastic/go-ucfg/yaml"
	"github.com/leehinman/spigot/internal/runner"
)

type MainConfig struct {
	Runners []*ucfg.Config `config:"runners" validate:"required"`
}

type Result struct {
	Done  bool
	Error error
}

var (
	defaultConfig = MainConfig{}
)

func run(cfg *ucfg.Config, results chan Result) {
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
	spigotConfig := defaultConfig
	config, err := yaml.NewConfigWithFile("./spigot.yml", ucfg.PathSep("."))
	if err != nil {
		panic(err)
	}
	err = config.Unpack(&spigotConfig)
	if err != nil {
		panic(err)
	}

	// rand.Seed(time.Now().UnixNano())

	resultCh := make(chan Result)

	for _, runner_cfg := range spigotConfig.Runners {
		runner_cfg := runner_cfg
		go func() {
			run(runner_cfg, resultCh)
		}()
	}

	for i := 0; i < len(spigotConfig.Runners); i++ {
		r := <-resultCh
		if !r.Done {
			panic(r.Error)
		}
	}
}
