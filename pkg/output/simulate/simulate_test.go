package simulate

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
		want  string
	}{
		"One Event": {
			input: []string{"a"},
			want:  "{\"events\":[{\"message\":\"a\"}]}\n",
		},
		"Two Events": {
			input: []string{"a", "b"},
			want:  "{\"events\":[{\"message\":\"a\"},{\"message\":\"b\"}]}\n",
		},
	}
	for name, tc := range tests {
		var buf bytes.Buffer
		var wc = &myWriteCloser{&buf}

		events := []event{}
		r := &Output{
			events:       events,
			pWriteCloser: wc,
		}
		for _, line := range tc.input {
			_, err := r.Write([]byte(line))
			assert.Nil(t, err)
		}
		r.Close()
		assert.Equal(t, []byte(tc.want), buf.Bytes(), name)
	}
}
