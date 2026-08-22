package winlog

import (
	"flag"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/elastic/go-ucfg"
	"github.com/stretchr/testify/assert"
)

var (
	update = flag.Bool("update", false, "update golden files")
)

func readGoldenFile(t *testing.T, filename string, expected []byte, update bool) []byte {
	t.Helper()
	goldenPath := filepath.Join("testdata", filename)

	if update {
		f, err := os.Create(goldenPath)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		if _, err := f.Write(expected); err != nil {
			t.Fatal(err)
		}

		return expected
	}

	data, err := ioutil.ReadFile(goldenPath)
	if err != nil {
		t.Fatal(err)
	}

	return data
}

// Update golden files by running:
//
//	go test ./pkg/generator/winlog -update
func TestGenerator_Next(t *testing.T) {
	tests := map[string]struct {
		config       map[string]interface{}
		expectedFile string
	}{
		"event4624": {
			config:       map[string]interface{}{"event_id": event4624},
			expectedFile: "event4624.xml",
		},
		"event4634": {
			config:       map[string]interface{}{"event_id": event4634},
			expectedFile: "event4634.xml",
		},
		"event4723": {
			config:       map[string]interface{}{"event_id": event4723},
			expectedFile: "event4723.xml",
		},
		"event4741": {
			config:       map[string]interface{}{"event_id": event4741},
			expectedFile: "event4741.xml",
		},
		"event4743": {
			config:       map[string]interface{}{"event_id": event4743},
			expectedFile: "event4743.xml",
		},
		"event4768": {
			config:       map[string]interface{}{"event_id": event4768},
			expectedFile: "event4768.xml",
		},
	}

	testTime, err := time.Parse(time.RFC3339, "1970-01-02T03:04:05+07:00")
	assert.NoError(t, err)

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			rand.Seed(1)
			// Reset package state. This test MUST NOT run in parallel.
			serviceSIDMap = map[string]string{}
			userSIDMap = map[string]string{}

			g, err := New(ucfg.MustNewFrom(tc.config), 0)
			assert.NoError(t, err)

			g.(*Generator).staticTime = &testTime

			got, err := g.Next()
			assert.NoError(t, err)

			expected := readGoldenFile(t, tc.expectedFile, got, *update)

			assert.Equal(t, string(expected), string(got))
		})
	}
}

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}
