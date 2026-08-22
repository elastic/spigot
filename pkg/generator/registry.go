package generator

import (
	"fmt"

	"github.com/elastic/go-ucfg"
)

// Factory is the function signature of each generators New function.
// Given a config it returns a generator or an error.
type Factory = func(*ucfg.Config, int64) (Generator, error)

var registry = make(map[string]Factory)

// Register associates a generator name with the generator factory.
func Register(name string, factory Factory) error {
	if _, exists := registry[name]; exists {
		return fmt.Errorf("Error registering input '%s': already registered", name)
	}
	registry[name] = factory
	return nil
}

// GetFactory retrieves a factory for a given name, or returns an
// error if there isn't a factory associated with that name.
func GetFactory(name string) (Factory, error) {
	factory, exists := registry[name]
	if !exists {
		return nil, fmt.Errorf("Input %s not registered", name)
	}
	return factory, nil
}
