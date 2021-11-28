package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/elastic/go-ucfg"
	"github.com/elastic/go-ucfg/yaml"
	"github.com/leehinman/spigot/internal/generator"
	"github.com/leehinman/spigot/internal/generator/asa"
	"github.com/leehinman/spigot/internal/generator/vpcflow"
	"github.com/leehinman/spigot/internal/output"
	"github.com/leehinman/spigot/internal/output/file"
	"github.com/leehinman/spigot/internal/output/s3"
	"github.com/leehinman/spigot/internal/output/syslog"
)

type Config struct {
	Workers    int            `config:"workers"`
	Records    int            `config:"records"`
	Interval   time.Duration  `config:"interval"`
	Outputs    []*ucfg.Config `config:"outputs" validate:"required"`
	Generators []*ucfg.Config `config:"generators" validate: "required"`
}

var (
	defaultConfig = Config{
		Workers:  1,
		Records:  1024,
		Interval: 5 * time.Second,
	}
)

func run(gen generator.Generator, out output.Output, lines int) {
	for i := 0; i < lines; i++ {
		b, err := gen.Next()
		if err != nil {
			panic(err)
		}
		_, err = out.Write(b)
		if err != nil {
			panic(err)
		}
	}
	out.Close()
}

func outputFromConfig(cfgs []*ucfg.Config) (out output.Output, err error) {
	for _, cfg := range cfgs {
		c := output.OutputConfig{}
		if err := cfg.Unpack(&c); err != nil {
			return nil, err
		}
		if !c.Enabled {
			continue
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
	return nil, fmt.Errorf("No output configured")
}

func generatorFromConfig(cfgs []*ucfg.Config) (gen generator.Generator, err error) {
	for _, cfg := range cfgs {
		c := generator.GeneratorConfig{}
		if err := cfg.Unpack(&c); err != nil {
			return nil, err
		}
		if !c.Enabled {
			continue
		}
		switch c.Type {
		case "vpcflow":
			return vpcflow.New(cfg)
		case "asa":
			return asa.New(cfg)
		default:
			return nil, fmt.Errorf("Unknown generator: %s", c.Type)
		}
	}
	return nil, fmt.Errorf("No generator configured")
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

	ticker := time.NewTicker(spigotConfig.Interval)
	for ; true; <-ticker.C {
		var wg sync.WaitGroup

		for i := 0; i < spigotConfig.Workers; i++ {
			out, err := outputFromConfig(spigotConfig.Outputs)
			if err != nil {
				panic(err)
			}
			gen, err := generatorFromConfig(spigotConfig.Generators)
			if err != nil {
				panic(err)
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				run(gen, out, spigotConfig.Records)
			}()
		}
		wg.Wait()
	}
}
