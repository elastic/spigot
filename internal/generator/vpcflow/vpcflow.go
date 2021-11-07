package vpcflow

import (
	"bytes"
	"math/rand"
	"net"
	"text/template"
	"time"
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

func New() (v *Vpcflow, err error) {
	v = &Vpcflow{}
	t, err := template.New("vpcflow").Parse("2 123456789010 eni-1235b8ca123456789 {{.SrcAddr}} {{.DstAddr}} {{.SrcPort}} {{.DstPort}} {{.Protocol}} {{.Packets}} {{.Bytes}} {{.Start}} {{.End}} {{.Action}} {{.LogStatus}}")
	if err != nil {
		return nil, err
	}
	v.Template = t
	return v, nil
}

func (v *Vpcflow) Next() ([]byte, error) {
	var buf bytes.Buffer

	v.SrcAddr = generateAddr()
	v.DstAddr = generateAddr()
	v.SrcPort = rand.Intn(65536)
	v.DstPort = rand.Intn(65536)
	v.Protocol = rand.Intn(256)
	v.Packets = rand.Intn(1048576)
	v.Bytes = v.Packets * 1500
	v.Start = time.Now().Unix() - int64(rand.Intn(60))
	v.End = v.Start + int64(v.Bytes/800)
	v.Action = ACTIONS[rand.Intn(2)]
	if v.Packets == 0 {
		v.LogStatus = STATUSES[2]
	} else {
		v.LogStatus = STATUSES[rand.Intn(2)]
	}

	err := v.Template.Execute(&buf, v)
	if err != nil {
		panic(err)
	}
	return buf.Bytes(), err

}

func generateAddr() net.IP {
	u32 := rand.Uint32()
	return net.IPv4(byte(u32&0xff), byte((u32>>8)&0xff), byte((u32>>16)&0xff), byte((u32>>24)&0xff))
}
