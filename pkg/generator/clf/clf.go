// Package clf generates Common Log Format (clf) log messages.
//
// Configuration:
//
//	combined: (bool, optional) If true, generate Combined Log Format records,
//	          which add referer and user-agent fields.
//
//	- generator:
//	    type: clf
//	    combined: true
package clf

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

const (
	// Name is the name used in the configuration file and the registry.
	Name = "clf"

	commonTemplate   = "{{.Host}} {{.Ident}} {{.AuthUser}} {{.Date}} {{.Request}} {{.Status}} {{.Bytes}}"
	combinedTemplate = "{{.Host}} {{.Ident}} {{.AuthUser}} {{.Date}} {{.Request}} {{.Status}} {{.Bytes}} {{.Referer}} {{.UserAgent}}"
	timestampFmt     = "[02/Jan/2006:15:04:05 -0700]"
)

// Record holds the random fields for a Common Log Format record.
type Record struct {
	Host      net.IP
	Ident     string
	AuthUser  string
	Date      string
	Request   string
	Status    string
	Bytes     string
	Referer   string
	UserAgent string
}

// Generator provides a Common Log Format record generator.
type Generator struct {
	Record Record

	tmpl       *template.Template
	staticTime *time.Time
	buf        bytes.Buffer
	combined   bool
}

// Next produces the next Common Log Format record.
//
// Example:
//
// 127.0.0.1 - - [10/Oct/2000:13:55:36 -0700] "GET /random-100.html HTTP/1.0" 200 2326
func (g *Generator) Next() ([]byte, error) {
	g.randomize()

	g.buf.Reset()
	if err := g.tmpl.Execute(&g.buf, &g.Record); err != nil {
		return nil, err
	}

	return g.buf.Bytes(), nil
}

func (g *Generator) randomize() {
	now := time.Now()
	if g.staticTime != nil {
		now = *g.staticTime
	}

	g.Record = Record{
		Host:     random.IPv4(),
		Ident:    "-",
		AuthUser: "-",
		Date:     now.Format(timestampFmt),
		Request: fmt.Sprintf(
			`"%s %s %s"`,
			random.HTTPMethod(),
			fmt.Sprintf("/random-%d.html", rand.Intn(100)),
			random.HTTPVersion(),
		),
		Status: strconv.Itoa(random.HTTPStatus()),
		Bytes:  strconv.Itoa(rand.Intn(10000)),
	}
	if g.combined {
		g.Record.Referer = "-"
		g.Record.UserAgent = `"` + random.UserAgent() + `"`
	}
}

// New is the factory for Common Log Format objects.
func New(cfg *ucfg.Config, _ int64) (generator.Generator, error) {
	var err error

	c := defaultConfig()
	if err := cfg.Unpack(&c); err != nil {
		return nil, err
	}

	g := Generator{
		combined: c.Combined,
	}

	if g.combined {
		g.tmpl, err = template.New("clf").Funcs(generator.FunctionMap).Parse(combinedTemplate)
	} else {
		g.tmpl, err = template.New("clf").Funcs(generator.FunctionMap).Parse(commonTemplate)
	}
	if err != nil {
		return nil, err
	}

	return &g, nil
}

func init() {
	if err := generator.Register(Name, New); err != nil {
		panic(err)
	}
}
