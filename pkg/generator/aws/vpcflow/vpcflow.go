// Package vpcflow generates version 2 AWS vpcflow log messages
//
// For the configuration file there are no options so only the following is needed:
//
//   - generator:
//     type: "aws:vpcflow"
package vpcflow

import (
	"bytes"
	"math/rand"
	"net"
	"text/template"
	"time"

	"github.com/elastic/go-ucfg"
	"github.com/elastic/spigot/pkg/generator"
	"github.com/elastic/spigot/pkg/random"
)

// Name is the name used in the configuration file and the registry.
const Name = "aws:vpcflow"

var (
	actions         = [...]string{"ACCEPT", "REJECT"}
	statuses        = [...]string{"OK", "SKIPDATA", "NODATA"}
	vpcFlowTemplate = "2 123456789010 eni-1235b8ca123456789 {{.SrcAddr}} {{.DstAddr}} {{.SrcPort}} {{.DstPort}} {{.Protocol}} {{.Packets}} {{.Bytes}} {{.Start}} {{.End}} {{.Action}} {{.LogStatus}}"
)

// Vpcflow holds the random fields for a vpcflow record.
type Vpcflow struct {
	//	version     string
	//	accountId   string
	//	interfaceId string
	SrcAddr   net.IP
	DstAddr   net.IP
	SrcPort   int
	DstPort   int
	Protocol  int
	Packets   int
	Bytes     int
	Start     int64
	End       int64
	Action    string
	LogStatus string
	template  *template.Template
}

func init() {
	generator.Register(Name, New)
}

// New is the Factory for Vpcflow objects.
func New(cfg *ucfg.Config, _ int64) (generator.Generator, error) {
	c := defaultConfig()
	if err := cfg.Unpack(&c); err != nil {
		return nil, err
	}

	v := &Vpcflow{}

	t, err := template.New("vpcflow").Funcs(generator.FunctionMap).Parse(vpcFlowTemplate)
	if err != nil {
		return nil, err
	}

	v.template = t
	v.randomize()

	return v, nil
}

// Next produces the next vpcflow record.
//
// Example:
//
// 2 123456789010 eni-1235b8ca123456789 172.31.16.139 172.31.16.21 20641 22 6 20 4249 1418530010 1418530070 ACCEPT OK
func (v *Vpcflow) Next() ([]byte, error) {
	var buf bytes.Buffer

	err := v.template.Execute(&buf, v)
	if err != nil {
		return nil, err
	}

	v.randomize()

	return buf.Bytes(), err

}

func (v *Vpcflow) randomize() {
	v.SrcAddr = random.IPv4()
	v.DstAddr = random.IPv4()
	v.SrcPort = random.Port()
	v.DstPort = random.Port()
	v.Protocol = rand.Intn(256)
	v.Packets = rand.Intn(1048576)
	v.Bytes = v.Packets * 1500
	v.End = time.Now().Unix()
	v.Start = v.End - int64(rand.Intn(60))
	v.Action = actions[rand.Intn(2)]
	if v.Packets == 0 {
		v.LogStatus = statuses[2]
	} else {
		v.LogStatus = statuses[rand.Intn(2)]
	}
}
