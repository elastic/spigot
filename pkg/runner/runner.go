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
	"time"

	"github.com/elastic/go-ucfg"
	"github.com/elastic/spigot/pkg/generator"
	_ "github.com/elastic/spigot/pkg/include"
	"github.com/elastic/spigot/pkg/output"
)

// Runner holds the config, output and generator.
type Runner struct {
	config    config
	generator generator.Generator
	output    output.Output
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
		return r, err
	}

	r.config = c

	o, err := output.New(c.Output)
	if err != nil {
		return r, err
	}

	r.output = o

	g, err := generator.New(c.Generator, int64(c.Records))
	if err != nil {
		return r, err
	}
	r.generator = g

	return r, nil
}

// Execute runs the runner
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
		if err := r.output.NewInterval(); err != nil {
			return err
		}
	}
	return r.output.Close()
}
