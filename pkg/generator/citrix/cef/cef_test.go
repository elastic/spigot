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
		{seed: 1, want: `Jan 2 03:04:05 <mark.emerg> 118.9.14.112 CEF:1|Citrix|NetScalar|NS10.0|APPFW|APPFW_FIELDCONSISTENCY|5|src=95.181.74.208 geolocation=NorthAmerica.US.Arizona.Tucson.*.* spt=24561 method=GET request=http://vpx247.example.net/FFC/CreditCardMind.html msg=Disallow Illegal URL. cn1=445 cn2=23237 cs1=pr_ffc cs2=PPE3 cs3=448615bbda08313f6a8eb668d20bf505 cs4=ALERT cs5=2022 cs6=phishing act=not blocked`},
		{seed: 3, want: `Jan 02 03:04:05 <local5.error> 244.161.164.196 CEF:1|Citrix|NetScalar|NS11.0|APPFW|APPFW_SAFECOMMERCE_XFORM|7|src=157.155.176.203 geolocation=Unknown spt=29735 method=GET request=http://aaron.stratum8.net/FFC/wwwboard/passwd.txt msg=Field consistency check failed for field passwd cn1=278 cn2=29074 cs1=pr_ffc cs2=PPE4 cs3=06d37841b74bcbbdf8987a19dcddc8e9 cs4=ALERT cs5=2022 cs6=web-cgi act=blocked`},
		{seed: 4, want: `Jan 2 03:04:05 <local4.error> 63.132.159.242 CEF:1|Citrix|NetScalar|NS11.0|APPFW|APPFW_SAFECOMMERCE|8|src=201.132.96.184 spt=25717 method=GET request=http://aaron.stratum8.net/FFC/wwwboard/passwd.txt msg=Field consistency check failed for field passwd cn1=586 cn2=78840 cs1=pr_ffc cs2=PPE8 cs3=ab6d79345fe5e99adf9ddd3d1dbfe5db cs4=INFO cs5=2022 cs6=sql-injection act=transformed`},
	}
	now, err := time.Parse(time.RFC3339, "1970-01-02T03:04:05Z")
	if err != nil {
		t.Fatal(err)
	}
	templ, err := template.New("cef").Funcs(generator.FunctionMap).Parse(tmpl)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	c := &CEF{templates: []*template.Template{templ}}
	for _, test := range tests {
		rand.Seed(test.seed)
		c.randomize()
		c.Timestamp = now
		got, err := c.Next()
		if err != nil {
			t.Errorf("unexpected error for c.Next() with seed=%d: %v", err, test.seed)
		}
		if string(got) != test.want {
			t.Errorf("unexpected result for seed=%d:\ngot: %s\nwant:%s", test.seed, got, test.want)
		}
	}
}
