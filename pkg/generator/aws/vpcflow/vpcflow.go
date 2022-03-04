package vpcflow

import (
	"bytes"
	"math/rand"
	"net"
	"text/template"
	"time"

	"github.com/elastic/go-ucfg"
	"github.com/leehinman/spigot/pkg/generator"
	"github.com/leehinman/spigot/pkg/random"
)

var ACTIONS = [...]string{"ACCEPT", "REJECT"}
var STATUSES = [...]string{"OK", "SKIPDATA", "NODATA"}

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
	Template  *template.Template
}

func init() {
	generator.Register("aws:vpcflow", New)
}

func New(cfg *ucfg.Config) (generator.Generator, error) {
	c := defaultConfig()
	if err := cfg.Unpack(&c); err != nil {
		return nil, err
	}
	v := &Vpcflow{}
	t, err := template.New("vpcflow").Parse("2 123456789010 eni-1235b8ca123456789 {{.SrcAddr}} {{.DstAddr}} {{.SrcPort}} {{.DstPort}} {{.Protocol}} {{.Packets}} {{.Bytes}} {{.Start}} {{.End}} {{.Action}} {{.LogStatus}}")
	if err != nil {
		return nil, err
	}
	v.Template = t
	return v, nil
}

func (v *Vpcflow) Next() ([]byte, error) {
	var buf bytes.Buffer

	v.SrcAddr = random.IPv4()
	v.DstAddr = random.IPv4()
	v.SrcPort = random.Port()
	v.DstPort = random.Port()
	v.Protocol = rand.Intn(256)
	v.Packets = rand.Intn(1048576)
	v.Bytes = v.Packets * 1500
	v.End = time.Now().Unix()
	v.Start = v.End - int64(rand.Intn(60))
	v.Action = ACTIONS[rand.Intn(2)]
	if v.Packets == 0 {
		v.LogStatus = STATUSES[2]
	} else {
		v.LogStatus = STATUSES[rand.Intn(2)]
	}

	err := v.Template.Execute(&buf, v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err

}
