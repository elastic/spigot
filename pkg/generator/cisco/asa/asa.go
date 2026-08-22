// Package asa implements the generator for Cisco ASA logs.
//
// Configuration file supports including timestamps in log messages
//
//	generator:
//	  type: cisco:asa
//	  include_timestamp: true
package asa

import (
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"text/template"
	"time"

	"github.com/elastic/go-ucfg"
	"github.com/elastic/spigot/pkg/generator"
	"github.com/elastic/spigot/pkg/random"
)

// Name is the name of the generator in the configuration file and registry
const Name = "cisco:asa"

var (
	asa106023    = "{{if .IncludeTimestamp}}{{.Timestamp.Format \"Jan 02 2006 03:04:05\"}}: {{end}}%ASA-4-106023: Deny {{.Protocol | ToLower}} src {{.SrcInt}}:{{.SrcAddr}}/{{.SrcPort}} dst {{.DstInt}}:{{.DstAddr}}/{{.DstPort}} type {{.Type}} code {{.Code}} by {{.AccessGroup | ToLower}} \"{{.AclId}}\" [0x8ed66b60, 0xf8852875]"
	asa302013    = "{{if .IncludeTimestamp}}{{.Timestamp.Format \"Jan 02 2006 03:04:05\"}}: {{end}}%ASA-6-302013: Built {{.Direction}} TCP connection {{.ConnectionId}} for {{.SrcInt}}:{{.SrcAddr}}/{{.SrcPort}} ({{.Map1Addr}}/{{.Map1Port}}) to {{.DstInt}}:{{.DstAddr}}/{{.DstPort}} ({{.Map2Addr}}/{{.Map2Port}})"
	asa302014    = "{{if .IncludeTimestamp}}{{.Timestamp.Format \"Jan 02 2006 03:04:05\"}}: {{end}}%ASA-6-302014: Teardown TCP connection {{.ConnectionId}} for {{.SrcInt}}:{{.SrcAddr}}/{{.SrcPort}} to {{.DstInt}}:{{.DstAddr}}/{{.DstPort}} duration {{.Duration}} bytes {{.Bytes}} {{.Reason}}"
	asa305011    = "{{if .IncludeTimestamp}}{{.Timestamp.Format \"Jan 02 2006 03:04:05\"}}: {{end}}%ASA-6-305011: Built {{.TranslationType}} {{.Protocol}} translation from {{.SrcInt}}:{{.SrcAddr}}/{{.SrcPort}} to {{.DstInt}}:{{.DstAddr}}/{{.DstPort}}"
	msgTemplates = [...]string{
		asa106023,
		asa302013,
		asa302014,
		asa305011,
	}
	directions       = [...]string{"inbound", "outbound"}
	protocols        = [...]string{"TCP", "UDP"}
	translationTypes = [...]string{"dynamic", "static"}
	reasons          = [...]string{
		"Conn-timeout",
		"Deny Terminate",
		"Failover primary closed",
		"FIN Timeout",
		"Flow closed by inspection",
		"Flow terminated by IPS",
		"Flow reset by IPS",
		"Flow terminated by TCP Intercept",
		"Flow timed out",
		"Flow timed out with reset",
		"Flow is a loopback",
		"Free the flow created as result of packet injection",
		"Invalid SYN",
		"IPS fail-close",
		"No interfaces associated with zone",
		"No valid adjacency",
		"Pinhole Timeout",
		"Probe maximum retries of retransmission exceeded",
		"Probe maximum retransmission time elapsed",
		"Probe received RST",
		"Probe received FIN",
		"Probe completed",
		"Route change",
		"SYN Control",
		"SYN Timeout",
		"TCP bad retransmission",
		"TCP FINs",
		"TCP Invalid SYN",
		"TCP Reset - APPLIANCE",
		"TCP Reset - I",
		"TCP Reset - O",
		"TCP segment partial overlap",
		"TCP unexpected window size variation",
		"Tunnel has been torn down",
		"Unauth Deny",
		"Unknown",
		"Xlate Clear",
	}
)

type Asa struct {
	AccessGroup      string
	AclId            string
	Bytes            int
	Code             int
	ConnectionId     int
	Direction        string
	DstAddr          net.IP
	DstInt           string
	DstPort          int
	DstUser          string
	Duration         string
	IncludeTimestamp bool
	Map1Addr         net.IP
	Map1Port         int
	Map2Addr         net.IP
	Map2Port         int
	Protocol         string
	Reason           string
	SrcAddr          net.IP
	SrcInt           string
	SrcPort          int
	SrcUser          string
	Timestamp        time.Time
	TranslationType  string
	Type             int
	templates        []*template.Template
}

func init() {
	generator.Register(Name, New)
}

// New is Factory for the asa generator
func New(cfg *ucfg.Config, _ int64) (generator.Generator, error) {
	c := defaultConfig()
	if err := cfg.Unpack(&c); err != nil {
		return nil, err
	}

	a := &Asa{
		IncludeTimestamp: c.IncludeTimestamp,
	}
	a.randomize()

	for i, v := range msgTemplates {
		t, err := template.New(strconv.Itoa(i)).Funcs(generator.FunctionMap).Parse(v)
		if err != nil {
			return nil, err
		}
		a.templates = append(a.templates, t)
	}

	return a, nil
}

// Next produces the next asa log entry
func (a *Asa) Next() ([]byte, error) {
	var buf bytes.Buffer

	err := a.templates[rand.Intn(len(a.templates))].Execute(&buf, a)
	if err != nil {
		return nil, err
	}

	a.randomize()

	return buf.Bytes(), err
}

func (a *Asa) randomize() {
	a.SrcInt = "SrcInt"
	a.SrcUser = "SrcUser"
	a.DstInt = "DstInt"
	a.DstUser = "DstUser"
	a.AccessGroup = "Access-Group"
	a.AclId = "AclId"
	a.Protocol = protocols[rand.Intn(len(protocols))]
	a.TranslationType = translationTypes[rand.Intn(len(translationTypes))]
	a.ConnectionId = rand.Intn(65536)
	a.Duration = fmt.Sprintf("%01d:%02d:%02d", rand.Intn(4), rand.Intn(60), rand.Intn(60))
	a.Bytes = rand.Intn(65536)
	a.Reason = reasons[rand.Intn(len(reasons))]
	a.SrcAddr = random.IPv4()
	a.SrcPort = random.Port()
	a.DstAddr = random.IPv4()
	a.DstPort = random.Port()
	a.Type = rand.Intn(64)
	a.Code = rand.Intn(64)
	a.Direction = directions[rand.Intn(len(directions))]
	a.Map1Addr = random.IPv4()
	a.Map1Port = random.Port()
	a.Map2Addr = random.IPv4()
	a.Map2Port = random.Port()
	a.Timestamp = time.Now()
}
