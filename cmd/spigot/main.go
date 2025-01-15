package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/elastic/go-ucfg"
	"github.com/elastic/go-ucfg/yaml"
	"github.com/elastic/spigot/pkg/runner"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

type Config struct {
	Runners []*ucfg.Config `config:"runners" validate:"required"`
}

type Result struct {
	Done         bool
	Error        error
	Metrics      metricdata.ResourceMetrics
	Destinations []string
}

func execute_runner(cfg *ucfg.Config, results chan Result, done chan struct{}) {
	r, err := runner.New(cfg)
	if err != nil {
		results <- Result{Error: err}
		return
	}
	m, dst, err := r.Execute(done)
	if err != nil {
		results <- Result{Error: err}
		return
	}
	results <- Result{Done: true, Metrics: m, Destinations: dst}
	return
}

func metricsToJson(rm metricdata.ResourceMetrics, dst []string) ([]byte, error) {
	out := map[string]interface{}{}
	out["destinations"] = dst
	for _, sm := range rm.ScopeMetrics {
		out["name"] = sm.Scope.Name
		for _, m := range sm.Metrics {
			var sum int64
			for _, dp := range m.Data.(metricdata.Sum[int64]).DataPoints {
				sum = sum + dp.Value
			}
			out[m.Name] = sum
		}
	}
	return json.Marshal(out)
}

func main() {
	var cfgFile string
	var randomize bool
	var printMetrics bool

	flag.StringVar(&cfgFile, "c", "./spigot.yml", "path to configuration file")
	flag.BoolVar(&randomize, "r", false, "seed random number generator with current time")
	flag.BoolVar(&printMetrics, "m", false, "print metrics at end of run")
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
	doneCh := make(chan struct{})
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	go func() {
		for range signalCh {
			close(doneCh)
		}
	}()

	for _, rCfg := range c.Runners {
		rCfg := rCfg
		go func() {
			execute_runner(rCfg, resultCh, doneCh)
		}()
	}

	for i := 0; i < len(c.Runners); i++ {
		r := <-resultCh
		if !r.Done {
			panic(r.Error)
		}
		if printMetrics {

			j, err := metricsToJson(r.Metrics, r.Destinations)
			if err != nil {
				panic(fmt.Errorf("error converting metrics to json: %w", err))
			}
			fmt.Fprintf(os.Stdout, "%s\n", string(j))
		}
	}
}
