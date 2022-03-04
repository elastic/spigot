package file

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type myWriteCloser struct {
	io.Writer
}

func (*myWriteCloser) Close() error {
	return nil
}

func TestWrite(t *testing.T) {
	tests := map[string]struct {
		input []string
		delim string
		want  string
	}{
		"OneLine,NewLine": {
			input: []string{"a"},
			delim: "\n",
			want:  "a\n",
		},
		"OneLine,Tab": {
			input: []string{"a"},
			delim: "\t",
			want:  "a\t",
		},
		"TwoLine,NewLine": {
			input: []string{"a", "b"},
			delim: "\n",
			want:  "a\nb\n",
		},
		"TwoLine,Tab": {
			input: []string{"a", "b"},
			delim: "\t",
			want:  "a\tb\t",
		},
	}
	for name, tc := range tests {
		var buf bytes.Buffer
		var wc = &myWriteCloser{&buf}

		f := &FileOutput{
			delimiter: tc.delim,
			pWC:       wc,
		}
		for _, line := range tc.input {
			_, err := f.Write([]byte(line))
			assert.Nil(t, err)
		}
		assert.Equal(t, []byte(tc.want), buf.Bytes(), name)
	}
}
