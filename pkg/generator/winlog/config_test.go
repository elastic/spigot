package winlog

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
			errorString: "'Bob' is not a valid value for 'type' expected 'winlog' accessing config",
		},
		"Invalid ID": {
			config:      map[string]interface{}{"event_id": 1},
			hasError:    true,
			errorString: "'1' is not a valid value for 'event_id' expected one of [4624 4634 4723 4741 4743 4768] accessing config",
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
				assert.Equal(t, tc.errorString, err.Error())
			}
			if !tc.hasError {
				assert.NoError(t, err)
			}
		})
	}
}
