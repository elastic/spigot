// Package firewall generates Fortinet Firewall log messages
//
// For the configuration file there are no options so only the following is needed:
//
//   - generator:
//     type: "fortinet:firewall"
package firewall

import (
	"bytes"
	"math/rand"
	"net"
	"strconv"
	"text/template"
	"time"

	"github.com/elastic/go-ucfg"
	"github.com/elastic/spigot/pkg/generator"
	"github.com/elastic/spigot/pkg/random"
)

// Name is the name used in the configuration file and the registry.
const Name = "fortinet:firewall"

var (
	eventUserTemplate      = "date={{.Date.UTC.Format \"2006-01-02\"}} time={{.Date.UTC.Format \"03:04:05\"}} devname=\"{{.DevName}}\" devid=\"{{.DevId}}\" logid=\"{{.LogId}}\" type=\"event\" subtype=\"user\" level=\"{{.Level}}\" vd=\"{{.Vd}}\" eventtime={{.Date.Unix}} tz=\"{{.Timezone}}\" logdesc=\"FSSO logon authentication status\" srcip={{.SrcIp}} user=\"{{.User}}\" server=\"{{.Server}}\" action=\"FSSO-logon\" msg=\"FSSO-logon event from FSSO_{{.Server}}: user {{.User}} logged on {{.SrcIp}}\""
	eventSystemTemplate    = "date={{.Date.UTC.Format \"2006-01-02\"}} time={{.Date.UTC.Format \"03:04:05\"}} devname=\"{{.DevName}}\" devid=\"{{.DevId}}\" logid=\"{{.LogId}}\" type=\"event\" subtype=\"system\" level=\"{{.Level}}\" vd=\"{{.Vd}}\" eventtime={{.Date.Unix}} tz=\"{{.Timezone}}\" logdesc=\"FortiSandbox AV database updated\" version=\"1.522479\" msg=\"FortiSandbox AV database updated\""
	utmDnsTemplate         = "date={{.Date.UTC.Format \"2006-01-02\"}} time={{.Date.UTC.Format \"03:04:05\"}} devname=\"{{.DevName}}\" devid=\"{{.DevId}}\" logid=\"{{.LogId}}\" type=\"utm\" subtype=\"dns\" eventtype=\"dns-query\" level=\"{{.Level}}\" vd=\"{{.Vd}}\" eventtime={{.Date.Unix}} tz=\"{{.Timezone}}\" policyid={{.PolicyId}} sessionid={{.SessionId}} srcip={{.SrcIp}} srcport={{.SrcPort}} srcintf=\"{{.Interface1}}\" srcintfrole=\"{{.InterfaceRole1}}\" dstip={{.DstIp}} dstport=53 dstintf=\"{{.Interface2}}\" dstintfrole=\"{{.InterfaceRole2}}\" proto={{.Protocol}} profile=\"elastictest\" xid={{.XId}} qname=\"{{.QueryName}}\" qtype=\"{{.QueryType}}\" qtypeval=1 qclass=\"IN\""
	trafficForwardTemplate = "date={{.Date.UTC.Format \"2006-01-02\"}} time={{.Date.UTC.Format \"03:04:05\"}} devname=\"{{.DevName}}\" devid=\"{{.DevId}}\" logid=\"{{.LogId}}\" type=\"traffic\" subtype=\"forward\" level=\"{{.Level}}\" vd=\"{{.Vd}}\" eventtime={{.Date.Unix}} srcip={{.SrcIp}} srcport={{.SrcPort}} srcintf=\"{{.Interface1}}\" srcintfrole=\"{{.InterfaceRole1}}\" dstip={{.DstIp}} dstport={{.DstPort}} dstintf=\"{{.Interface2}}\" dstintfrole=\"{{.InterfaceRole2}}\" sessionid={{.SessionId}} proto={{.Protocol}} action=\"{{.TrafficAction}}\" policyid={{.PolicyId}} policytype=\"policy\" service=\"SNMP\" dstcountry=\"Reserved\" srccountry=\"Reserved\" trandisp=\"noop\" duration={{.Duration}} sentbyte={{.SentBytes}} rcvdbyte={{.SentBytes}} sentpkt={{.SentPackets}} appcat=\"unscanned\" crscore=30 craction=131072 crlevel=\"high\""
	msgTemplates           = [...]string{
		eventUserTemplate,
		eventSystemTemplate,
		utmDnsTemplate,
		trafficForwardTemplate,
	}
	users          = [...]string{"user01", "user02", "user03", "user04", "user05", "user06", "user07"}
	levels         = [...]string{"warning", "notice", "information", "error"}
	interfaces     = [...]string{"int0", "int1", "int2", "int3", "int4", "int5", "int6", "int7"}
	roles          = [...]string{"lan", "wan", "internal", "external", "inbound", "outbound"}
	protocols      = [...]int{6, 17}
	queries        = [...]string{"example.com", "google.com", "amazon.com", "elastic.co", "apple.com", "facebook.com", "microsoft.com"}
	queryTypes     = [...]string{"A", "AAAA"}
	servers        = [...]string{"srv0", "srv1", "srv2", "srv3", "srv4", "srv5", "srv6", "srv7"}
	trafficActions = [...]string{"deny", "accept"}
)

