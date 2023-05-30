package shipper

import (
	"fmt"
	"time"
)

type config struct {
	Type                string        `config:"type" validate:"required"`
	Address             string        `config:"address" validate:"required"`
	InputId             string        `config:"input_id" validate:"required"`
	StreamId            string        `config:"stream_id" validate:"required"`
	DataStreamType      string        `config:"datastream_type" validate:"required"`
	DataStreamDataset   string        `config:"datastream_dataset" validate:"required"`
	DataStreamNamespace string        `config:"datastream_namespace" validate:"required"`
	Timeout             time.Duration `config:"timeout" validate:"required"`
}

func defaultConfig() config {
	return config{
		Type:                Name,
		Address:             "127.0.0.1:5351",
		InputId:             "spigot",
		StreamId:            "spigot",
		DataStreamType:      "logs",
		DataStreamDataset:   "logs",
		DataStreamNamespace: "default",
		Timeout:             15 * time.Second,
	}
}

func (c *config) Validate() error {
	if c.Type != Name {
		return fmt.Errorf("'%s' is not a valid value for 'type' expected '%s'", c.Type, Name)
	}
	return nil
}
