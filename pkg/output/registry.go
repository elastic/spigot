package output

import (
	"fmt"

	"github.com/elastic/go-ucfg"
)

type Factory = func(*ucfg.Config) (Output, error)

var registry = make(map[string]Factory)

func Register(name string, factory Factory) error {
	if _, exists := registry[name]; exists {
		return fmt.Errorf("Error registering input '%s': already registered", name)
	}
	registry[name] = factory
	return nil
}

func GetFactory(name string) (Factory, error) {
	factory, exists := registry[name]
	if !exists {
		return nil, fmt.Errorf("Input %s not registered", name)
	}
	return factory, nil
}
