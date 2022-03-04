package rally

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
		"OneLine": {
			input: []string{"a"},
			want:  "{\"message\":\"a\"}\n",
		},
		"TwoLine": {
			input: []string{"a", "b"},
			want:  "{\"message\":\"a\"}\n{\"message\":\"b\"}\n",
		},
	}
	for name, tc := range tests {
		var buf bytes.Buffer
		var wc = &myWriteCloser{&buf}

		r := &RallyOutput{
			pWriteCloser: wc,
		}
		for _, line := range tc.input {
			_, err := r.Write([]byte(line))
			assert.Nil(t, err)
		}
		assert.Equal(t, []byte(tc.want), buf.Bytes(), name)
	}
}
