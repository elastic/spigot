package cef

import (
	"fmt"
	"net/http"
	"time"
)

type config struct {
	Type        string   `config:"type" validate:"required"`
	CEFVersions []int    `config:"cef"`
	Vendors     []string `config:"vendors" validate:"required"`
	Products    []string `config:"products" validate:"required"`
	Versions    []string `config:"versions" validate:"required"`
	Classes     []string `config:"classes" validate:"required"`
	Names       []string `config:"names" validate:"required"`
	Severities  []int    `config:"severities"`

	Users      []string `config:"users"`
	Privs      []string `config:"privs"`
	Methods    []string `config:"methods"`
	Interfaces []string `config:"interfaces"`
	TimeZones  []string `config:"timezones"`
	Actions    []string `config:"actions"`
	Messages   []string `config:"messages"`
	Words      []string `config:"words"`
	Text       string   `config:"text"`

	Must    []string `config:"must_include"`
	Exclude []string `config:"must_exclude"`
	Max     int      `config:"max_extensions"`

	Now      func() time.Time
	ZeroUUID bool
}

func defaultConfig() config {
	return config{
		Type: Name,
		Now:  time.Now,
	}
}

func (c *config) Validate() error {
	if c.Type != Name {
		return fmt.Errorf("'%s' is not a valid value for 'type' expected '%s'", c.Type, Name)
	}
	if len(c.CEFVersions) == 0 {
		c.CEFVersions = []int{0}
	}
	if len(c.Severities) == 0 {
		c.Severities = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	}
	if len(c.Words) == 0 {
		c.Words = words
	}
	if len(c.Text) == 0 {
		c.Text = loremIpsum
	}
	if len(c.Users) == 0 {
		c.Users = users
	}
	if len(c.Privs) == 0 {
		c.Privs = privs
	}
	if len(c.Methods) == 0 {
		c.Methods = methods
	}
	if len(c.Interfaces) == 0 {
		c.Interfaces = interfaces
	}
	if len(c.Actions) == 0 {
		c.Actions = actions
	}
	if len(c.Messages) == 0 {
		c.Messages = messages
	}
	return nil
}

const loremIpsum = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

var (
	interfaces = []string{
		"eth0",
		"eth1",
	}
	users = []string{
		"alice",
		"bob",
		"eve",
		"mallory",
	}
	privs = []string{
		"Administrator",
		"User",
		"Guest",
	}
	timeZones = []string{
		"Europe/London",
		"Europe/Paris",
		"America/New_York",
	}
	methods = []string{
		http.MethodConnect,
		http.MethodDelete,
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
	}
	messages = []string{
		"Signature violation rule ID 807: web-cgi /wwwboard/passwd.txt access",
		"Disallow Illegal URL.",
		"Transformed (xout) potential credit card numbers seen in server response",
		"Maximum number of potential credit card numbers seen",
		"Field consistency check failed for field passwd",
	}
	actions = []string{
		"Accept", "Bypass", "Drop",
	}
	words = []string{
		"aware",
		"accessible",
		"aggressive",
		"alike",
		"average",
		"bathe",
		"behave",
		"bite-sized",
		"boy",
		"bulb",
		"buzz",
		"caption",
		"cars",
		"channel",
		"choke",
		"class",
		"dazzling",
		"disagreeable",
		"dramatic",
		"ducks",
		"dynamic",
		"elite",
		"enchanting",
		"encouraging",
		"end",
		"excite",
		"excuse",
		"exotic",
		"eyes",
		"fearless",
		"fix",
		"flat",
		"four",
		"futuristic",
		"gifted",
		"great",
		"green",
		"hanging",
		"happy",
		"hard",
		"horse",
		"house",
		"huge",
		"hysterical",
		"identify",
		"ink",
		"kill",
		"linen",
		"living",
		"long-term",
		"lunchroom",
		"man",
		"march",
		"melodic",
		"monkey",
		"muddled",
		"murky",
		"neat",
		"nutritious",
		"obnoxious",
		"obsequious",
		"oil",
		"old",
		"parched",
		"payment",
		"pedal",
		"pine",
		"pointless",
		"poor",
		"preach",
		"previous",
		"psychedelic",
		"radiate",
		"ragged",
		"rely",
		"ring",
		"romantic",
		"seashore",
		"seat",
		"share",
		"skate",
		"slip",
		"soak",
		"solid",
		"spoon",
		"spray",
		"squeamish",
		"stir",
		"stranger",
		"tacit",
		"test",
		"thinkable",
		"thoughtless",
		"tiny",
		"tough",
		"towering",
		"unadvised",
		"vagabond",
		"women",
	}
)
