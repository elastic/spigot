// Package runner provides the glue to link a generator to an output and to execute.
//
//	Configuration.
//
//	"generator" and "output" are required, and are the configs of the
//	specific types.
//
//	"records" is optional, default is 1024.  This is the number of log
//	records to write per interval.
//
//	"interval" is optional and is a go duration.  If no interval is
//	given then the runner is executed once.  If an interval is given
//	then at each interval the runner is executed.
//
//	Example:
//
//	  generator:
//	    type: "aws:vpcflow"
//	  output:
//	    type: file
//	    directory: "/var/tmp"
//	    pattern: "spigot_asa_*.log"
//	    delimiter: "\n"
//	  interval: 5s
//	  records: 2
//
//	This would write 2 vpcflow log entries to a file in the
//	/var/tmp/spigot_asa_<random>.log file every 5 seconds.
package runner

import (
	"context"
	"fmt"
	"time"

	"github.com/elastic/go-ucfg"
	"github.com/elastic/spigot/pkg/generator"
	_ "github.com/elastic/spigot/pkg/include"
	"github.com/elastic/spigot/pkg/output"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// Runner holds the config, output and generator.
type Runner struct {
	config              config
	generator           generator.Generator
	output              output.Output
	metricReader        sdkmetric.Reader
	metricMeterProvider *sdkmetric.MeterProvider
	name                string
	destinations        []string
}

type Metrics struct {
	recordsCounter  metric.Int64Counter
	recordsBytes    metric.Int64Counter
	intervalCounter metric.Int64Counter
}

type config struct {
	Generator *ucfg.Config  `config:"generator" validate:"required"`
	Output    *ucfg.Config  `config:"output" validate:"required"`
	Interval  time.Duration `config:"interval"`
	Records   int           `config:"records"`
}

func defaultConfig() config {
	c := config{
		Records: 1024,
	}
	return c
}

// New is Factory for creating a new runner
func New(cfg *ucfg.Config) (Runner, error) {
	r := Runner{}
	c := defaultConfig()
	err := cfg.Unpack(&c)
	if err != nil {
		return r, fmt.Errorf("error unpacking config: %w", err)
	}

	r.config = c
	o, err := output.New(c.Output)
	if err != nil {
		return r, fmt.Errorf("error creating output: %w", err)
	}

	r.output = o
	r.name = r.output.Name()
	r.destinations = appendIfMissing(r.destinations, r.output.Destination())

	g, err := generator.New(c.Generator)
	if err != nil {
		return r, err
	}
	r.generator = g

	r.metricReader = sdkmetric.NewManualReader()
	r.metricMeterProvider = sdkmetric.NewMeterProvider(sdkmetric.WithReader(r.metricReader))
	return r, nil
}

// Execute runs the runner
func (r *Runner) Execute(done chan struct{}) (metricdata.ResourceMetrics, []string, error) {
	m := &Metrics{}
	ctx := context.Background()
	var ticker *time.Ticker = nil
	if r.config.Interval > 0 {
		ticker = time.NewTicker(r.config.Interval)
	}

	meter := r.metricMeterProvider.Meter(r.name)
	recordsCounter, err := meter.Int64Counter("records")
	if err != nil {
		return metricdata.ResourceMetrics{}, r.destinations, fmt.Errorf("error creating runner records counter: %w", err)
	}
	m.recordsCounter = recordsCounter

	recordsBytes, err := meter.Int64Counter("bytes")
	if err != nil {
		return metricdata.ResourceMetrics{}, r.destinations, fmt.Errorf("error creating runner bytes counter: %w", err)
	}
	m.recordsBytes = recordsBytes

	intervalCounter, err := meter.Int64Counter("intervals")
	if err != nil {
		return metricdata.ResourceMetrics{}, r.destinations, fmt.Errorf("error createing runner intervals counter: %w", err)
	}
	m.intervalCounter = intervalCounter

	if r.config.Interval == 0 {
		for i := 0; i < r.config.Records; i++ {
			b, err := r.generator.Next()
			if err != nil {
				return metricdata.ResourceMetrics{}, r.destinations, fmt.Errorf("error calling generator Next: %w", err)
			}
			n, err := r.output.Write(b)
			if err != nil {
				return metricdata.ResourceMetrics{}, r.destinations, fmt.Errorf("error calling output Write: %w", err)
			}
			m.recordsCounter.Add(ctx, 1)
			m.recordsBytes.Add(ctx, int64(n))
		}
	} else {
		for {
			select {
			case <-ticker.C:
				for i := 0; i < r.config.Records; i++ {
					b, err := r.generator.Next()
					if err != nil {
						return metricdata.ResourceMetrics{}, r.destinations, fmt.Errorf("error calling generator Next: %w", err)
					}
					n, err := r.output.Write(b)
					if err != nil {
						return metricdata.ResourceMetrics{}, r.destinations, fmt.Errorf("error calling output Write: %w", err)
					}
					m.recordsCounter.Add(ctx, 1)
					m.recordsBytes.Add(ctx, int64(n))
				}
				if err := r.output.NewInterval(); err != nil {
					return metricdata.ResourceMetrics{}, r.destinations, fmt.Errorf("error calling output NewInterval: %w", err)
				}
				r.destinations = appendIfMissing(r.destinations, r.output.Destination())
				m.intervalCounter.Add(ctx, 1)
			case <-done:
				goto printMetrics
			}
		}
	}

printMetrics:

	rm := metricdata.ResourceMetrics{}
	if err := r.metricReader.Collect(ctx, &rm); err != nil {
		panic(fmt.Errorf("error collecting runner metrics: %w", err))
	}

	return rm, r.destinations, r.output.Close()
}

func appendIfMissing(s []string, e string) []string {
	for _, element := range s {
		if element == e {
			return s
		}
	}
	return append(s, e)
}
