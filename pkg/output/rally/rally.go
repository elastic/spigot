package rally

import (
	"encoding/json"
	"io"
	"os"

	"github.com/elastic/go-ucfg"
	"github.com/leehinman/spigot/pkg/output"
)

const OutputName = "rally"

type RallyOutput struct {
	pWriteCloser io.WriteCloser
}

type entry struct {
	Message string `json:"message"`
}

func init() {
	output.Register(OutputName, New)
}

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
	return &RallyOutput{pWriteCloser: pOsFile}, nil
}

func (r *RallyOutput) Write(b []byte) (int, error) {
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

func (r *RallyOutput) Close() error {
	return r.pWriteCloser.Close()
}
