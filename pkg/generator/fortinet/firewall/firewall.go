package firewall

import (
	"bytes"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/elastic/go-ucfg"
	"github.com/leehinman/spigot/pkg/generator"
	"github.com/leehinman/spigot/pkg/random"
)

var (
	EventUserTemplate      = "date={{.Date.UTC.Format \"2006-01-02\"}} time={{.Date.UTC.Format \"03:04:05\"}} devname=\"{{.DevName}}\" devid=\"{{.DevId}}\" logid=\"{{.LogId}}\" type=\"event\" subtype=\"user\" level=\"{{.Level}}\" vd=\"{{.Vd}}\" eventtime={{.Date.Unix}} tz=\"{{.Timezone}}\" logdesc=\"FSSO logon authentication status\" srcip={{.SrcIp}} user=\"{{.User}}\" server=\"{{.Server}}\" action=\"FSSO-logon\" msg=\"FSSO-logon event from FSSO_{{.Server}}: user {{.User}} logged on {{.SrcIp}}\""
	EventSystemTemplate    = "date={{.Date.UTC.Format \"2006-01-02\"}} time={{.Date.UTC.Format \"03:04:05\"}} devname=\"{{.DevName}}\" devid=\"{{.DevId}}\" logid=\"{{.LogId}}\" type=\"event\" subtype=\"system\" level=\"{{.Level}}\" vd=\"{{.Vd}}\" eventtime={{.Date.Unix}} tz=\"{{.Timezone}}\" logdesc=\"FortiSandbox AV database updated\" version=\"1.522479\" msg=\"FortiSandbox AV database updated\""
	UtmDnsTemplate         = "date={{.Date.UTC.Format \"2006-01-02\"}} time={{.Date.UTC.Format \"03:04:05\"}} devname=\"{{.DevName}}\" devid=\"{{.DevId}}\" logid=\"{{.LogId}}\" type=\"utm\" subtype=\"dns\" eventtype=\"dns-query\" level=\"{{.Level}}\" vd=\"{{.Vd}}\" eventtime={{.Date.Unix}} tz=\"{{.Timezone}}\" policyid={{.PolicyId}} sessionid={{.SessionId}} srcip={{.SrcIp}} srcport={{.SrcPort}} srcintf=\"{{.Interface1}}\" srcintfrole=\"{{.InterfaceRole1}}\" dstip={{.DstIp}} dstport=53 dstintf=\"{{.Interface2}}\" dstintfrole=\"{{.InterfaceRole2}}\" proto={{.Protocol}} profile=\"elastictest\" xid={{.XId}} qname=\"{{.QueryName}}\" qtype=\"{{.QueryType}}\" qtypeval=1 qclass=\"IN\""
	TrafficForwardTemplate = "date={{.Date.UTC.Format \"2006-01-02\"}} time={{.Date.UTC.Format \"03:04:05\"}} devname=\"{{.DevName}}\" devid=\"{{.DevId}}\" logid=\"{{.LogId}}\" type=\"traffic\" subtype=\"forward\" level=\"{{.Level}}\" vd=\"{{.Vd}}\" eventtime={{.Date.Unix}} srcip={{.SrcIp}} srcport={{.SrcPort}} srcintf=\"{{.Interface1}}\" srcintfrole=\"{{.InterfaceRole1}}\" dstip={{.DstIp}} dstport={{.DstPort}} dstintf=\"{{.Interface2}}\" dstintfrole=\"{{.InterfaceRole2}}\" sessionid={{.SessionId}} proto={{.Protocol}} action=\"{{.TrafficAction}}\" policyid={{.PolicyId}} policytype=\"policy\" service=\"SNMP\" dstcountry=\"Reserved\" srccountry=\"Reserved\" trandisp=\"noop\" duration={{.Duration}} sentbyte={{.SentBytes}} rcvdbyte={{.SentBytes}} sentpkt={{.SentPackets}} appcat=\"unscanned\" crscore=30 craction=131072 crlevel=\"high\""
	MsgTemplates           = [...]string{
		EventUserTemplate,
		EventSystemTemplate,
		UtmDnsTemplate,
		TrafficForwardTemplate,
	}
	FuncMap = template.FuncMap{
		"ToLower": strings.ToLower,
		"ToUpper": strings.ToUpper,
	}
	Users          = [...]string{"user01", "user02", "user03", "user04", "user05", "user06", "user07"}
	Levels         = [...]string{"warning", "notice", "information", "error"}
	Interfaces     = [...]string{"int0", "int1", "int2", "int3", "int4", "int5", "int6", "int7"}
	Roles          = [...]string{"lan", "wan", "internal", "external", "inbound", "outbound"}
	Protocols      = [...]int{6, 17}
	Queries        = [...]string{"example.com", "google.com", "amazon.com", "elastic.co", "apple.com", "facebook.com", "microsoft.com"}
	QueryTypes     = [...]string{"A", "AAAA"}
	Servers        = [...]string{"srv0", "srv1", "srv2", "srv3", "srv4", "srv5", "srv6", "srv7"}
	TrafficActions = [...]string{"deny", "accept"}
)

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
	generator.Register("fortinet:firewall", New)
}

func New(cfg *ucfg.Config) (generator.Generator, error) {
	f := &Firewall{
		DevName:  "testswitch3",
		DevId:    "testrouter",
		LogId:    "0123456789",
		Timezone: "-0500",
	}
	f.randomize()
	for i, v := range MsgTemplates {
		t, err := template.New(strconv.Itoa(i)).Funcs(FuncMap).Parse(v)
		if err != nil {
			return nil, err
		}
		f.Templates = append(f.Templates, t)
	}
	return f, nil
}

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
	f.Date = time.Now()
	f.Vd = "root"
	f.User = Users[rand.Intn(len(Users))]
	f.Server = Servers[rand.Intn(len(Servers))]
	f.SrcIp = random.IPv4()
	f.SrcPort = random.Port()
	f.DstIp = random.IPv4()
	f.DstPort = random.Port()
	f.PolicyId = rand.Intn(256)
	f.SessionId = rand.Intn(65536)
	f.Interface1 = Interfaces[rand.Intn(len(Interfaces))]
	f.Interface2 = Interfaces[rand.Intn(len(Interfaces))]
	f.InterfaceRole1 = Roles[rand.Intn(len(Roles))]
	f.InterfaceRole2 = Roles[rand.Intn(len(Roles))]
	f.Protocol = Protocols[rand.Intn(len(Protocols))]
	f.QueryName = Queries[rand.Intn(len(Queries))]
	f.QueryType = QueryTypes[rand.Intn(len(QueryTypes))]
	f.XId = rand.Intn(256)
	f.Level = Levels[rand.Intn(len(Levels))]
	f.TrafficAction = TrafficActions[rand.Intn(len(TrafficActions))]
	f.SentPackets = rand.Intn(65536)
	f.SentBytes = f.SentPackets * 1500
	f.Duration = rand.Intn(1024)
}
