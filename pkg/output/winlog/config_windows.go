package winlog

import "fmt"

const gigabyte = 1 << 30

type config struct {
	Type               string `config:"type" validate:"required"`
	Provider           string `config:"provider" validate:"required"`
	Source             string `config:"source" validate:"required"`
	EventCreateMsgFile string `config:"event_create_msg_file" validate:"required"`
	WinlogSizeInBytes  int    `config:"winlog_size_in_bytes" validate:"required"`
	Templated          bool   `config:"templated"`
	PersistsEvents     bool   `config:"persist_events"`
}

func defaultConfig() config {
	return config{
		Type:               Name,
		Provider:           "WinlogbeatTest",
		Source:             "Benchmark",
		EventCreateMsgFile: "%SystemRoot%\\System32\\EventCreate.exe",
		WinlogSizeInBytes:  gigabyte,
	}
}

func (c *config) Validate() error {
	if c.Type != Name {
		return fmt.Errorf("%s is not a valid type for %s", c.Type, Name)
	}
	if c.WinlogSizeInBytes < 1 {
		return fmt.Errorf("winlog_size_in_bytes must be positive")
	}
	return nil
}
