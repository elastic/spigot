// Package cef implements the generator for generic CEF logs.
//
// Configuration file supports including timestamps in log messages
//
//	  generator:
//	    type: "generic:cef"
//	    max_extensions: 20
//		   vendors: ["VaporCorp", ...]
//		   products: ["VaporWare"]
//		   versions: ["0.1", "0.1-alpha"]
//		   classes: ["APPSS"]
//		   names: ["APPSS_UL", "APPSS_LTUL"]
//	    must_include: ["src", "spt", "dst", "dpt",...]
//	    must_exclude: ["art",...]
package cef

import (
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/elastic/go-ucfg"
	"github.com/elastic/spigot/pkg/generator"
	"github.com/elastic/spigot/pkg/random"
	"github.com/google/uuid"
)

// Name is the name of the generator in the configuration file and registry
const Name = "generic:cef"

var (
	tmpl         = `CEF:{{.CEFVersion}}|{{.Vendor}}|{{.Product}}|{{.Version}}|{{.Class}}|{{.Name}}|{{.Severity}}|{{range $i, $v := .Extensions}}{{if $i}} {{end}}{{$v}}{{end}}`
	msgTemplates = []string{
		tmpl,
	}
)

type CEF struct {
	CEFVersion int
	Vendor     string
	Product    string
	Version    string
	Class      string
	Name       string
	Severity   int

	Extensions []string

	config
	templates []*template.Template
}

func init() {
	generator.Register(Name, New)
}

// New returns a new CEF log line generator.
func New(cfg *ucfg.Config, _ int64) (generator.Generator, error) {
	config := defaultConfig()
	if err := cfg.Unpack(&config); err != nil {
		return nil, err
	}

	c := &CEF{config: config}
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
	c.CEFVersion = randInt(c.CEFVersions)
	c.Vendor = randString(c.Vendors)
	c.Product = randString(c.Products)
	c.Version = randString(c.Versions)
	c.Class = randString(c.Classes)
	c.Name = randString(c.Names)
	c.Severity = randInt(c.Severities)

	c.Extensions = c.Extensions[:0]
	if c.Max == 0 {
		return
	}
	have := make(map[string]bool)
	for _, x := range c.Exclude {
		have[x] = true
	}
	for _, m := range c.Must {
		c.addExtension(m, have)
	}
	perm := rand.Perm(len(extensions))
	max := rand.Intn(c.Max)
	for _, p := range perm {
		if len(c.Extensions) >= max {
			break
		}
		c.addExtension(extensions[p], have)
	}
	rand.Shuffle(len(c.Extensions), func(i, j int) { c.Extensions[i], c.Extensions[j] = c.Extensions[j], c.Extensions[i] })
}

func (c *CEF) addExtension(abbrev string, have map[string]bool) {
	cand := extensionMapping[abbrev]
	if cand.Wants == "" && !have[cand.Abbrev] {
		c.Extensions = append(c.Extensions, cand.Render(c.config))
		have[cand.Abbrev] = true
		return
	}
	var add []mappedField
	seen := make(map[string]bool)
	for {
		if seen[cand.Abbrev] {
			panic(fmt.Sprintf("cycle including %s", cand.Abbrev))
		}
		seen[cand.Abbrev] = true
		add = append(add, cand)
		if cand.Wants == "" {
			if len(add) > c.Max {
				add = add[len(add)-c.Max:]
			}
			for _, a := range add {
				if have[a.Abbrev] {
					continue
				}
				c.Extensions = append(c.Extensions, a.Render(c.config))
				have[a.Abbrev] = true
			}
			break
		}
		cand = extensionMapping[cand.Wants]
	}
}

func randInt(i []int) int {
	return i[rand.Intn(len(i))]
}

func randString(s []string) string {
	return s[rand.Intn(len(s))]
}

type mappedField struct {
	Abbrev string
	Target string
	Value  func(config) fmt.Stringer
	Wants  string
}

