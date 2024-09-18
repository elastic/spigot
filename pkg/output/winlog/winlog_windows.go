// Package winlog implements the output of logs windows event log.
//
// Configuration file supports either writing writing passed in messages
// as events with IDs 1-1000 for simple purposes, or to pass the values
// as event templates. Take into account that `event_create_msg_file` should
// exist and contain the appropiate messages to render the required events.
// By default it will use `%SystemRoot%\System32\EventCreate.exe` which allow
// the rendering of events IDs in the 1-1000 range without further configuration.
// The channel will be cleared and removed after running, if `persist_events` is set
// to `true` the events will remain after the execution.
// If you use it with `spigot`'s `winlog` generator set `event_create_msg_file` to
// `winlog-generator` so it uses the right file.

// To manually remove the event log you can run `Remove-EventLog -LogName "WinlogbeatTest"`
// in PowerShell.
//
//	output:
//	  type: winlog
//	  event_create_msg_file: "%SystemRoot%\\system32\\adtschema.dll"
//	  templated: true
//	  persist_events: true
package winlog

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os/exec"
	"strconv"
	"time"

	"github.com/andrewkroh/sys/windows/svc/eventlog"
	"github.com/elastic/beats/v7/winlogbeat/sys/wineventlog"
	"github.com/elastic/go-ucfg"
	"github.com/leehinman/spigot/pkg/output"
	"go.uber.org/multierr"
)

// OutputName is the name of the output in the configuration file and registry
const Name = "winlog"

type Output struct {
	config
	log *eventlog.Log
}

func init() {
	output.Register(Name, New)
}

// New is the Factory for creating a new winlog output.
func New(cfg *ucfg.Config) (output.Output, error) {
	var err error

	c := defaultConfig()
	if err = cfg.Unpack(&c); err != nil {
		return nil, err
	}
	if c.EventCreateMsgFile == "winlog-generator" {
		c.EventCreateMsgFile = TemplateMessageFile
	}

	log, err := createLog(c.Provider, c.Source, c.EventCreateMsgFile)
	if err != nil {
		return nil, err
	}

	if err := setLogSize(c.Provider, c.WinlogSizeInBytes); err != nil {
		return nil, err
	}

	return &Output{
		config: c,
		log:    log,
	}, nil
}

func (o *Output) Write(b []byte) (n int, err error) {
	if o.log == nil {
		return 0, errors.New("the output is closed and unusable")
	}
	var etype uint16
	var eid uint32
	var messages []string
	switch o.Templated {
	case true:
		tpl := &EventTemplate{}
		if err := json.Unmarshal(b, tpl); err != nil {
			return 0, err
		}
		etype = tpl.EventType
		eid = tpl.EventID
		messages = tpl.Messages
	case false:
		etype = eventlog.Info
		eid = uint32(rand.Int63() % 1000)
		messages = []string{string(b)}
	}
	return len(b), safeWriteEvent(o.log, etype, eid, messages)
}

// Close the winlog. Writes after this will fail.
func (o *Output) Close() error {
	if o.PersistsEvents {
		return o.log.Close()
	}

	err := multierr.Combine(
		o.log.Close(),
		wineventlog.EvtClearLog(wineventlog.NilHandle, o.Provider, ""),
		eventlog.RemoveSource(o.Provider, o.Source),
		eventlog.RemoveProvider(o.Provider),
	)
	o.log = nil
	return err
}

func (o *Output) NewInterval() error {
	return nil
}

// createLog creates a new event log and returns a handle for writing events
// to the log.
func createLog(name, source, messageFile string) (*eventlog.Log, error) {
	existed, err := eventlog.Install(name, source, messageFile, true, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		return nil, err
	}

	if existed {
		wineventlog.EvtClearLog(wineventlog.NilHandle, name, "")
	}

	log, err := eventlog.Open(source)
	if err != nil {
		return nil, multierr.Combine(
			err,
			eventlog.RemoveSource(name, source),
			eventlog.RemoveProvider(name),
		)
	}

	return log, nil
}

func safeWriteEvent(log *eventlog.Log, etype uint16, eid uint32, msgs []string) error {
	deadline := time.Now().Add(time.Second * 10)
	for {
		err := log.Report(etype, eid, msgs)
		if err == nil {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("failed to write event to event log: %w", err)
		}
	}
}

// setLogSize set the maximum number of bytes that an event log can hold.
func setLogSize(provider string, sizeBytes int) error {
	output, err := exec.Command("wevtutil.exe", "sl", "/ms:"+strconv.Itoa(sizeBytes), provider).CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set log size: %w, %s", err, string(output))
	}
	return nil
}
