package winlog

const TemplateMessageFile = "%SystemRoot%\\system32\\adtschema.dll"

type EventTemplate struct {
	EventType uint16   `json:"event_type"`
	EventID   uint32   `json:"event_id"`
	Messages  []string `json:"messages"`
}
