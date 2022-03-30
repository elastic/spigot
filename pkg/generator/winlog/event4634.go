package winlog

import (
	"math/rand"
	"strconv"
)

const event4634 = 4634

// randomize4634 generates a random event with
// ID 4634 (An account was logged off).
func randomize4634(g *Generator) Event {
	domain := RandomDomain()
	computerName := RandomComputerName(domain)

	target := RandomUser()

	evt := RandomEvent(event4634, g.getTime())
	evt.Provider = Provider{
		Name: "Microsoft-Windows-Security-Auditing",
		GUID: "{54849625-5478-4994-A5BA-3E3B0328C30D}",
	}
	evt.Channel = "Security"
	evt.Computer = computerName
	evt.EventData = EventData{
		Data: []KeyValue{
			{Key: "TargetUserSid", Value: RandomUserSID(target)},
			{Key: "TargetUserName", Value: target},
			{Key: "TargetDomainName", Value: domain},
			{Key: "TargetLogonId", Value: "0x" + strconv.FormatInt(int64(rand.Intn(65536)), 16)},
			{Key: "LogonType", Value: "2"},
		},
	}

	return evt
}
