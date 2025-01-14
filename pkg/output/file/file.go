// Package file implements the output of logs to a file.
//
// Configuration file supports either writing to a file or a directory
// with random names.  delimiter is required.  This is the string to
// write between log entries.  Normally a new line "\n"
//
//	output:
//	  type: file
//	  filename: "/var/tmp/rally.ndjson"
//	  delimiter: "/n"
//
// or
//
//	output:
//	  type: file
//	  directory: "/var/tmp"
//	  pattern: "rally_*"
//	  delimiter: "\r\n"
//
// directory and pattern are used in os.CreateTemp call
package file

import (
	"os"
	"path/filepath"

	"github.com/elastic/go-ucfg"
	"github.com/elastic/spigot/pkg/output"
)

// OutputName is the name of the output in the configuration file and registry
const Name = "file"

// Output stores pointer to an io.WriteCloser.  This is where the log
// entries will be written.  It also stores the delimiter that will be
// added between log entries.
type Output struct {
	name         string
	delimiter    string
	pWriteCloser *os.File
	directory    string
	pattern      string
}

func init() {
	output.Register(Name, New)
}

// New is the Factory for creating a new file output.  Calling this
// results in a file handle being opened to write the log data to.
func New(cfg *ucfg.Config) (output.Output, error) {
	var pOsFile *os.File
	var err error
	var name string

	c := defaultConfig()
	if err = cfg.Unpack(&c); err != nil {
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
		delimiter:    c.Delimiter,
		directory:    c.Directory,
		pattern:      c.Pattern,
		name:         name,
	}
	return &out, nil
}

// Name returns the name of the output
func (o *Output) Name() string {
	return o.name
}

// Destination return the path to the file being written to
func (o *Output) Destination() string {
	if o.pWriteCloser != nil {
		return o.pWriteCloser.Name()
	}
	return ""
}

// Write writes the log entry to the file handle that is opened with
// new and appends the delimiter.
func (o *Output) Write(b []byte) (n int, err error) {
	j, err := o.pWriteCloser.Write(b)
	if err != nil {
		return j, err
	}
	k, err := o.pWriteCloser.Write([]byte(o.delimiter))
	return j + k, err
}

// Close closes the io.WriteCloser.  Writes after this will fail.
func (o *Output) Close() error {
	return o.pWriteCloser.Close()
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
