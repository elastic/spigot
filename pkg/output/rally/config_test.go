package rally

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigs(t *testing.T) {
	tests := map[string]struct {
		c           config
		hasError    bool
		errorString string
	}{
		"Valid Filename": {
			c:           config{Type: OutputName, Filename: "output.log"},
			hasError:    false,
			errorString: "",
		},
		"Valid Dir and Pattern": {
			c:           config{Type: OutputName, Directory: "/var/tmp", Pattern: "output_"},
			hasError:    false,
			errorString: "",
		},
		"Wrong type": {
			c:           config{Type: "malory", Filename: "output.log"},
			hasError:    true,
			errorString: "malory is not a valid type for rally",
		},
		"Dir and filename set": {
			c:           config{Type: OutputName, Directory: "/var/tmp", Filename: "output.log"},
			hasError:    true,
			errorString: "if filename is set, directory and pattern must not be",
		},
		"Pattern and filename set": {
			c:           config{Type: OutputName, Pattern: "output_", Filename: "output.log"},
			hasError:    true,
			errorString: "if filename is set, directory and pattern must not be",
		},
		"Only Directory set": {
			c:           config{Type: OutputName, Directory: "/var/tmp"},
			hasError:    true,
			errorString: "directory and pattern must both be set",
		},
		"Only Pattern set": {
			c:           config{Type: OutputName, Pattern: "output_"},
			hasError:    true,
			errorString: "directory and pattern must both be set",
		},
		"Dir, Pattern and Filename set": {
			c:           config{Type: OutputName, Pattern: "output_", Filename: "output.log", Directory: "/var/tmp"},
			hasError:    true,
			errorString: "if filename is set, directory and pattern must not be",
		},
		"Only type set": {
			c:           config{Type: OutputName},
			hasError:    true,
			errorString: "you must specify filename or directory and pattern",
		},
	}
	for name, tc := range tests {
		err := tc.c.Validate()
		if tc.hasError {
			assert.NotNil(t, err, name)
			assert.Equal(t, err.Error(), tc.errorString, name)
		}
		if !tc.hasError {
			assert.Nil(t, err, name)
		}
	}
}
