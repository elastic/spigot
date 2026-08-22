// Package firewall generates AWS Network Firewall log messages. It currently supports generating
// netflow and alert events, and supports generating detailed TCP netflow and HTTP alert logs.
//
// Configuration:
//
//		 event_type: Specify the type of event to generate, or leave blank for random.
//		             Valid values are: alert, netflow.
//
//	  - generator:
//	      type: aws:firewall
//		     event_type: netflow
package firewall

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/elastic/go-ucfg"
	"github.com/elastic/spigot/pkg/generator"
	"github.com/elastic/spigot/pkg/random"
)

// Name is the name used in the configuration file and the registry.
const Name = "aws:firewall"

const timestampFmt = "2006-01-02T15:04:05.999999-0700"

const (
	EventTypeAlert   = "alert"
	EventTypeNetflow = "netflow"

	ProtocolICMP = "ICMP"
	ProtocolTCP  = "TCP"
	ProtocolUDP  = "UDP"

	AppProtoNone = ""
	AppProtoHTTP = "http"

	AlertActionAllowed = "allowed"
	AlertActionBlocked = "blocked"
)

var (
	eventTypes   = [...]string{EventTypeAlert, EventTypeNetflow}
	protocols    = [...]string{ProtocolICMP, ProtocolTCP, ProtocolUDP}
	tcpAppProtos = [...]string{AppProtoNone, AppProtoHTTP}
	alertActions = [...]string{AlertActionAllowed, AlertActionBlocked}
)

// HTTPData provides fields for HTTP records.
type HTTPData struct {
	Hostname      string `json:"hostname"`
	URL           string `json:"URL"`
	HTTPUserAgent string `json:"http_user_agent"`
	HTTPMethod    string `json:"http_method"`
	Protocol      string `json:"protocol"`
	Length        int    `json:"length"`
}

// AlertData provides fields for alert event records.
type AlertData struct {
	Action      string `json:"action"`
	SignatureID int    `json:"signature_id"`
	Rev         int    `json:"rev"`
	Signature   string `json:"signature"`
	Category    string `json:"category"`
	Severity    int    `json:"severity"`
}

// NetflowData provides fields for netflow event records
type NetflowData struct {
	Pkts   int    `json:"pkts"`
	Bytes  int    `json:"bytes"`
	Start  string `json:"start"`
	End    string `json:"end"`
	Age    int    `json:"age"`
	MinTTL int    `json:"min_ttl"`
	MaxTTL int    `json:"max_ttl"`
}

// TCPData provides fields for TCP-based records.
type TCPData struct {
	TCPFlags string `json:"tcp_flags"`
	Fin      bool   `json:"fin,omitempty"`
	Syn      bool   `json:"syn,omitempty"`
	Rst      bool   `json:"rst,omitempty"`
	Psh      bool   `json:"psh,omitempty"`
	Ack      bool   `json:"ack,omitempty"`
	Urg      bool   `json:"urg,omitempty"`
}

type EventData struct {
	Timestamp string       `json:"timestamp"`
	FlowID    int          `json:"flow_id"`
	EventType string       `json:"event_type"`
	SrcIP     net.IP       `json:"src_ip"`
	SrcPort   int          `json:"src_port"`
	DstIP     net.IP       `json:"dst_ip"`
	DstPort   int          `json:"dst_port"`
	Proto     string       `json:"proto"`
	AppProto  string       `json:"app_proto,omitempty"`
	Alert     *AlertData   `json:"alert,omitempty"`
	Netflow   *NetflowData `json:"netflow,omitempty"`
	HTTP      *HTTPData    `json:"http,omitempty"`
	TCP       *TCPData     `json:"tcp,omitempty"`
}

// Firewall holds the random fields for a firewall record.
type Firewall struct {
	FirewallName     string    `json:"firewall_name"`
	AvailabilityZone string    `json:"availability_zone"`
	Event            EventData `json:"event"`
	EventTimestamp   string    `json:"event_timestamp"`
}

// Generator provides an AWS Firewall record generator.
type Generator struct {
	Data Firewall

	eventType string
}

