package vpcflow

import (
	"bytes"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNext(t *testing.T) {
	tests := []struct {
		SrcAddr   net.IP
		DstAddr   net.IP
		SrcPort   int
		DstPort   int
		Protocol  int
		Packets   int
		Bytes     int
		Action    string
		LogStatus string
	}{
		{
			SrcAddr:   net.ParseIP("66.4.203.154"),
			DstAddr:   net.ParseIP("30.52.197.240"),
			SrcPort:   19911,
			DstPort:   1211,
			Protocol:  129,
			Packets:   643462,
			Bytes:     965193000,
			Action:    "ACCEPT",
			LogStatus: "OK",
		},
	}

	for _, tc := range tests {
		v, err := New()
		assert.Nil(t, err)
		got, err := v.Next()
		assert.Nil(t, err)
		assert.NotEmpty(t, got)
		assert.Equal(t, tc.SrcAddr, v.SrcAddr)
		assert.Equal(t, tc.DstAddr, v.DstAddr)
		assert.Equal(t, tc.SrcPort, v.SrcPort)
		assert.Equal(t, tc.DstPort, v.DstPort)
		assert.Equal(t, tc.Protocol, v.Protocol)
		assert.Equal(t, tc.Packets, v.Packets)
		assert.Equal(t, tc.Bytes, v.Bytes)
		assert.Equal(t, tc.Action, v.Action)
		assert.Equal(t, tc.LogStatus, v.LogStatus)
	}
}

func TestTemplate(t *testing.T) {
	v, err := New()
	assert.Nil(t, err)
	v.SrcAddr = net.ParseIP("66.4.203.154")
	v.DstAddr = net.ParseIP("30.52.197.240")
	v.SrcPort = 2048
	v.DstPort = 80
	v.Protocol = 2
	v.Packets = 10
	v.Bytes = 800
	v.Start = 1024
	v.End = 2048
	v.Action = "ACCEPT"
	v.LogStatus = "OK"
	want := "2 123456789010 eni-1235b8ca123456789 66.4.203.154 30.52.197.240 2048 80 2 10 800 1024 2048 ACCEPT OK"
	var buf bytes.Buffer

	err = v.Template.Execute(&buf, v)
	assert.Nil(t, err)
	assert.Equal(t, []byte(want), buf.Bytes())
}
