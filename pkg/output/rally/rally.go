// Package rally outputs logs to ndjson suitable for use by https://github.com/elastic/rally
//
// Configuration file supports either writing to a file or a directory with random names.
//
//	output:
//	  type: rally
//	  filename: "/var/tmp/rally.ndjson"
//
// or
//
//	output:
//	  type: rally
//	  directory: "/var/tmp"
//	  pattern: "rally_*"
//
// directory and pattern are used in os.CreateTemp call
package rally

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/elastic/go-ucfg"
	"github.com/elastic/spigot/pkg/output"
)

// Name is the name of the output in the configuration file and registry
const Name = "rally"

// RallyOutput stores pointer to an io.WriteCloser.  This is where the
// log entries will be written.
type Output struct {
	name         string
	pWriteCloser *os.File
	directory    string
	pattern      string
}

type entry struct {
	Message string `json:"message"`
}

func init() {
	output.Register(Name, New)
}

// New is the Factory for creating a new rally output.  Calling this
// results in a file handle being opened to write the data to.
func New(cfg *ucfg.Config) (output.Output, error) {
	var pOsFile *os.File
	var err error
	var name string

	c := defaultConfig()
	if err := cfg.Unpack(&c); err != nil {
		return nil, err
	}
	if c.Directory != "" && c.Pattern != "" {
		pOsFile, err = os.CreateTemp(c.Directory, c.Pattern)
		if err != nil {
			return nil, err
		}
		name = Name + "_" + filepath.Join(c.Directory, c.Pattern)
	}
	if c.Filename != "" {
		pOsFile, err = os.Create(c.Filename)
		if err != nil {
			return nil, err
		}
		name = Name + "_" + pOsFile.Name()
	}
	out := Output{
		pWriteCloser: pOsFile,
		directory:    c.Directory,
		pattern:      c.Pattern,
		name:         name,
	}
	return &out, nil
}

// Name return the name of the output
func (r *Output) Name() string {
	return r.name
}

// Destination return the path to the file being written to
func (o *Output) Destination() string {
	if o.pWriteCloser != nil {
		return o.pWriteCloser.Name()
	}
	return ""
}

// Write formats the log for rally and writes the data to the file
// handle that was opened with New
func (r *Output) Write(b []byte) (int, error) {
	e := &entry{
		Message: string(b),
	}
	jsonBytes, err := json.Marshal(e)
	if err != nil {
		return 0, err
	}
	n, err := r.pWriteCloser.Write(jsonBytes)
	if err != nil {
		return n, err
	}
	k, err := r.pWriteCloser.Write([]byte("\n"))
	return n + k, err
}

// Close closes the io.WriteCloser.  Writes after this will fail.
func (r *Output) Close() error {
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
	return nil
}
