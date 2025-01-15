package vpcflow

import (
	"math/rand"
	"testing"
	"text/template"

	"github.com/elastic/spigot/pkg/generator"
	"github.com/stretchr/testify/assert"
)

func TestNext(t *testing.T) {
	tests := map[string]struct {
		template string
		expected string
	}{
		"vpcflow v2": {
			template: vpcFlowTemplate,
			expected: "2 123456789010 eni-1235b8ca123456789 66.4.203.154 30.52.197.240 19911 1211 129 643462 965193000 2 42 ACCEPT OK",
		},
	}

	for name, tc := range tests {
		rand.Seed(1)
		v := &Vpcflow{}
		tmpl, err := template.New(name).Funcs(generator.FunctionMap).Parse(tc.template)
		assert.Nil(t, err, name)
		v.template = tmpl
		v.randomize()
		v.End = 42
		v.Start = 2
		got, err := v.Next()
		assert.Nil(t, err, name)
		assert.Equal(t, []byte(tc.expected), got, name)
	}
}