func (f mappedField) Render(c config) string {
	if f.Value == nil {
		return ""
	}
	return fmt.Sprintf("%s=%s", f.Abbrev, f.Value(c))
}

func init() {
	for k, v := range extensionMapping {
		if v.Value == nil {
			continue
		}

		v.Abbrev = k
		extensionMapping[k] = v

		extensions = append(extensions, k)
		if strings.HasSuffix(k, "Label") {
			_, ok := extensionMapping[v.Wants]
			if !ok {
				panic(fmt.Sprintf("label mapping without target: %v", v.Target))
			}
		}
	}
	sort.Strings(extensions)
}

var extensions []string

// extensionMapping is a mapping of CEF key names to full field names and data
// types. This mapping was generated from tables contained in:
//   - "Micro Focus Security ArcSight Common Event Format Version 25"
//     dated September 28, 2017.
//   - "Check Point Log Exporter CEF Field Mappings"
//     dated November 23, 2018.
//   - "HPE Security ArcSight Common Event Format Version 23"
//     dated May 16, 2016.
var extensionMapping = map[string]mappedField{
	"agt": {
		Target: "agentAddress",
		Value:  func(c config) fmt.Stringer { return ipv4Value{} },
	},
	"agentDnsDomain": {
		Target: "agentDnsDomain",
		Value:  func(c config) fmt.Stringer { return domainValue(c.Words) },
	},
	"ahost": {
		Target: "agentHostName",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"aid": {
		Target: "agentId",
		Value:  func(c config) fmt.Stringer { return uuidValue{zero: c.ZeroUUID} },
	},
	"amac": {
		Target: "agentMacAddress",
		Value:  func(c config) fmt.Stringer { return hwaddrValue{6} },
	},
	"agentNtDomain": {
		Target: "agentNtDomain",
		Value:  func(c config) fmt.Stringer { return domainValue(c.Words) },
	},
	"art": {
		Target: "agentReceiptTime",
		Value: func(c config) fmt.Stringer {
			return integerValue{int(c.Now().Add(-time.Hour).UnixMilli()), int(c.Now().UnixMilli())}
		},
	},
	"atz": {
		Target: "agentTimeZone",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.TimeZones) },
	},
	"agentTranslatedAddress": {
		Target: "agentTranslatedAddress",
		Value:  func(c config) fmt.Stringer { return ipv4Value{} },
	},
	"agentTranslatedZoneExternalID": {
		Target: "agentTranslatedZoneExternalID",
		Value:  func(c config) fmt.Stringer { return uuidValue{zero: c.ZeroUUID} },
	},
	"agentTranslatedZoneURI": {
		Target: "agentTranslatedZoneURI",
		Value:  func(c config) fmt.Stringer { return urlValue(c.Words) },
	},
	"at": {
		Target: "agentType",
		Value:  func(c config) fmt.Stringer { return keywordValue{"local", "network"} },
	},
	"av": {
		Target: "agentVersion",
		Value:  func(c config) fmt.Stringer { return integerValue{0, 5} },
	},
	"agentZoneExternalID": {
		Target: "agentZoneExternalID",
		Value:  func(c config) fmt.Stringer { return uuidValue{zero: c.ZeroUUID} },
	},
	"agentZoneURI": {
		Target: "agentZoneURI",
		Value:  func(c config) fmt.Stringer { return urlValue(c.Words) },
	},
	"app": {
		Target: "applicationProtocol",
		Value: func(c config) fmt.Stringer {
			return keywordValue{"tcp", "TCP", "udp", "UDP", "sip", "SIP", "http", "HTTP"}
		},
	},
	"cnt": {
		Target: "baseEventCount",
		Value:  func(c config) fmt.Stringer { return integerValue{0, 1e3} },
	},
	"in": {
		Target: "bytesIn",
		Value:  func(c config) fmt.Stringer { return integerValue{0, 1e5} },
	},
	"out": {
		Target: "bytesOut",
		Value:  func(c config) fmt.Stringer { return integerValue{0, 1e5} },
	},
	"customerExternalID": {
		Target: "customerExternalID",
		Value:  func(c config) fmt.Stringer { return uuidValue{zero: c.ZeroUUID} },
	},
	"customerURI": {
		Target: "customerURI",
		Value:  func(c config) fmt.Stringer { return urlValue(c.Words) },
	},
	"dst": {
		Target: "destinationAddress",
		Value:  func(c config) fmt.Stringer { return ipv4Value{} },
	},
	"destinationDnsDomain": {
		Target: "destinationDnsDomain",
		Value:  func(c config) fmt.Stringer { return domainValue(c.Words) },
	},
	"dlat": {
		Target: "destinationGeoLatitude",
		Value:  func(c config) fmt.Stringer { return floatValue{-180, 180} },
	},
	"dlong": {
		Target: "destinationGeoLongitude",
		Value:  func(c config) fmt.Stringer { return floatValue{-90, 90} },
	},
	"dhost": {
		Target: "destinationHostName",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"dmac": {
		Target: "destinationMacAddress",
		Value:  func(c config) fmt.Stringer { return hwaddrValue{6} },
	},
	"dntdom": {
		Target: "destinationNtDomain",
		Value:  func(c config) fmt.Stringer { return domainValue(c.Words) },
	},
	"dpt": {
		Target: "destinationPort",
		Wants:  "dst",
		Value:  func(c config) fmt.Stringer { return integerValue{0, 65535} },
	},
	"dpid": {
		Target: "destinationProcessId",
		Value:  func(c config) fmt.Stringer { return integerValue{0, 65535} },
	},
	"dproc": {
		Target: "destinationProcessName",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"destinationServiceName": {
		Target: "destinationServiceName",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"destinationTranslatedAddress": {
		Target: "destinationTranslatedAddress",
		Value:  func(c config) fmt.Stringer { return ipv4Value{} },
	},
	"destinationTranslatedPort": {
		Target: "destinationTranslatedPort",
		Value:  func(c config) fmt.Stringer { return integerValue{0, 65535} },
	},
	"destinationTranslatedZoneExternalID": {
		Target: "destinationTranslatedZoneExternalID",
		Value:  func(c config) fmt.Stringer { return uuidValue{zero: c.ZeroUUID} },
	},
	"destinationTranslatedZoneURI": {
		Target: "destinationTranslatedZoneURI",
		Value:  func(c config) fmt.Stringer { return urlValue(c.Words) },
	},
	"duid": {
		Target: "destinationUserId",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Users) },
	},
	"duser": {
		Target: "destinationUserName",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Users) },
	},
	"dpriv": {
		Target: "destinationUserPrivileges",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Privs) },
	},
	"destinationZoneExternalID": {
		Target: "destinationZoneExternalID",
		Value:  func(c config) fmt.Stringer { return uuidValue{zero: c.ZeroUUID} },
	},
	"destinationZoneURI": {
		Target: "destinationZoneURI",
		Value:  func(c config) fmt.Stringer { return urlValue(c.Words) },
	},
	"act": {
		Target: "deviceAction",
		Value:  func(c config) fmt.Stringer { return keywordValue(actions) },
	},
	"dvc": {
		Target: "deviceAddress",
		Value:  func(c config) fmt.Stringer { return ipv4Value{} },
	},
	"cfp1": {
		Target: "deviceCustomFloatingPoint1",
		Value:  func(c config) fmt.Stringer { return floatValue{0, 100} },
	},
	"cfp1Label": {
		Target: "deviceCustomFloatingPoint1Label",
		Wants:  "cfp1",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"cfp2": {
		Target: "deviceCustomFloatingPoint2",
		Value:  func(c config) fmt.Stringer { return floatValue{0, 100} },
	},
	"cfp2Label": {
		Target: "deviceCustomFloatingPoint2Label",
		Wants:  "cfp2",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"cfp3": {
		Target: "deviceCustomFloatingPoint3",
		Value:  func(c config) fmt.Stringer { return floatValue{0, 100} },
	},
	"cfp3Label": {
		Target: "deviceCustomFloatingPoint3Label",
		Wants:  "cfp3",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"cfp4": {
		Target: "deviceCustomFloatingPoint4",
		Value:  func(c config) fmt.Stringer { return floatValue{0, 100} },
	},
	"cfp4Label": {
		Target: "deviceCustomFloatingPoint4Label",
		Wants:  "cfp4",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"deviceCustomDate1": {
		Target: "deviceCustomDate1",
		Value: func(c config) fmt.Stringer {
			return integerValue{int(c.Now().Add(-time.Hour).UnixMilli()), int(c.Now().UnixMilli())}
		},
	},
	"deviceCustomDate1Label": {
		Target: "deviceCustomDate1Label",
		Wants:  "deviceCustomDate1",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"deviceCustomDate2": {
		Target: "deviceCustomDate2",
		Value: func(c config) fmt.Stringer {
			return integerValue{int(c.Now().Add(-time.Hour).UnixMilli()), int(c.Now().UnixMilli())}
		},
	},
	"deviceCustomDate2Label": {
		Target: "deviceCustomDate2Label",
		Wants:  "deviceCustomDate2",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"c6a1": {
		Target: "deviceCustomIPv6Address1",
		Value:  func(c config) fmt.Stringer { return ipv6Value{} },
	},
	"c6a1Label": {
		Target: "deviceCustomIPv6Address1Label",
		Wants:  "c6a1",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"c6a2": {
		Target: "deviceCustomIPv6Address2",
		Value:  func(c config) fmt.Stringer { return ipv6Value{} },
	},
	"c6a2Label": {
		Target: "deviceCustomIPv6Address2Label",
		Wants:  "c6a2",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"c6a3": {
		Target: "deviceCustomIPv6Address3",
		Value:  func(c config) fmt.Stringer { return ipv6Value{} },
	},
	"c6a3Label": {
		Target: "deviceCustomIPv6Address3Label",
		Wants:  "c6a3",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"c6a4": {
		Target: "deviceCustomIPv6Address4",
		Value:  func(c config) fmt.Stringer { return ipv6Value{} },
	},
	"C6a4Label": {
		Target: "deviceCustomIPv6Address4Label",
		Wants:  "c6a4",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"cn1": {
		Target: "deviceCustomNumber1",
		Value:  func(c config) fmt.Stringer { return integerValue{0, 1000} },
	},
	"cn1Label": {
		Target: "deviceCustomNumber1Label",
		Wants:  "cn1",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"cn2": {
		Target: "deviceCustomNumber2",
		Value:  func(c config) fmt.Stringer { return integerValue{0, 1000} },
	},
	"cn2Label": {
		Target: "deviceCustomNumber2Label",
		Wants:  "cn2",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"cn3": {
		Target: "deviceCustomNumber3",
		Value:  func(c config) fmt.Stringer { return integerValue{0, 1000} },
	},
	"cn3Label": {
		Target: "deviceCustomNumber3Label",
		Wants:  "cn3",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"cs1": {
		Target: "deviceCustomString1",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"cs1Label": {
		Target: "deviceCustomString1Label",
		Wants:  "cs1",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"cs2": {
		Target: "deviceCustomString2",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"cs2Label": {
		Target: "deviceCustomString2Label",
		Wants:  "cs2",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"cs3": {
		Target: "deviceCustomString3",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"cs3Label": {
		Target: "deviceCustomString3Label",
		Wants:  "cs3",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"cs4": {
		Target: "deviceCustomString4",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"cs4Label": {
		Target: "deviceCustomString4Label",
		Wants:  "cs4",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"cs5": {
		Target: "deviceCustomString5",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"cs5Label": {
		Target: "deviceCustomString5Label",
		Wants:  "cs5",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"cs6": {
		Target: "deviceCustomString6",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"cs6Label": {
		Target: "deviceCustomString6Label",
		Wants:  "cs6",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"deviceDirection": {
		Target: "deviceDirection",
		Value:  func(c config) fmt.Stringer { return integerValue{0, 1} },
	},
	"deviceDnsDomain": {
		Target: "deviceDnsDomain",
		Value:  func(c config) fmt.Stringer { return domainValue(c.Words) },
	},
	"cat": {
		Target: "deviceEventCategory",
	},
	"deviceExternalId": {
		Target: "deviceExternalId",
		Value:  func(c config) fmt.Stringer { return uuidValue{zero: c.ZeroUUID} },
	},
	"deviceFacility": {
		Target: "deviceFacility",
		Value: func(c config) fmt.Stringer {
			return keywordValue{"auth", "authpriv", "cron", "daemon", "kern", "lpr", "mail", "mark", "news", "syslog", "user", "uucp", "local0", "local1", "local2", "local3", "local4", "local5", "local6", "local7"}
		},
	},
	"dvchost": {
		Target: "deviceHostName",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"deviceInboundInterface": {
		Target: "deviceInboundInterface",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Interfaces) },
	},
	"dvcmac": {
		Target: "deviceMacAddress",
		Value:  func(c config) fmt.Stringer { return hwaddrValue{6} },
	},
	"deviceNtDomain": {
		Target: "deviceNtDomain",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"DeviceOutboundInterface": {
		Target: "deviceOutboundInterface",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Interfaces) },
	},
	"DevicePayloadId": {
		Target: "devicePayloadId",
		Value:  func(c config) fmt.Stringer { return uuidValue{zero: c.ZeroUUID} },
	},
	"dvcpid": {
		Target: "deviceProcessId",
		Value:  func(c config) fmt.Stringer { return integerValue{0, 65535} },
	},
	"deviceProcessName": {
		Target: "deviceProcessName",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"rt": {
		Target: "deviceReceiptTime",
		Value: func(c config) fmt.Stringer {
			return integerValue{int(c.Now().Add(-time.Hour).UnixMilli()), int(c.Now().UnixMilli())}
		},
	},
	"dtz": {
		Target: "deviceTimeZone",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.TimeZones) },
		Wants:  "rt",
	},
	"deviceTranslatedAddress": {
		Target: "deviceTranslatedAddress",
		Value:  func(c config) fmt.Stringer { return ipv4Value{} },
	},
	"deviceTranslatedZoneExternalID": {
		Target: "deviceTranslatedZoneExternalID",
		Value:  func(c config) fmt.Stringer { return uuidValue{zero: c.ZeroUUID} },
	},
	"deviceTranslatedZoneURI": {
		Target: "deviceTranslatedZoneURI",
		Value:  func(c config) fmt.Stringer { return urlValue(c.Words) },
	},
	"deviceZoneExternalID": {
		Target: "deviceZoneExternalID",
		Value:  func(c config) fmt.Stringer { return uuidValue{zero: c.ZeroUUID} },
	},
	"deviceZoneURI": {
		Target: "deviceZoneURI",
		Value:  func(c config) fmt.Stringer { return urlValue(c.Words) },
	},
	"end": {
		Target: "endTime",
		Value: func(c config) fmt.Stringer {
			return integerValue{int(c.Now().Add(-time.Hour).UnixMilli()), int(c.Now().UnixMilli())}
		},
	},
	"eventId": {
		Target: "eventId",
		Value:  func(c config) fmt.Stringer { return integerValue{0, 1e5} },
	},
	"outcome": {
		Target: "eventOutcome",
		Value:  func(c config) fmt.Stringer { return keywordValue{"success", "failure"} },
	},
	"externalId": {
		Target: "externalId",
		Value:  func(c config) fmt.Stringer { return uuidValue{zero: c.ZeroUUID} },
	},
	"fileCreateTime": {
		Target: "fileCreateTime",
		Value: func(c config) fmt.Stringer {
			return integerValue{int(c.Now().Add(-time.Hour).UnixMilli()), int(c.Now().UnixMilli())}
		},
	},
	"fileHash": {
		Target: "fileHash",
		Value:  func(c config) fmt.Stringer { return hashValue{16} },
	},
	"fileId": {
		Target: "fileId",
		Value:  func(c config) fmt.Stringer { return uuidValue{zero: c.ZeroUUID} },
	},
	"fileModificationTime": {
		Target: "fileModificationTime",
		Value: func(c config) fmt.Stringer {
			return integerValue{int(c.Now().Add(-time.Hour).UnixMilli()), int(c.Now().UnixMilli())}
		},
	},
	"flexNumber1": {
		Target: "deviceFlexNumber1",
		Value:  func(c config) fmt.Stringer { return integerValue{0, 1e5} },
	},
	"flexNumber1Label": {
		Target: "deviceFlexNumber1Label",
		Wants:  "flexNumber1",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"flexNumber2": {
		Target: "deviceFlexNumber2",
		Value:  func(c config) fmt.Stringer { return integerValue{0, 1e5} },
	},
	"flexNumber2Label": {
		Target: "deviceFlexNumber2Label",
		Wants:  "flexNumber2",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},

	"fname": {
		Target: "filename",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"filePath": {
		Target: "filePath",
	},
	"filePermission": {
		Target: "filePermission",
	},
	"fsize": {
		Target: "fileSize",
		Value:  func(c config) fmt.Stringer { return integerValue{0, 1e5} },
	},
	"fileType": {
		Target: "fileType",
		Value:  func(c config) fmt.Stringer { return keywordValue{"directory", "regular", "pipe", "socket"} },
	},
	"flexDate1": {
		Target: "flexDate1",
		Value: func(c config) fmt.Stringer {
			return integerValue{int(c.Now().Add(-time.Hour).UnixMilli()), int(c.Now().UnixMilli())}
		},
	},
	"flexDate1Label": {
		Target: "flexDate1Label",
		Wants:  "flexDate1",
	},
	"flexString1": {
		Target: "flexString1",
	},
	"flexString2": {
		Target: "flexString2",
	},
	"flexString1Label": {
		Target: "flexString1Label",
		Wants:  "flexString1",
	},
	"flexString2Label": {
		Target: "flexString2Label",
		Wants:  "flexString2",
	},
	"msg": {
		Target: "message",
		Value:  func(c config) fmt.Stringer { return keywordValue(messages) },
	},
	"oldFileCreateTime": {
		Target: "oldFileCreateTime",
		Value: func(c config) fmt.Stringer {
			return integerValue{int(c.Now().Add(-time.Hour).UnixMilli()), int(c.Now().UnixMilli())}
		},
	},
	"oldFileHash": {
		Target: "oldFileHash",
		Value:  func(c config) fmt.Stringer { return hashValue{16} },
	},
	"oldFileId": {
		Target: "oldFileId",
	},
	"oldFileModificationTime": {
		Target: "oldFileModificationTime",
		Value: func(c config) fmt.Stringer {
			return integerValue{int(c.Now().Add(-time.Hour).UnixMilli()), int(c.Now().UnixMilli())}
		},
	},
	"oldFileName": {
		Target: "oldFileName",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"oldFilePath": {
		Target: "oldFilePath",
	},
	"oldFilePermission": {
		Target: "oldFilePermission",
	},
	"oldFileSize": {
		Target: "oldFileSize",
		Value:  func(c config) fmt.Stringer { return integerValue{0, 1e5} },
	},
	"oldFileType": {
		Target: "oldFileType",
		Value:  func(c config) fmt.Stringer { return keywordValue{"directory", "regular", "pipe", "socket"} },
	},
	"rawEvent": {
		Target: "rawEvent",
		Value:  func(c config) fmt.Stringer { return textValue{1, 500} },
	},
	"reason": {
		Target: "Reason",
		Value:  func(c config) fmt.Stringer { return keywordValue{"bad password", "unknown user", "banned"} },
	},
	"requestClientApplication": {
		Target: "requestClientApplication",
		Value:  func(c config) fmt.Stringer { return stringerise(random.UserAgent) },
	},
	"requestContext": {
		Target: "requestContext",
		Value:  func(c config) fmt.Stringer { return urlValue(c.Words) },
	},
	"requestCookies": {
		Target: "requestCookies",
		Value:  func(c config) fmt.Stringer { return uuidValue{zero: c.ZeroUUID} },
	},
	"requestMethod": {
		Target: "requestMethod",
		Value: func(c config) fmt.Stringer {
			return keywordValue{http.MethodConnect, http.MethodDelete, http.MethodGet, http.MethodPost, http.MethodPut}
		},
	},
	"request": {
		Target: "requestUrl",
		Value:  func(c config) fmt.Stringer { return urlValue(c.Words) },
	},
	"src": {
		Target: "sourceAddress",
		Value:  func(c config) fmt.Stringer { return ipv4Value{} },
	},
	"sourceDnsDomain": {
		Target: "sourceDnsDomain",
		Value:  func(c config) fmt.Stringer { return domainValue(c.Words) },
	},
	"slat": {
		Target: "sourceGeoLatitude",
		Value:  func(c config) fmt.Stringer { return floatValue{-180, 180} },
	},
	"slong": {
		Target: "sourceGeoLongitude",
		Value:  func(c config) fmt.Stringer { return floatValue{-90, 90} },
	},
	"shost": {
		Target: "sourceHostName",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"smac": {
		Target: "sourceMacAddress",
		Value:  func(c config) fmt.Stringer { return hwaddrValue{6} },
	},
	"sntdom": {
		Target: "sourceNtDomain",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"spt": {
		Target: "sourcePort",
		Value:  func(c config) fmt.Stringer { return integerValue{min: 0, max: 65535} },
	},
	"spid": {
		Target: "sourceProcessId",
		Value:  func(c config) fmt.Stringer { return integerValue{min: 0, max: 65535} },
	},
	"sproc": {
		Target: "sourceProcessName",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"sourceServiceName": {
		Target: "sourceServiceName",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Words) },
	},
	"sourceTranslatedAddress": {
		Target: "sourceTranslatedAddress",
		Value:  func(c config) fmt.Stringer { return ipv4Value{} },
	},
	"sourceTranslatedPort": {
		Target: "sourceTranslatedPort",
		Value:  func(c config) fmt.Stringer { return integerValue{min: 0, max: 65535} },
	},
	"sourceTranslatedZoneExternalID": {
		Target: "sourceTranslatedZoneExternalID",
		Value:  func(c config) fmt.Stringer { return uuidValue{zero: c.ZeroUUID} },
	},
	"sourceTranslatedZoneURI": {
		Target: "sourceTranslatedZoneURI",
		Value:  func(c config) fmt.Stringer { return urlValue(c.Words) },
	},
	"suid": {
		Target: "sourceUserId",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Users) },
	},
	"suser": {
		Target: "sourceUserName",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Users) },
	},
	"spriv": {
		Target: "sourceUserPrivileges",
		Value:  func(c config) fmt.Stringer { return keywordValue(c.Privs) },
	},
	"sourceZoneExternalID": {
		Target: "sourceZoneExternalID",
		Value:  func(c config) fmt.Stringer { return uuidValue{zero: c.ZeroUUID} },
	},
	"sourceZoneURI": {
		Target: "sourceZoneURI",
		Value:  func(c config) fmt.Stringer { return urlValue(c.Words) },
	},
	"start": {
		Target: "startTime",
		Value: func(c config) fmt.Stringer {
			return integerValue{int(c.Now().Add(-time.Hour).UnixMilli()), int(c.Now().UnixMilli())}
		},
	},
	"proto": {
		Target: "transportProtocol",
		Value:  func(c config) fmt.Stringer { return keywordValue{"tcp", "TCP", "udp", "UDP"} },
	},
	"type": {
		Target: "type",
		Value:  func(c config) fmt.Stringer { return integerValue{0, 15} },
	},

	// This is an ArcSight categorization field that is commonly used, but its
	// short name is not contained in the documentation used for the above list.
	"catdt": {
		Target: "categoryDeviceType",
		Value:  func(c config) fmt.Stringer { return keywordValue{"Operating system", "Network-based IDS/IPS"} },
	},
	"mrt": {
		Target: "managerReceiptTime",
		Value: func(c config) fmt.Stringer {
			return integerValue{int(c.Now().Add(-time.Hour).UnixMilli()), int(c.Now().UnixMilli())}
		},
	},
}

type urlValue []string

func (u urlValue) String() string {
	return fmt.Sprintf("%s://%s/%s%s%s",
		keywordValue{"http", "https"}, domainValue(u), keywordValue(u), keywordValue{"/", "?"}, keywordValue(u),
	)
}

type domainValue []string

func (d domainValue) String() string {
	return fmt.Sprintf("%s.%s.%s", keywordValue(d), keywordValue(d), keywordValue{"com", "org", "co"})
}

type keywordValue []string

func (k keywordValue) String() string {
	return k[rand.Intn(len(k))]
}

type uuidValue struct {
	zero bool
}

func (u uuidValue) String() string {
	if u.zero {
		uuid, _ := uuid.NewRandomFromReader(bytes.NewReader(make([]byte, 16)))
		return uuid.String()
	}
	return uuid.New().String()
}

type hashValue struct {
	bytes int
}

func (h hashValue) String() string {
	buf := make([]byte, h.bytes)
	rand.Read(buf)
	return fmt.Sprintf("%0*x", h.bytes, buf)
}

type hwaddrValue struct {
	bytes int
}

func (a hwaddrValue) String() string {
	buf := make(net.HardwareAddr, a.bytes)
	rand.Read(buf)
	return buf.String()
}

type ipv4Value struct{}

func (ipv4Value) String() string {
	buf := make(net.IP, 4)
	rand.Read(buf)
	return buf.String()
}

type ipv6Value struct{}

func (ipv6Value) String() string {
	buf := make(net.IP, 16)
	rand.Read(buf)
	for i := range buf {
		if rand.Float64() < 0.3 {
			buf[i] = 0
		}
	}
	return buf.String()
}

type timeValue struct {
	min, max int64
	format   string
}

func newTimeValue(min, max time.Time, format string) timeValue {
	return timeValue{min: min.UnixMilli(), max: max.UnixMilli(), format: format}
}

func (t timeValue) String() string {
	return time.UnixMilli(rand.Int63n(t.max-t.min) + t.min).Format(t.format)
}

type integerValue struct {
	min, max int
}

func (t integerValue) String() string {
	return strconv.Itoa(rand.Intn(t.max-t.min+1) + t.min)
}

type floatValue struct {
	min, max float64
}

func (t floatValue) String() string {
	return strconv.FormatFloat(((t.max-t.min)*rand.Float64())+t.min, 'f', -1, 64)
}

type textValue struct {
	min, max int
}

func (t textValue) String() string {
	idx := rand.Intn(t.max-t.min) + t.min
	words := strings.Split(loremIpsum, " ")
	if idx >= len(words) {
		return loremIpsum
	}
	return strings.Join(words[:idx], " ")
}

type stringerise func() string

func (s stringerise) String() string { return s() }
