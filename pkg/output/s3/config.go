package s3

import "fmt"

type config struct {
	Type      string `config:"type" validate:"required"`
	Bucket    string `config:"bucket" validate:"required"`
	Region    string `config:"region" validate:"required"`
	Delimiter string `config:"delimiter"`
	Prefix    string `config:"prefix" validate:"required"`
}

func defaultConfig() config {
	return config{
		Type:      Name,
		Delimiter: "\n",
	}
}

func (c *config) Validate() error {
	if c.Type != Name {
		return fmt.Errorf("'%s' is not a valid value for 'type' expected '%s'", c.Type, Name)
	}
	return nil
}
