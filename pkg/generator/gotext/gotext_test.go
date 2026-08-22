package gotext

import (
	"testing"

	"github.com/elastic/go-ucfg"
	"github.com/stretchr/testify/assert"
)

func TestNext(t *testing.T) {

	type config struct {
		Type   string                 `config:"type" validate:"required"`
		Config map[string]interface{} `config:"config" validate:"required"`
	}

	testCases := []struct {
		genConfig map[string]interface{}
		errString string
		expected  string
	}{
		{
			genConfig: map[string]interface{}{
				"type": "gotext",
				"config": map[string]interface{}{
					"name": "test",
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
			},
			errString: "",
			expected:  "",
		},
	}

	for _, tc := range testCases {
		var c config
		u_config, err := ucfg.NewFrom(tc.genConfig)
		assert.Nil(t, err)
		err = u_config.Unpack(&c)
		assert.Nil(t, err)
		g, err := New(u_config, 100)
		assert.Nil(t, err)
		for i := 0; i < 100; i++ {
			_, err := g.Next()
			assert.Nil(t, err)
		}
	}
}
