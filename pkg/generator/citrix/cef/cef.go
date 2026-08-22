// Package cef implements the generator for Citrix CEF logs.
//
//	generator:
//	  type: citrix:cef
package cef

import (
	"bytes"
	"encoding/hex"
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

// Details from https://docs.citrix.com/en-us/citrix-adc/downloads/cef-log-components.pdf,
// https://docs.citrix.com/en-us/citrix-adc/current-release/application-firewall/logs.html
// and https://support.citrix.com/article/CTX136146/common-event-format-cef-logging-support-in-the-application-firewall.

// Name is the name of the generator in the configuration file and registry
const Name = "citrix:cef"

var (
	tmpl         = `{{.Timestamp.Format .TimeLayout}} <{{.Facility}}.{{.Priority}}> {{.Addr}} CEF:{{.CEFVersion}}|{{.Vendor}}|{{.Product}}|{{.Version}}|{{.Module}}|{{.Violation}}|{{.Severity}}|src={{.SrcAddr}} {{with .Geo}}geolocation={{.}} {{end}}spt={{.SrcPort}} method={{.Method}} request={{.Request}} msg={{.Message}} cn1={{.EventID}} cn2={{.TxID}} cs1={{.Profile}} cs2={{.PPEID}} cs3={{.SessID}} cs4={{.SeverityLabel}} cs5={{.Year}} {{with .ViolationCategory}}cs6={{.}} {{end}}act={{.Action}}`
	msgTemplates = []string{
		tmpl,
	}
	timeLayouts = []string{
		"Jan 02 15:04:05",
		"Jan 2 15:04:05",
	}
	facilities = []string{
		"auth", "authpriv", "cron", "daemon", "kern", "lpr", "mail", "mark", "news", "syslog", "user", "uucp", "local0", "local1", "local2", "local3", "local4", "local5", "local6", "local7",
	}
	priorities = []string{
		"debug", "info", "notice", "warning", "warn", "err", "error", "crit", "alert", "emerg", "panic",
	}
	vendors = []string{
		"Citrix",
	}
	products = []string{
		"NetScalar",
	}
	versions = []string{
		"NS10.0",
		"NS11.0",
	}
	modules = []string{
		"APPFW",
	}
	violations = []string{
		"APPFW_FIELDCONSISTENCY",
		"APPFW_SAFECOMMERCE",
		"APPFW_SAFECOMMERCE_XFORM",
		"APPFW_SIGNATURE_MATCH",
		"APPFW_STARTURL",
	}
	locations = []string{
		"",
		"Unknown",
		"NorthAmerica.US.Arizona.Tucson.*.*",
		"NorthAmerica.US.Arizona.Phoenix.*.*",
		"NorthAmerica.US.California.SanFrancisco.*.*",
	}
	methods = []string{
		"GET", "POST",
	}
	requests = []string{
		`http://aaron.stratum8.net/FFC/login.html`,
		`http://aaron.stratum8.net/FFC/login.php?login_name=abc&passwd=123456789234&drinking_pref=on&text_area=&loginButton=ClickToLogin&as_sfid=AAAAAAWIahZuYoIFbjBhYMP05mJLTwEfIY0a7AKGMg3jIBaKmwtK4t7M7lNxOgj7Gmd3SZc8KUj6CR6a7W5kIWDRHN8PtK1Zc-txHkHNx1WknuG9DzTuM7t1THhluevXu9I4kp8%3D&as_fid=feeec8758b41740eedeeb6b35b85dfd3d5def30c`,
		`http://aaron.stratum8.net/FFC/wwwboard/passwd.txt`,
		`http://aaron.stratum8.net/FFC/CreditCardMind.html`,
		`http://vpx247.example.net/FFC/CreditCardMind.html`,
		`http://vpx247.example.net/FFC/login_post.html?abc\=def`,
		`http://vpx247.example.net/FFC/wwwboard/passwd.txt`,
	}
	messages = []string{
		"Signature violation rule ID 807: web-cgi /wwwboard/passwd.txt access",
		"Disallow Illegal URL.",
		"Transformed (xout) potential credit card numbers seen in server response",
		"Maximum number of potential credit card numbers seen",
		"Field consistency check failed for field passwd",
	}
	profiles = []string{
		"pr_ffc",
	}
	severityLabels = []string{
		"INFO", "ALERT",
	}
	violationCategory = []string{
		"",
		"web-cgi",
		"sql-injection",
		"phishing",
	}
	actions = []string{
		"blocked", "not blocked", "transformed",
	}
)

type CEF struct {
	Timestamp  time.Time
	TimeLayout string

	Facility string
	Priority string

	Addr net.IP

	CEFVersion int
	Vendor     string
	Product    string
	Version    string
	Module     string
	Violation  string
	Severity   int

	SrcAddr           net.IP
	Geo               string
	SrcPort           int
	Method            string
	Request           string
	Message           string
	EventID           int
	TxID              int
	Profile           string
	PPEID             string
	SessID            string
	SeverityLabel     string
	Year              int
	ViolationCategory string
	Action            string

	templates []*template.Template
}

func init() {
	generator.Register(Name, New)
}

// New returns a new Citrix CEF log line generator.
func New(cfg *ucfg.Config, _ int64) (generator.Generator, error) {
	def := defaultConfig()
	if err := cfg.Unpack(&def); err != nil {
		return nil, err
	}

	c := &CEF{}
	c.randomize()

	for i, v := range msgTemplates {
		t, err := template.New(strconv.Itoa(i)).Funcs(generator.FunctionMap).Parse(v)
		if err != nil {
			return nil, err
		}
		c.templates = append(c.templates, t)
	}

	return c, nil
}

// Next produces the next CEF log entry.
func (c *CEF) Next() ([]byte, error) {
	var buf bytes.Buffer

	err := c.templates[rand.Intn(len(c.templates))].Execute(&buf, c)
	if err != nil {
		return nil, err
	}

	c.randomize()

	return buf.Bytes(), err
}

func (c *CEF) randomize() {
	c.Timestamp = time.Now()
	c.TimeLayout = randString(timeLayouts)

	c.Facility = randString(facilities)
	c.Priority = randString(priorities)

	c.Addr = random.IPv4()

	c.CEFVersion = rand.Intn(2)
	c.Vendor = randString(vendors)
	c.Product = randString(products)
	c.Version = randString(versions)
	c.Module = randString(modules)
	c.Violation = randString(violations)
	c.Severity = rand.Intn(10) + 1

	c.SrcAddr = random.IPv4()
	c.Geo = randString(locations)
	c.SrcPort = random.Port()
	c.Method = randString(methods)
	c.Request = randString(requests)
	c.Message = randString(messages)
	c.EventID = rand.Intn(1000)
	c.TxID = rand.Intn(100000)
	c.Profile = randString(profiles)
	c.PPEID = fmt.Sprintf("PPE%d", rand.Intn(9)+1)
	sessID := make([]byte, 16)
	rand.Read(sessID)
	c.SessID = hex.EncodeToString(sessID)
	c.SeverityLabel = randString(severityLabels)
	c.Year = c.Timestamp.Year()
	c.ViolationCategory = randString(violationCategory)
	c.Action = randString(actions)
}

func randString(s []string) string {
	return s[rand.Intn(len(s))]
}
