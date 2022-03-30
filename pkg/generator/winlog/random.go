package winlog

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var (
	mapMu sync.Mutex

	serviceSIDMap = map[string]string{}
	userSIDMap    = map[string]string{}
)

// RandomUser generates a random user name.
func RandomUser() string {
	return "user" + strconv.Itoa(rand.Intn(100))
}

// RandomComputerName generates a random computer name. If domain is provided,
// it will be app.
func RandomComputerName(domain string) string {
	name := "COMPUTER-" + strconv.Itoa(rand.Intn(1000))
	if domain != "" {
		name += "." + domain
	}

	return name
}

// RandomDomain generates a random domain.
func RandomDomain() string {
	return "DOMAIN-" + strconv.Itoa(rand.Intn(10))
}

// RandomSID generates a random SID.
func RandomSID() string {
	return fmt.Sprintf(
		"S-1-5-21-%d-%d-%d-%d",
		rand.Intn(1<<32),
		rand.Intn(1<<32),
		rand.Intn(1<<32),
		rand.Intn(1<<16),
	)
}

// RandomServiceSID generates a random SID for a service with name. If a SID
// has already been generated for this name, it will be returned.
func RandomServiceSID(name string) string {
	mapMu.Lock()
	defer mapMu.Unlock()

	if sid, ok := serviceSIDMap[name]; ok {
		return sid
	}

	sid := RandomSID()
	serviceSIDMap[name] = sid

	return sid
}

// RandomUserSID generates a random SID for a user with name. If a SID
// has already been generated for this name, it will be returned.
func RandomUserSID(name string) string {
	mapMu.Lock()
	defer mapMu.Unlock()

	if sid, ok := userSIDMap[name]; ok {
		return sid
	}

	sid := RandomSID()
	userSIDMap[name] = sid

	return sid
}

func RandomEvent(eventID uint32, now time.Time) Event {
	return Event{
		EventID: EventID{
			ID: eventID,
		},
		Task:     uint16(rand.Intn(65536)),
		Keywords: 0x8020000000000000,
		TimeCreated: TimeCreated{
			SystemTime: now,
		},
		RecordID:    rand.Uint64(),
		Correlation: Correlation{},
		Execution: Execution{
			ProcessID: uint32(rand.Intn(65536)),
			ThreadID:  uint32(rand.Intn(65536)),
		},
	}
}
