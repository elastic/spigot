package clf

import (
	"math/rand"
	"testing"
	"time"

	"github.com/elastic/go-ucfg"
	"github.com/stretchr/testify/assert"
)

func TestGenerator_Next(t *testing.T) {
	tests := map[string]struct {
		config   map[string]interface{}
		expected string
	}{
		"common": {
			config:   map[string]interface{}{"combined": false},
			expected: `66.4.203.154 - - [02/Jan/1970:03:04:05 +0700] "GET /random-47.html HTTP/2" 200 1318`,
		},
		"combined": {
			config:   map[string]interface{}{"combined": true},
			expected: `66.4.203.154 - - [02/Jan/1970:03:04:05 +0700] "GET /random-47.html HTTP/2" 200 1318 - "Mozilla/5.0 (Macintosh; Intel Mac OS X 12.3; rv:98.0) Gecko/20100101 Firefox/98.0"`,
		},
	}

	testTime, err := time.Parse(time.RFC3339, "1970-01-02T03:04:05+07:00")
	assert.NoError(t, err)

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			rand.Seed(1)

			g, err := New(ucfg.MustNewFrom(tc.config), 0)
			assert.NoError(t, err)

			g.(*Generator).staticTime = &testTime

			got, err := g.Next()
			assert.NoError(t, err)

			assert.Equal(t, tc.expected, string(got))
		})
	}
}
