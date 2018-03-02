package ver

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
)

const (
	consulUnitTestUrl = "https://:5c7e1872-77b1-974a-973c-64a8316f3833@consul.konsorten.net/tpo"
)

func TestRetrieveBuildFromConsul(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	rand.Seed(time.Now().UnixNano())

	// prepare dummy data
	vi := MakeVersionInformation()

	vi.SetSemVersion("1.2.3")
	vi.Revision = fmt.Sprintf("abcdef%v", rand.Int31n(8999999)+1000000)

	// get next build id
	err := RetrieveBuildFromConsul(consulUnitTestUrl, "konsorten/build-info/build-info-unit-test", vi)

	if err != nil {
		t.Fatalf("Failed to retrieve build id from consul: %v", err)
	}

	t.Logf("%v\n", vi.VersionString())

	if vi.Build < 1 {
		t.Fatalf("Build ID not updated")
	}

	lastBuildId := vi.Build

	// get same build id
	vi.Build = 0

	err = RetrieveBuildFromConsul(consulUnitTestUrl, "konsorten/build-info/build-info-unit-test", vi)

	if err != nil {
		t.Fatalf("Failed to retrieve build id from consul: %v", err)
	}

	t.Logf("%v\n", vi.VersionString())

	if vi.Build != lastBuildId {
		t.Fatalf("Build ID did change for same revision")
	}
}
