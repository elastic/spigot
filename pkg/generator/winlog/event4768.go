package winlog

import (
	"strconv"

	"github.com/elastic/spigot/pkg/random"
)

const event4768 = 4768

// randomize4768 generates a random event with
// ID 4768 (A Kerberos authentication ticket (TGT) was requested).
func randomize4768(g *Generator) Event {
	domain := RandomDomain()
	computerName := RandomComputerName(domain)

	target := RandomUser()

	evt := RandomEvent(event4768, g.getTime())
	evt.Provider = Provider{
		Name: "Microsoft-Windows-Security-Auditing",
		GUID: "{54849625-5478-4994-A5BA-3E3B0328C30D}",
	}
	evt.Channel = "Security"
	evt.Computer = computerName
	evt.EventData = EventData{
		Data: []KeyValue{
			{Key: "TargetUserName", Value: target},
			{Key: "TargetDomainName", Value: domain},
			{Key: "TargetSid", Value: RandomUserSID(target)},
			{Key: "ServiceName", Value: "krbtgt"},
			{Key: "TargetSid", Value: RandomServiceSID("krbtgt")},
			{Key: "TicketOptions", Value: "0x40810010"},
			{Key: "TicketEncryptionType", Value: "0x12"},
			{Key: "PreAuthType", Value: "15"},
			{Key: "IpAddress", Value: random.IPv4().String()},
			{Key: "IpPort", Value: strconv.Itoa(random.Port())},
			{Key: "CertIssuerName", Value: domain + "-CA-1"},
			{Key: "CertSerialNumber", Value: "1D0000000D292FBE3C6CDDAFA200020000000D"},
			{Key: "CertThumbprint", Value: "564DFAEE99C71D62ABC553E695BD8DBC46669413"},
		},
	}

	return evt
}
