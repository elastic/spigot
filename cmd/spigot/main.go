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
	Workers       int            `config:"workers"`
	Records       int            `config:"records"`
	Interval      time.Duration  `config:"interval"`
	FileConfig    file.Config    `config:"output_file"`
	S3Config      s3.Config      `config:"output_s3"`
	SyslogConfig  syslog.Config  `config:"output_syslog"`
	AsaConfig     asa.Config     `config:"generator_asa"`
	VpcflowConfig vpcflow.Config `config:"generator_vpcflow"`
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

func outputFromConfig(conf Config) (out output.Output, err error) {
	if conf.FileConfig.Enabled {
		return file.New(conf.FileConfig)
	}
	if conf.S3Config.Enabled {
		return s3.New(conf.S3Config)
	}
	if conf.SyslogConfig.Enabled {
		return syslog.New(conf.SyslogConfig)
	}
	fmt.Printf("%+v", conf)
	return nil, fmt.Errorf("No output configured")
}

func generatorFromConfig(conf Config) (gen generator.Generator, err error) {
	if conf.VpcflowConfig.Enabled {
		return vpcflow.New()
	}
	if conf.AsaConfig.Enabled {
		return asa.New(conf.AsaConfig)
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
			out, err := outputFromConfig(spigotConfig)
			if err != nil {
				panic(err)
			}
			gen, err := generatorFromConfig(spigotConfig)
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
