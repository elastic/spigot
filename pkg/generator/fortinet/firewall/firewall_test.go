package firewall

import (
	"math/rand"
	"testing"
	"text/template"
	"time"

	"github.com/elastic/spigot/pkg/generator"
	"github.com/stretchr/testify/assert"
)

func TestNext(t *testing.T) {
	tests := map[string]struct {
		template string
		expected string
	}{
		"EventUser": {template: eventUserTemplate,
			expected: "date=1970-01-02 time=03:04:05 devname=\"testswitch3\" devid=\"testrouter\" logid=\"0123456789\" type=\"event\" subtype=\"user\" level=\"error\" vd=\"root\" eventtime=97445 tz=\"-0500\" logdesc=\"FSSO logon authentication status\" srcip=142.155.32.170 user=\"user07\" server=\"srv7\" action=\"FSSO-logon\" msg=\"FSSO-logon event from FSSO_srv7: user user07 logged on 142.155.32.170\""},
		"EventSystem": {template: eventSystemTemplate,
			expected: "date=1970-01-02 time=03:04:05 devname=\"testswitch3\" devid=\"testrouter\" logid=\"0123456789\" type=\"event\" subtype=\"system\" level=\"error\" vd=\"root\" eventtime=97445 tz=\"-0500\" logdesc=\"FortiSandbox AV database updated\" version=\"1.522479\" msg=\"FortiSandbox AV database updated\""},
		"UtmDns": {template: utmDnsTemplate,
			expected: "date=1970-01-02 time=03:04:05 devname=\"testswitch3\" devid=\"testrouter\" logid=\"0123456789\" type=\"utm\" subtype=\"dns\" eventtype=\"dns-query\" level=\"error\" vd=\"root\" eventtime=97445 tz=\"-0500\" policyid=57 sessionid=53932 srcip=142.155.32.170 srcport=1211 srcintf=\"int0\" srcintfrole=\"internal\" dstip=2.11.181.108 dstport=53 dstintf=\"int4\" dstintfrole=\"wan\" proto=6 profile=\"elastictest\" xid=26 qname=\"elastic.co\" qtype=\"A\" qtypeval=1 qclass=\"IN\""},
		"TrafficForward": {template: trafficForwardTemplate,
			expected: "date=1970-01-02 time=03:04:05 devname=\"testswitch3\" devid=\"testrouter\" logid=\"0123456789\" type=\"traffic\" subtype=\"forward\" level=\"error\" vd=\"root\" eventtime=97445 srcip=142.155.32.170 srcport=1211 srcintf=\"int0\" srcintfrole=\"internal\" dstip=2.11.181.108 dstport=53638 dstintf=\"int4\" dstintfrole=\"wan\" sessionid=53932 proto=6 action=\"accept\" policyid=57 policytype=\"policy\" service=\"SNMP\" dstcountry=\"Reserved\" srccountry=\"Reserved\" trandisp=\"noop\" duration=994 sentbyte=24247500 rcvdbyte=24247500 sentpkt=16165 appcat=\"unscanned\" crscore=30 craction=131072 crlevel=\"high\""},
	}
	test_time, err := time.Parse(time.RFC3339, "1970-01-02T03:04:05Z")
	assert.Nil(t, err)
	for name, tc := range tests {
		rand.Seed(1)
		f := &Firewall{}
		f.randomize()
		templ, err := template.New(name).Funcs(generator.FunctionMap).Parse(tc.template)
		assert.Nil(t, err)
		f.Templates = []*template.Template{templ}
		f.Date = test_time
		got, err := f.Next()
		assert.Nil(t, err)
		assert.Equal(t, []byte(tc.expected), got, name)
	}
}
