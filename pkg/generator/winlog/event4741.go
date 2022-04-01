package winlog

import (
	"fmt"
	"math/rand"
	"strconv"
)

const event4741 = 4741

// randomize4741 generates a random event with
// ID 4741 (A computer account was created).
func randomize4741(g *Generator) Event {
	now := g.getTime()

	domain := RandomDomain()
	hostname := RandomComputerName("")
	computerName := hostname + "." + domain

	targetName := hostname + "$"
	subjectName := RandomUser()

	evt := RandomEvent(event4741, g.getTime())
	evt.Provider = Provider{
		Name: "Microsoft-Windows-Security-Auditing",
		GUID: "{54849625-5478-4994-A5BA-3E3B0328C30D}",
	}
	evt.Channel = "Security"
	evt.Computer = computerName
	evt.EventData = EventData{
		Data: []KeyValue{
			{Key: "TargetUserSid", Value: RandomUserSID(targetName)},
			{Key: "TargetUserName", Value: targetName},
			{Key: "TargetDomainName", Value: domain},
			{Key: "SubjectUserSid", Value: RandomUserSID(subjectName)},
			{Key: "SubjectUserName", Value: subjectName},
			{Key: "SubjectDomainName", Value: domain},
			{Key: "SubjectLogonId", Value: "0x" + strconv.FormatInt(int64(rand.Intn(65536)), 16)},
			{Key: "PrivilegeList", Value: "-"},
			{Key: "SamAccountName", Value: hostname + "$"},
			{Key: "DisplayName", Value: "-"},
			{Key: "UserPrincipalName", Value: "-"},
			{Key: "HomeDirectory", Value: "-"},
			{Key: "HomePath", Value: "-"},
			{Key: "ScriptPath", Value: "-"},
			{Key: "ProfilePath", Value: "-"},
			{Key: "UserWorkstations", Value: "-"},
			{Key: "PasswordLastSet", Value: now.Format("2/1/2006 03:04:05 PM")},
			{Key: "AccountExpires", Value: "%%1794"},
			{Key: "PrimaryGroupId", Value: strconv.Itoa(rand.Intn(10000))},
			{Key: "AllowedToDelegateTo", Value: "-"},
			{Key: "OldUacValue", Value: "0x0"},
			{Key: "NewUacValue", Value: "0x80"},
			{Key: "UserAccountControl", Value: "%%2087"},
			{Key: "UserParameters", Value: "-"},
			{Key: "SidHistory", Value: "-"},
			{Key: "LogonHours", Value: "%%1793"},
			{Key: "DnsHostName", Value: computerName},
			{Key: "ServicePrincipalNames", Value: fmt.Sprintf(
				"HOST/%s RestrictedKrbHost/%s HOST/%s RestrictedKrbHost/%s",
				computerName,
				computerName,
				hostname,
				hostname,
			)},
		},
	}

	return evt
}
