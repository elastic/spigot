package clf

import (
	"testing"

	"github.com/elastic/go-ucfg"
	"github.com/stretchr/testify/assert"
)

func TestConfigs(t *testing.T) {
	tests := map[string]struct {
		config      map[string]interface{}
		hasError    bool
		errorString string
	}{
		"Valid Type": {
			config:      map[string]interface{}{"type": Name},
			hasError:    false,
			errorString: "",
		},
		"Invalid Type": {
			config:      map[string]interface{}{"type": "Bob"},
			hasError:    true,
			errorString: "'Bob' is not a valid value for 'type' expected 'clf' accessing config",
		},
		"No Type": {
			config:      map[string]interface{}{"type": ""},
			hasError:    true,
			errorString: "string value is not set accessing 'type'",
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cfg, err := ucfg.NewFrom(tc.config)
			assert.NoError(t, err)

			_, err = New(cfg, 0)
			if tc.hasError {
				assert.Error(t, err)
				assert.Equal(t, err.Error(), tc.errorString)
			}
			if !tc.hasError {
				assert.NoError(t, err)
			}
		})
	}
}
