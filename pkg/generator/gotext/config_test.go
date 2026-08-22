package gotext

import (
	"testing"

	"github.com/elastic/go-ucfg"
	"github.com/stretchr/testify/assert"
)

func TestConfigs(t *testing.T) {

	var gcfg GeneratorConfig

	testCases := []struct {
		genConfig map[string]interface{}
		errString string
		expected  string
	}{
		{
			genConfig: map[string]interface{}{
				"name": "asa",
				"formats": []string{
					"%Basic-test: Deny {{.AccessGroup | ToLower}} because {{.AclId}} said so",
				},
				"fields": []map[string]interface{}{
					{
						"name": "AccessGroup",
						"type": "string",
						"choices": []string{
							"Access-Group",
						},
					},
					{
						"name": "AclId",
						"type": "string",
						"choices": []string{
							"AclId",
						},
					},
					{
						"name": "Direction",
						"type": "string",
						"choices": []string{
							"inbound",
							"outbound",
						},
					},
				},
			},
			errString: "",
			expected:  "",
		},
	}

	for _, tc := range testCases {
		c, err := ucfg.NewFrom(tc.genConfig)
		assert.Nil(t, err)
		err = c.Unpack(&gcfg)
		assert.Nil(t, err)
	}
}
