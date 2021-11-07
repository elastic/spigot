package file

import (
	"io"
	"os"
)

type FileOutput struct {
	delimiter string
	pWC       io.WriteCloser
}

func New(c Config) (f *FileOutput, err error) {
	pWC, err := os.CreateTemp(c.Directory, c.Pattern)
	if err != nil {
		return nil, err
	}
	f = &FileOutput{
		delimiter: c.Delimiter,
		pWC:       pWC,
	}
	return f, nil
}

func (f *FileOutput) Write(b []byte) (n int, err error) {
	j, err := f.pWC.Write(b)
	if err != nil {
		return j, err
	}
	k, err := f.pWC.Write([]byte(f.delimiter))
	return j + k, err
}

func (f *FileOutput) Close() error {
	return f.pWC.Close()
}
