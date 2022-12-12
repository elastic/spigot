// Package simulate outputs logs to json suitable for use as a
// https://github.com/elastic/elastic-package test pipeline input
// events file.
//
// Configuration file supports either writing to a file or a directory
// with random names.
//
//	output:
//	  type: simulate
//	  filename: "/var/tmp/simulate.json"
//
// or
//
//	output:
//	  type: simulate
//	  directory: "/var/tmp"
//	  pattern: "simulate_*"
//
// directory and pattern are used in os.CreateTemp call
package simulate

import (
	"encoding/json"
	"io"
	"os"

	"github.com/elastic/go-ucfg"
	"github.com/leehinman/spigot/pkg/output"
)

// Name is the name of the output in the configuration file and registry
const Name = "simulate"

// Output stores pointer to an io.WriteCloser.  This is where the
// log entries will be written.
type Output struct {
	events       []event
	pWriteCloser io.WriteCloser
	directory    string
	pattern      string
}

type event struct {
	Message string `json:"message"`
}

type inputEvents struct {
	Events []event `json:"events"`
}

func init() {
	output.Register(Name, New)
}

// New is the Factory for creating a new simulate output.  Calling this
// results in a file handle being opened to write the data to.
func New(cfg *ucfg.Config) (output.Output, error) {
	var pOsFile *os.File
	var err error

	c := defaultConfig()
	if err := cfg.Unpack(&c); err != nil {
		return nil, err
	}
	if c.Directory != "" && c.Pattern != "" {
		pOsFile, err = os.CreateTemp(c.Directory, c.Pattern)
		if err != nil {
			return nil, err
		}
	}
	if c.Filename != "" {
		pOsFile, err = os.Create(c.Filename)
		if err != nil {
			return nil, err
		}
	}
	events := []event{}

	out := Output{
		pWriteCloser: pOsFile,
		events:       events,
		directory:    c.Directory,
		pattern:      c.Pattern,
	}
	return &out, nil
}

// Write formats the event and adds it to the internal slice of events.
// actual writing will happen when Close is called.
func (r *Output) Write(b []byte) (int, error) {
	e := event{
		Message: string(b),
	}
	r.events = append(r.events, e)
	return 0, nil
}

// Close Marshals the internal slice of events and writes the JSON to
// the io.WriteCloser.  Adds a newline at end of JSON data.
// Writes after this will fail.
func (r *Output) Close() error {
	inputEvents := &inputEvents{Events: r.events}

	jsonBytes, err := json.Marshal(inputEvents)
	if err != nil {
		return err
	}

	_, err = r.pWriteCloser.Write(jsonBytes)
	if err != nil {
		return err
	}

	_, err = r.pWriteCloser.Write([]byte("\n"))
	if err != nil {
		return err
	}

	return r.pWriteCloser.Close()
}

func (o *Output) NewInterval() error {
	if o.directory == "" && o.pattern == "" {
		return nil
	}
	if err := o.Close(); err != nil {
		return err
	}
	pOsFile, err := os.CreateTemp(o.directory, o.pattern)
	if err != nil {
		return err
	}
	o.pWriteCloser = pOsFile
	o.events = o.events[:0]
	return nil
}
