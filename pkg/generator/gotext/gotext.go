// Package gotext implements the generator for generic logs.
//
// Configuration file supports including timestamps in log messages
//
//	generator:
//	  type: gotext
//	  include_timestamp: true
package gotext

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/elastic/go-ucfg"
	"github.com/elastic/spigot/pkg/generator"
	"github.com/elastic/spigot/pkg/random"
)

// Name is the name of the generator in the configuration file and registry
const Name = "gotext"

func init() {
	generator.Register(Name, New)
}

var (
	FunctionMap = template.FuncMap{
		"ToLower":    strings.ToLower,
		"ToUpper":    strings.ToUpper,
		"RandomIPv4": RandomIPv4,
		"RandomPort": RandomPort,
		"RandomInt":  RandomInt,
		"Percent":    Percent,
		"PlusInt":    PlusInt,
		"TimesInt":   TimesInt,
	}
)

type TimestampFormatter struct {
	cur time.Time
}

func (tf TimestampFormatter) seconds(precision int) string {
	var returnVal string

	secs := tf.cur.Unix()
	returnVal = fmt.Sprintf("%d", secs)

	if precision == 0 {
		return returnVal
	} else if precision > 9 {
		precision = 9
	}

	nanos := tf.cur.UnixNano() % 1_000_000_000
	partialsString := fmt.Sprintf("%09d", nanos)

	returnVal = fmt.Sprintf("%d.%s", secs, partialsString[0:precision])

	return returnVal
}

func (tf TimestampFormatter) Seconds(format any, _ string) string {

	precision := ToIntWithDefault(format, 0)
	return tf.seconds(precision)
}

func (tf TimestampFormatter) Seconds3(format, _ string) string {
	return tf.seconds(3)
}
func (tf TimestampFormatter) Seconds6(format, _ string) string {
	return tf.seconds(6)
}
func (tf TimestampFormatter) Seconds9(format, _ string) string {
	return tf.seconds(9)
}
func (tf TimestampFormatter) Format(format, _ string) string {
	return tf.cur.Format(format)
}

func ToIntWithDefault(input any, defaultValue int) int {

	if v, ok := input.(int); ok {
		return v
	}
	switch v := input.(type) {
	case string:
		result, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("Could not convert %v (%T) to int: %v\n", input, input, err)
			return 1
		}
		return result
	default:
	}
	return defaultValue

}
func ToInt(input any) int {
	return ToIntWithDefault(input, 1)
}

func PlusInt(a, b any) string {

	return fmt.Sprintf("%v", ToInt(a)+ToInt(b))
}
func TimesInt(a, b any) string {

	return fmt.Sprintf("%v", ToInt(a)*ToInt(b))
}

func Percent(numerator, denominator any) string {
	fnum := float64(ToInt(numerator))
	dnum := float64(ToInt(denominator))

	result := fnum / dnum * 100.0
	return fmt.Sprintf("%8.6f", result)
}

func RandomDuration() string {
	// return the string interpretation of that value
	return fmt.Sprintf("%01d:%02d:%02d", rand.Intn(4), rand.Intn(60), rand.Intn(60))
}

func RandomInt(maximum int) string {
	randval := 0
	if maximum > 0 {
		// get a random value
		randval = rand.Intn(maximum)
	} else if maximum < 0 {
		// get a random value
		randval = -1 * rand.Intn(-1*maximum)
	}

	// return the string interpretation of that value
	return strconv.Itoa(randval)
}

func RandomIPv4() string {
	return random.IPv4().String()
}

func RandomPort() string {
	return strconv.Itoa(random.Port())
}

type Template struct {
	Format string
	Tpl    *template.Template
}

type Field struct {
	Name     string   `config:"name"`
	Type     string   `config:"type"`
	Choices  []string `config:"choices"`
	template *Template
}

type GoText struct {
	Name          string
	Fields        []Field
	templates     []Template
	timeFormatter TimestampFormatter
}

func (g *GoText) Next() ([]byte, error) {
	var buf bytes.Buffer

	object := make(map[string]any)

	object["Timestamp"] = g.timeFormatter.Format("2006-01-02 15:04:05", "")
	timeFormatter := TimestampFormatter{cur: g.timeFormatter.cur}
	g.timeFormatter.cur = g.timeFormatter.cur.Add(time.Millisecond * 10)

	object["TimestampFormatter"] = &timeFormatter

	// loop over each field
	for _, f := range g.Fields {
		object[f.Name] = f.randomize(object)
	}

	// are there formats?
	if len(g.templates) < 1 {
		fmt.Printf("i am %v; %v\n", g.Name, g.templates)
		return nil, fmt.Errorf("this has no templates to process; bailing")
	}
	index := rand.Intn(len(g.templates))

	// attempt to generate each one
	err := g.templates[index].Tpl.Execute(&buf, object)
	if err != nil {
		log.Fatal("Failed to execute template", g.templates[index].Format, "with error", err)
		return nil, err
	}

	return buf.Bytes(), err
}

// New is Factory for the gotext generator
func New(cfg *ucfg.Config, numRecords int64) (generator.Generator, error) {
	c := defaultConfig()
	if err := cfg.Unpack(&c); err != nil {
		return nil, err
	}

	gotextConfig := c.Config

	startTime := time.Now().Add(-time.Duration(numRecords) * 10 * time.Millisecond)

	timeFormatter := TimestampFormatter{cur: startTime}

	// check variables
	// return
	g := &GoText{
		Name:          gotextConfig.Name,
		Fields:        nil,
		templates:     nil,
		timeFormatter: timeFormatter,
	}

	for i, v := range gotextConfig.Fields {
		f := Field{
			Name:    v.Name,
			Type:    v.Type,
			Choices: v.Choices,
		}

		// if there is a Template field
		if v.Template != nil {
			t, err := template.New(strconv.Itoa(i)).Funcs(FunctionMap).Parse(*v.Template)
			if err != nil {
				return nil, err
			}

			f.template = &Template{
				Format: *v.Template,
				Tpl:    t,
			}
		}
		g.Fields = append(g.Fields, f)
	}

	for i, v := range gotextConfig.Formats {
		t, err := template.New(strconv.Itoa(i)).Funcs(FunctionMap).Parse(*v)
		if err != nil {
			return nil, err
		}
		g.templates = append(g.templates, Template{
			Format: *v,
			Tpl:    t,
		})
	}

	return g, nil
}