// Firewall holds the random fields for a firewall record
type Firewall struct {
	Date           time.Time
	DevId          string
	DevName        string
	Direction      string
	DstIp          net.IP
	DstPort        int
	Duration       int
	Interface1     string
	Interface2     string
	InterfaceRole1 string
	InterfaceRole2 string
	Level          string
	LogId          string
	PolicyId       int
	Protocol       int
	QueryName      string
	QueryType      string
	ReceivedBytes  int
	SentBytes      int
	SentPackets    int
	Server         string
	SessionId      int
	SrcIp          net.IP
	SrcPort        int
	Templates      []*template.Template
	Timezone       string
	TrafficAction  string
	User           string
	Vd             string
	XId            int
}

func init() {
	generator.Register(Name, New)
}

// New is the Factory for Firewall objects.
func New(cfg *ucfg.Config, _ int64) (generator.Generator, error) {
	c := defaultConfig()
	if err := cfg.Unpack(&c); err != nil {
		return nil, err
	}

	f := &Firewall{}
	f.randomize()

	for i, v := range msgTemplates {
		t, err := template.New(strconv.Itoa(i)).Funcs(generator.FunctionMap).Parse(v)
		if err != nil {
			return nil, err
		}
		f.Templates = append(f.Templates, t)
	}
	return f, nil
}

// Next produces the next firewall record.
//
// Example:
//
// date=1970-01-02 time=03:04:05 devname=\"testswitch3\" devid=\"testrouter\" logid=\"0123456789\" type=\"event\" subtype=\"user\" level=\"error\" vd=\"root\" eventtime=97445 tz=\"-0500\" logdesc=\"FSSO logon authentication status\" srcip=142.155.32.170 user=\"user07\" server=\"srv7\" action=\"FSSO-logon\" msg=\"FSSO-logon event from FSSO_srv7: user user07 logged on 142.155.32.170\"
func (f *Firewall) Next() ([]byte, error) {
	var buf bytes.Buffer

	err := f.Templates[rand.Intn(len(f.Templates))].Execute(&buf, f)
	if err != nil {
		return nil, err
	}

	//randomize after evaluating template to make testing easier
	f.randomize()
	return buf.Bytes(), err
}

func (f *Firewall) randomize() {
	f.DevName = "testswitch3"
	f.DevId = "testrouter"
	f.LogId = "0123456789"
	f.Timezone = "-0500"
	f.Date = time.Now()
	f.Vd = "root"
	f.User = users[rand.Intn(len(users))]
	f.Server = servers[rand.Intn(len(servers))]
	f.SrcIp = random.IPv4()
	f.SrcPort = random.Port()
	f.DstIp = random.IPv4()
	f.DstPort = random.Port()
	f.PolicyId = rand.Intn(256)
	f.SessionId = rand.Intn(65536)
	f.Interface1 = interfaces[rand.Intn(len(interfaces))]
	f.Interface2 = interfaces[rand.Intn(len(interfaces))]
	f.InterfaceRole1 = roles[rand.Intn(len(roles))]
	f.InterfaceRole2 = roles[rand.Intn(len(roles))]
	f.Protocol = protocols[rand.Intn(len(protocols))]
	f.QueryName = queries[rand.Intn(len(queries))]
	f.QueryType = queryTypes[rand.Intn(len(queryTypes))]
	f.XId = rand.Intn(256)
	f.Level = levels[rand.Intn(len(levels))]
	f.TrafficAction = trafficActions[rand.Intn(len(trafficActions))]
	f.SentPackets = rand.Intn(65536)
	f.SentBytes = f.SentPackets * 1500
	f.Duration = rand.Intn(1024)
}
