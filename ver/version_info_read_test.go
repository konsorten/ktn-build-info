package ver

import (
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"
)

func TestTryReadVersionInfoYAML_simple(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	// change to test data dir
	currDir, _ := os.Getwd()
	os.Chdir("examples/simple")
	defer os.Chdir(currDir)

	// read version information
	vi := MakeVersionInformation()

	found, err := TryReadVersionInfoYAML(nil, vi)

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

func TestTryReadVersionInfoYAML_complex(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	// change to test data dir
	currDir, _ := os.Getwd()
	os.Chdir("examples/complex/a/b")
	defer os.Chdir(currDir)

	// read version information
	vi := MakeVersionInformation()

	found, err := TryReadVersionInfoYAML(nil, vi)

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
