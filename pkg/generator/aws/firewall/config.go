package firewall

import (
	"fmt"
	"strings"
)

type config struct {
	Type      string `config:"type" validate:"required"`
	EventType string `config:"event_type"`
}

func defaultConfig() config {
	return config{
		Type: Name,
	}
}

func (c *config) Validate() error {
	if c.Type != Name {
		return fmt.Errorf("'%s' is not a valid value for 'type' expected '%s'", c.Type, Name)
	}
	if !(c.EventType == "" || c.EventType == EventTypeAlert || c.EventType == EventTypeNetflow) {
		return fmt.Errorf("'%s' is not a valid value for 'event_type' expected '%s'", c.EventType, strings.Join([]string{EventTypeAlert, EventTypeNetflow}, ", "))
	}

	return nil
}
