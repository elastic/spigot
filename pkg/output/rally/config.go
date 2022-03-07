package rally

import "fmt"

type config struct {
	Type      string `config:"type" validate:"required"`
	Directory string `config:"directory"`
	Pattern   string `config:"pattern"`
	Filename  string `config:"filename"`
}

func defaultConfig() config {
	return config{
		Type: Name,
	}
}

func (c *config) Validate() error {
	if c.Type != Name {
		return fmt.Errorf("%s is not a valid type for %s", c.Type, Name)
	}
	if c.Filename != "" && (c.Directory != "" || c.Pattern != "") {
		return fmt.Errorf("if filename is set, directory and pattern must not be")
	}
	if (c.Directory != "" && c.Pattern == "") || (c.Directory == "" && c.Pattern != "") {
		return fmt.Errorf("directory and pattern must both be set")
	}
	if c.Filename == "" && c.Directory == "" && c.Pattern == "" {
		return fmt.Errorf("you must specify filename or directory and pattern")
	}
	return nil
}
