package asa

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
		"106023": {template: asa106023, expected: "%ASA-4-106023: Deny udp src SrcInt:144.254.210.24/18340 dst DstInt:141.249.228.131/23215 type 34 code 49 by access-group \"AclId\" [0x8ed66b60, 0xf8852875]"},
		"302013": {template: asa302013, expected: "%ASA-6-302013: Built inbound TCP connection 19911 for SrcInt:144.254.210.24/18340 (53.42.9.120/30347) to DstInt:141.249.228.131/23215 (43.185.8.75/16165)"},
		"302014": {template: asa302014, expected: "%ASA-6-302014: Teardown TCP connection 19911 for SrcInt:144.254.210.24/18340 to DstInt:141.249.228.131/23215 duration 3:01:18 bytes 52025 Xlate Clear"},
		"305011": {template: asa305011, expected: "%ASA-6-305011: Built static UDP translation from SrcInt:144.254.210.24/18340 to DstInt:141.249.228.131/23215"},
	}
	for name, tc := range tests {
		rand.Seed(1)
		a := &Asa{}
		templ, err := template.New(name).Funcs(generator.FunctionMap).Parse(tc.template)
		assert.Nil(t, err)
		a.templates = []*template.Template{templ}
		a.randomize()
		got, err := a.Next()
		assert.Nil(t, err, name)
		assert.Equal(t, []byte(tc.expected), got, name)
	}
}
