// Package winlog generates Windows Event Log XML records.
//
// Configuration:
//
//	event_id: (number, optional) If provided, generated events using this ID.
//	          Must be one of the registered event IDs. See 'eventRandomizers'
//	          for the list of valid IDs. If not provided, the generator will
//	          randomly select from the available list for each record.
//
//	- generator:
//	    type: winlog
//	    event_id: 4768
package winlog

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/elastic/go-ucfg"
	"github.com/elastic/spigot/pkg/generator"
	"github.com/elastic/spigot/pkg/output/winlog"
)

const Name = "winlog"

type randomizerFunc func(g *Generator) Event

var (
	eventRandomizers = map[int]randomizerFunc{
		event4624: randomize4624,
		event4634: randomize4634,
		event4723: randomize4723,
		event4741: randomize4741,
		event4743: randomize4743,
		event4768: randomize4768,
	}
	eventIDs []int // Populated at runtime based on 'eventRandomizers' keys.
)

// Event holds the random fields for a Windows Event Log record.
type Event struct {
	XMLName xml.Name `xml:"http://schemas.microsoft.com/win/2004/08/events/event Event"`

	Provider    Provider    `xml:"System>Provider"`
	EventID     EventID     `xml:"System>EventID"`
	Version     uint8       `xml:"System>Version"`
	Level       uint8       `xml:"System>Level"`
	Task        uint16      `xml:"System>Task"`
	Opcode      uint8       `xml:"System>Opcode"`
	Keywords    HexUint64   `xml:"System>Keywords"`
	TimeCreated TimeCreated `xml:"System>TimeCreated"`
	RecordID    uint64      `xml:"System>EventRecordID"`
	Correlation Correlation `xml:"System>Correlation"`
	Execution   Execution   `xml:"System>Execution"`
	Channel     string      `xml:"System>Channel"`
	Computer    string      `xml:"System>Computer"`
	Security    Security    `xml:"System>Security"`

	EventData EventData `xml:"EventData"`
}

func (e *Event) AsTemplate() winlog.EventTemplate {
	messages := make([]string, len(e.EventData.Data))
	for i, data := range e.EventData.Data {
		messages[i] = data.Value
	}
	return winlog.EventTemplate{
		EventType: uint16(e.Level),
		EventID:   e.EventID.ID,
		Messages:  messages,
	}
}

// Provider identifies the provider that logged the event.
type Provider struct {
	Name            string `xml:"Name,attr,omitempty"`
	GUID            string `xml:"GUID,attr,omitempty"`
	EventSourceName string `xml:"EventSourceName,attr,omitempty"`
}

// EventID is the identifier that the provider uses to identify a
// specific event type.
type EventID struct {
	Qualifiers uint16 `xml:"Qualifiers,attr,omitempty"`
	ID         uint32 `xml:",chardata"`
}

// HexUint64 is a uint64. When marshaled, it will be in hexadecimal format.
type HexUint64 uint64

func (v HexUint64) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	return enc.EncodeElement(fmt.Sprintf("%#x", v), start)
}

// TimeCreated contains the system time of when the event was logged.
type TimeCreated struct {
	SystemTime time.Time `xml:"SystemTime,attr"`
}

// Correlation contains activity identifiers that consumers can use to group
// related events together.
type Correlation struct {
	ActivityID        string `xml:"ActivityID,attr,omitempty"`
	RelatedActivityID string `xml:"RelatedActivityID,attr,omitempty"`
}

// Execution contains information about the process and thread that logged the event.
type Execution struct {
	ProcessID uint32 `xml:"ProcessID,attr"`
	ThreadID  uint32 `xml:"ThreadID,attr"`

	// Only available for events logged to an event tracing log file (.etl file).
	ProcessorID   uint32 `xml:"ProcessorID,attr,omitempty"`
	SessionID     uint32 `xml:"SessionID,attr,omitempty"`
	KernelTime    uint32 `xml:"KernelTime,attr,omitempty"`
	UserTime      uint32 `xml:"UserTime,attr,omitempty"`
	ProcessorTime uint32 `xml:"ProcessorTime,attr,omitempty"`
}

// EventData contains the event data.
type EventData struct {
	Data []KeyValue `xml:",any"`
}

// KeyValue is a key value pair of strings.
type KeyValue struct {
	Key   string `xml:"Name,attr"`
	Value string `xml:",chardata"`
}

// Security represents the Windows Security Identifier for an account.
type Security struct {
	UserID string `xml:"UserID,attr,omitempty"`
}

// Generator provides a Windows Event XML record generator.
type Generator struct {
	Event Event

	eventID    *int
	staticTime *time.Time
	render     func(Event) ([]byte, error)
}

// Next produces the next Windows Event XML record.
func (g *Generator) Next() ([]byte, error) {
	var eventID int

	if g.eventID != nil {
		eventID = *g.eventID
	} else {
		eventID = eventIDs[rand.Intn(len(eventIDs))]
	}
	fn, ok := eventRandomizers[eventID]
	if !ok {
		return nil, fmt.Errorf("event ID %d is not registered with this generator", eventID)
	}

	g.Event = fn(g)

	return g.render(g.Event)
}

func (g *Generator) getTime() time.Time {
	if g.staticTime != nil {
		return *g.staticTime
	}

	return time.Now()
}

// New is the factory for Windows Event XML objects.
func New(cfg *ucfg.Config) (generator.Generator, error) {
	c := defaultConfig()
	if err := cfg.Unpack(&c); err != nil {
		return nil, err
	}

	g := Generator{}
	if c.EventID > 0 {
		g.eventID = &c.EventID
	}

	if c.AsTemplate {
		g.render = func(e Event) ([]byte, error) {
			return json.Marshal(e.AsTemplate())
		}
	} else {
		g.render = func(e Event) ([]byte, error) {
			return xml.MarshalIndent(&g.Event, "", "  ")
		}
	}

	return &g, nil
}

func init() {
	for k := range eventRandomizers {
		eventIDs = append(eventIDs, k)
	}
	sort.Ints(eventIDs)

	_ = generator.Register(Name, New)
}