func init() {
	_ = generator.Register(Name, New)
}

// New is the factory for AWS Firewall objects.
func New(cfg *ucfg.Config, _ int64) (generator.Generator, error) {
	c := defaultConfig()
	if err := cfg.Unpack(&c); err != nil {
		return nil, err
	}

	g := Generator{
		eventType: c.EventType,
	}

	return &g, nil
}

// Next produces the next AWS Firewall record.
func (g *Generator) Next() ([]byte, error) {
	var err error

	g.randomize()

	data, err := json.Marshal(&g.Data)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal %s data: %w", Name, err)
	}

	return data, nil
}

func (g *Generator) randomize() {
	now := time.Now()
	g.Data = Firewall{
		FirewallName:     fmt.Sprintf("Firewall-%d", rand.Intn(100)),
		AvailabilityZone: random.AWSAvailabilityZone(),
		EventTimestamp:   strconv.Itoa(int(now.Unix())),
		Event: EventData{
			Timestamp: now.Format(timestampFmt),
			FlowID:    rand.Int(),
			SrcIP:     random.IPv4(),
			SrcPort:   random.Port(),
			DstIP:     random.IPv4(),
			DstPort:   random.Port(),
			Proto:     protocols[rand.Intn(len(protocols))],
		},
	}

	if g.eventType == "" {
		g.Data.Event.EventType = eventTypes[rand.Intn(len(eventTypes))]
	} else {
		g.Data.Event.EventType = g.eventType
	}

	switch g.Data.Event.EventType {
	case EventTypeAlert:
		g.randomizeAlert()
	case EventTypeNetflow:
		g.randomizeNetflow(now)
	}
}

func (g *Generator) randomizeAlert() {
	signature := rand.Intn(1024)
	g.Data.Event.Alert = &AlertData{
		Action:      alertActions[rand.Intn(len(alertActions))],
		SignatureID: signature,
		Rev:         rand.Intn(1024),
		Signature:   fmt.Sprintf("Signature-%d", signature),
		Category:    fmt.Sprintf("Category-%d", rand.Intn(100)),
		Severity:    rand.Intn(6),
	}

	if g.Data.Event.Proto == ProtocolTCP {
		g.randomizeTCP()
	}
}

func (g *Generator) randomizeNetflow(now time.Time) {
	ttl := rand.Intn(256)
	start := now.Add(-time.Duration(rand.Intn(60)) * time.Minute)
	g.Data.Event.Netflow = &NetflowData{
		Pkts:   rand.Intn(100),
		Start:  start.Format(timestampFmt),
		End:    now.Format(timestampFmt),
		Age:    int(now.Sub(start).Seconds()),
		MinTTL: ttl,
		MaxTTL: ttl,
	}
	g.Data.Event.Netflow.Bytes = g.Data.Event.Netflow.Pkts*rand.Intn(1024) + 1
}

func (g *Generator) randomizeTCP() {
	g.Data.Event.AppProto = tcpAppProtos[rand.Intn(len(tcpAppProtos))]

	flags := rand.Intn(64)
	g.Data.Event.TCP = &TCPData{
		TCPFlags: fmt.Sprintf("%02d", flags),
		Fin:      flags&(1<<0) != 0,
		Syn:      flags&(1<<1) != 0,
		Rst:      flags&(1<<2) != 0,
		Psh:      flags&(1<<3) != 0,
		Ack:      flags&(1<<4) != 0,
		Urg:      flags&(1<<5) != 0,
	}

	if g.Data.Event.AppProto == AppProtoHTTP {
		g.randomizeHTTP()
	}
}

func (g *Generator) randomizeHTTP() {
	g.Data.Event.HTTP = &HTTPData{
		Hostname:      fmt.Sprintf("HTTPHost-%d", rand.Intn(100)),
		URL:           fmt.Sprintf("/random-%d.html", rand.Intn(100)),
		HTTPUserAgent: random.UserAgent(),
		HTTPMethod:    random.HTTPMethod(),
		Protocol:      random.HTTPVersion(),
		Length:        rand.Intn(1024),
	}
}
