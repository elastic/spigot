package vpcflow

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
		"Valid Type": {
			c:           map[string]interface{}{"type": Name},
			hasError:    false,
			errorString: "",
		},
		"Invalid Type": {
			c:           map[string]interface{}{"type": "Bob"},
			hasError:    true,
			errorString: "'Bob' is not a valid value for 'type' expected 'aws:vpcflow' accessing config",
		},
		"No Type": {
			c:           map[string]interface{}{"type": ""},
			hasError:    true,
			errorString: "string value is not set accessing 'type'",
		},
	}
	for name, tc := range tests {
		c, err := ucfg.NewFrom(tc.c)
		assert.Nil(t, err, name)
		_, err = New(c, 0)
		if tc.hasError {
			assert.NotNil(t, err, name)
			assert.Equal(t, err.Error(), tc.errorString, name)
		}
		if !tc.hasError {
			assert.Nil(t, err, name)
		}
	}
}
