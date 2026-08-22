package firewall

import (
	"encoding/json"
	"math/rand"
	"net"
	"testing"

	"github.com/elastic/go-ucfg"
	"github.com/stretchr/testify/assert"
)

func TestNext(t *testing.T) {
	tests := map[string]struct {
		config map[string]interface{}
		want   Firewall
	}{
		"netflow": {
			config: map[string]interface{}{
				"type":       Name,
				"event_type": "netflow",
			},
			want: Firewall{
				FirewallName:     "Firewall-81",
				AvailabilityZone: "eu-west-1c",
				Event: EventData{
					FlowID:    6129484611666145821,
					EventType: "netflow",
					SrcIP:     net.ParseIP("118.9.14.112"),
					SrcPort:   34177,
					DstIP:     net.ParseIP("12.163.211.175"),
					DstPort:   52025,
					Proto:     "UDP",
					Netflow: &NetflowData{
						Pkts:   94,
						Bytes:  64579,
						Age:    0,
						MinTTL: 72,
						MaxTTL: 72,
					},
				},
			},
		},
		"alert": {
			config: map[string]interface{}{
				"event_type": "alert",
			},
			want: Firewall{
				FirewallName:     "Firewall-81",
				AvailabilityZone: "eu-west-1c",
				Event: EventData{
					FlowID:    6129484611666145821,
					EventType: "alert",
					SrcIP:     net.ParseIP("118.9.14.112"),
					SrcPort:   34177,
					DstIP:     net.ParseIP("12.163.211.175"),
					DstPort:   52025,
					Proto:     "UDP",
					Alert: &AlertData{
						Action:      "allowed",
						SignatureID: 840,
						Rev:         198,
						Signature:   "Signature-840",
						Category:    "Category-11",
						Severity:    0,
					},
				},
			},
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			rand.Seed(1)
			var got Firewall
			g, err := New(ucfg.MustNewFrom(tc.config), 0)
			if err != nil {
				t.Fatal(err)
			}

			data, err := g.Next()
			assert.NoError(t, err)

			err = json.Unmarshal(data, &got)
			assert.NoError(t, err)

			// Clear dynamic fields.
			got.EventTimestamp = ""
			got.Event.Timestamp = ""
			if got.Event.Netflow != nil {
				got.Event.Netflow.Start = ""
				got.Event.Netflow.End = ""
			}

			assert.Equal(t, tc.want, got)
		})
	}
}

func BenchmarkGenerator_Next(b *testing.B) {
	b.ReportAllocs()

	rand.Seed(1)
	g, err := New(ucfg.New(), 0)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		_, _ = g.Next()

	}
}
