package cef

import (
	"math/rand"
	"testing"
	"text/template"
	"time"

	"github.com/leehinman/spigot/pkg/generator"
)

func TestNext(t *testing.T) {
	tests := []struct {
		seed int64
		want string
	}{
		{seed: 1, want: `CEF:0|Check Point|DNS Trace Log|NS11.0|APPFW|APPFW_SIGNATURE_MATCH|6|dvc=81.37.33.15 c6a4=ef1:c33b:0:f8:c93:8200:3:4ad2 C6a4Label=ring deviceNtDomain=boy oldFileHash=96be6fb77970466a5626fe33408cf9e8 destinationTranslatedZoneURI=https://huge.end.com/share/slip flexDate1=1273775460729 requestMethod=PUT`},
		{seed: 3, want: `CEF:0|Check Point|NetScalar|NS10.0|APPFW|APPFW_STARTURL|7|agentZoneExternalID=00000000-0000-4000-8000-000000000000 dvcpid=15874 slong=63.95873668457833 agentDnsDomain=behave.boy.com flexNumber2Label=preach flexNumber2=92827 rawEvent=Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum. requestCookies=00000000-0000-4000-8000-000000000000 duser=eve`},
		{seed: 4, want: `CEF:0|Check Point|DNS Trace Log|NS11.0|APPFW|APPFW_SAFECOMMERCE_XFORM|5|destinationTranslatedPort=31726 destinationDnsDomain=house.identify.co rt=1273775720347 outcome=failure start=1273773883151 deviceDnsDomain=rely.futuristic.co`},
	}
	templ, err := template.New("cef").Funcs(generator.FunctionMap).Parse(tmpl)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	c := &CEF{templates: []*template.Template{templ}, config: config{
		Type:     Name,
		Vendors:  vendors,
		Products: products,
		Versions: versions,
		Classes:  classes,
		Names:    names,
		Now: func() time.Time {
			return time.Date(2010, time.May, 13, 11, 53, 44, 0, time.FixedZone("-0700", -7*60*60))
		},
		ZeroUUID: true,
	}}
	c.config.Validate() // Populate the remaining fields with the defaults.
	for _, test := range tests {
		rand.Seed(test.seed)
		c.Max = 10
		c.randomize()
		got, err := c.Next()
		if err != nil {
			t.Errorf("unexpected error for c.Next() with seed=%d: %v", err, test.seed)
		}
		if string(got) != test.want {
			t.Errorf("unexpected result for seed=%d:\ngot: %q\nwant:%q", test.seed, got, test.want)
		}
	}
}

var (
	vendors = []string{
		"Check Point",
	}
	products = []string{
		"NetScalar",
		"DNS Trace Log",
	}
	versions = []string{
		"NS10.0",
		"NS11.0",
	}
	classes = []string{
		"APPFW",
	}
	names = []string{
		"APPFW_FIELDCONSISTENCY",
		"APPFW_SAFECOMMERCE",
		"APPFW_SAFECOMMERCE_XFORM",
		"APPFW_SIGNATURE_MATCH",
		"APPFW_STARTURL",
	}
)
