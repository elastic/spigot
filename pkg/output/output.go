// Package output provides basic interface for log output.  It's
// primary job is to wrap specific implementations of log outputs and
// provide a shared public interface.
package output

import (
	"github.com/elastic/go-ucfg"
)

// Output is the inteface that wraps the Write and Close methods.
type Output interface {
	Write(p []byte) (n int, err error)
	Close() error
}

type config struct {
	Type string `config:"type" validate:"required"`
}

// New creates a new instance of the output that is specified by the
// "type" in the ucfg.Config that is passed in.  If no matching output
// is found for that type than an error is returned.
func New(cfg *ucfg.Config) (Output, error) {
	c := config{}
	err := cfg.Unpack(&c)
	if err != nil {
		return nil, err
	}
	factory, err := GetFactory(c.Type)
	if err != nil {
		return nil, err
	}
	return factory(cfg)
}
