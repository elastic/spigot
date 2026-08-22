package gotext

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"

	"strconv"

	"github.com/elastic/spigot/pkg/random"
)

type config struct {
	Type   string          `config:"type" validate:"required"`
	Config GeneratorConfig `config:"config" validate:"required"`
}

type GcField struct {
	Name     string   `config:"name"`
	Type     string   `config:"type"`
	Choices  []string `config:"choices"`
	Template *string  `config:"tpl"`
}

type GeneratorConfig struct {
	Name    string    `config:"name" validate:"required"`
	Formats []*string `config:"formats"`
	Fields  []GcField `config:"fields"`
}

func defaultConfig() config {
	return config{
		Type:   Name,
		Config: GeneratorConfig{},
	}
}

func (c *config) Validate() error {
	if c.Type != Name {
		return fmt.Errorf("'%s' is not a valid value for 'type' expected '%s'", c.Type, Name)
	}

	return nil
}

func (f *Field) convert(in bytes.Buffer) any {

	switch f.Type {
	case "Port", "port", "int":
		asString := in.String()
		asInt, err := strconv.Atoi(asString)
		if err != nil {
			log.Fatalf("Could not convert %v to int: %v\n", in, err)
			return nil
		}
		return asInt
	default:
	}
	return in.String()
}

func (f *Field) randomize(object map[string]any) any {
	var buf bytes.Buffer

	// if there is a template, process that
	if f.template != nil {
		err := f.template.Tpl.Execute(&buf, object)
		if err != nil {
			log.Fatal("Failed to execute template", f.template.Format, "with error", err)
			return f.Type
		}

		// need to convert this to the type
		return f.convert(buf)
	}

	// if there are choices, select one at random
	if f.Choices != nil {
		count := len(f.Choices)
		if count > 0 {
			return f.Choices[rand.Intn(count)]
		}
	}

	// if there is a random definition, use that
	switch f.Type {
	case "int":
		return RandomInt(65535)
	case "IPv4", "IP", "ipv4":
		return RandomIPv4()
	case "Port", "port":
		return strconv.Itoa(random.Port())
	case "interface", "Intf", "intf":
		return fmt.Sprintf("%s%02d", f.Name, rand.Intn(16))
	case "Duration", "duration":
		return RandomDuration()
	default:
	}

	// otherwise, return the type as a string
	return f.Type

}
