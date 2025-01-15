package winlog

import (
	"math/rand"
	"strconv"

	"github.com/elastic/spigot/pkg/random"
)

const event4624 = 4624

// randomize4624 generates a random event with
// ID 4624 (An account was successfully logged on).
func randomize4624(g *Generator) Event {
	computerName := RandomComputerName("")

	subjectName := computerName + "$"
	targetName := RandomUser()

	evt := RandomEvent(event4624, g.getTime())
	evt.Provider = Provider{
		Name: "Microsoft-Windows-Security-Auditing",
		GUID: "{54849625-5478-4994-A5BA-3E3B0328C30D}",
	}
	evt.Channel = "Security"
	evt.Computer = computerName
	evt.EventData = EventData{
		Data: []KeyValue{
			{Key: "SubjectUserSid", Value: "S-1-5-18"},
			{Key: "SubjectUserName", Value: subjectName},
			{Key: "SubjectDomainName", Value: "WORKGROUP"},
			{Key: "SubjectLogonId", Value: "0x" + strconv.FormatInt(int64(rand.Intn(65536)), 16)},
			{Key: "TargetUserSid", Value: RandomUserSID(targetName)},
			{Key: "TargetUserName", Value: targetName},
			{Key: "TargetDomainName", Value: computerName},
			{Key: "TargetLogonId", Value: "0x" + strconv.FormatInt(int64(rand.Intn(65536)), 16)},
			{Key: "LogonType", Value: "2"},
			{Key: "LogonProcessName", Value: "User32"},
			{Key: "AuthenticationPackageName", Value: "Negotiate"},
			{Key: "WorkstationName", Value: computerName},
			{Key: "LogonGuid", Value: "{00000000-0000-0000-0000-000000000000}"},
			{Key: "TransmittedServices", Value: "-"},
			{Key: "LmPackageName", Value: "-"},
			{Key: "KeyLength", Value: "0"},
			{Key: "ProcessId", Value: "0x" + strconv.FormatInt(int64(rand.Intn(65536)), 16)},
			{Key: "ProcessName", Value: `C:\\Windows\\System32\\svchost.exe`},
			{Key: "IpAddress", Value: random.IPv4().String()},
			{Key: "IpPort", Value: strconv.Itoa(random.Port())},
			{Key: "ImpersonationLevel", Value: "%%1833"},
			{Key: "RestrictedAdminMode", Value: "-"},
			{Key: "TargetOutboundUserName", Value: "-"},
			{Key: "TargetOutboundDomainName", Value: "-"},
			{Key: "VirtualAccount", Value: "%%1843"},
			{Key: "TargetLinkedLogonId", Value: "0x0"},
			{Key: "ElevatedToken", Value: "%%1842"},
		},
	}

	return evt
}
