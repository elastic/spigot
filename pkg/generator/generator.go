// Package generator provides basic interface to log generators.  It's
// primary job is to wrap specific implementations of log generators
// and provide a shared public interface.
package generator

import (
	"strings"
	"text/template"

	"github.com/elastic/go-ucfg"
)

var (
	FunctionMap = template.FuncMap{
		"ToLower": strings.ToLower,
		"ToUpper": strings.ToUpper,
	}
)

// Generator is the interface that wraps the Next method.
//
// Next generates the next log message and returns it as an array of bytes.
type Generator interface {
	Next() ([]byte, error)
}

type config struct {
	Type string `config:"type" validate:"required"`
}

// New creates a new instance of the generator that is specified by
// the "type" in the ucfg.Config that is passed in.  If no matching
// generator is found for that type than an error is returned.
func New(cfg *ucfg.Config, numRecords int64) (Generator, error) {
	c := config{}
	err := cfg.Unpack(&c)
	if err != nil {
		return nil, err
	}
	factory, err := GetFactory(c.Type)
	if err != nil {
		return nil, err
	}
	return factory(cfg, numRecords)
}
