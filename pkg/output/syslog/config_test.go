//go:build !windows

package syslog

import (
	"testing"

	"github.com/elastic/go-ucfg"
	"github.com/stretchr/testify/assert"
)

func TestConfigs(t *testing.T) {
	tests := map[string]struct {
		c           map[string]interface{}
		hasError    bool
		errorString string
	}{
		"Invalid Type": {
			c:           map[string]interface{}{"type": "Bob", "network": "tcp", "host": "localhost", "port": "1234"},
			hasError:    true,
			errorString: "'Bob' is not a valid value for 'type' expected 'syslog' accessing config",
		},
		"No Type": {
			c:           map[string]interface{}{"type": "", "network": "tcp", "host": "localhost", "port": "1234"},
			hasError:    true,
			errorString: "string value is not set accessing 'type'",
		},
	}
	for name, tc := range tests {
		c, err := ucfg.NewFrom(tc.c)
		assert.Nil(t, err, name)
		_, err = New(c)
		if tc.hasError {
			assert.NotNil(t, err, name)
			assert.Equal(t, err.Error(), tc.errorString, name)
		}
		if !tc.hasError {
			assert.Nil(t, err, name)
		}
	}
}
