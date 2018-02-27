package ver

import (
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"
)

func TestTryReadPackageJSON(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	// change to test data dir
	currDir, _ := os.Getwd()
	os.Chdir("examples/npm")
	defer os.Chdir(currDir)

	// read version information
	vi := MakeVersionInformation()

	found, err := TryReadPackageJSON(nil, vi)

	if err != nil {
		t.Fatalf("Failed to read version info: %v", err)
	}

	if !found {
		t.Fatalf("%v not found", VersionInfoYamlFilename)
	}

	if !vi.IsValid() {
		t.Fatal("Invalid version information")
	}
}
