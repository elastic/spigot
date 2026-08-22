//go:build !windows

package syslog

import "fmt"

type config struct {
	Type     string `config:"type" validate:"required"`
	Facility string `config:"facility"`
	Severity string `config:"severity"`
	Tag      string `config:"tag"`
	Network  string `config:"network" validate:"required"`
	Host     string `config:"host" validate:"required"`
	Port     string `config:"port" validate:"required"`
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
	return nil
}
