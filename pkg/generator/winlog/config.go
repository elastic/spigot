package winlog

import (
	"fmt"
)

type config struct {
	Type       string `config:"type" validate:"required"`
	EventID    int    `config:"event_id"`
	AsTemplate bool   `config:"as_template"`
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
	if c.EventID != 0 {
		if _, ok := eventRandomizers[c.EventID]; !ok {
			return fmt.Errorf("'%d' is not a valid value for 'event_id' expected one of %v", c.EventID, eventIDs)
		}
	}

	return nil
}
