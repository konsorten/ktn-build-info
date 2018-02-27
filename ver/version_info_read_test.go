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
	found, err := TryReadVersionInfoYAML()

	if err != nil {
		t.Fatalf("Failed to read version info: %v", err)
	}

	if found == nil {
		t.Fatalf("%v not found", VersionInfoYamlFilename)
	}

	if !found.IsValid() {
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
	found, err := TryReadVersionInfoYAML()

	if err != nil {
		t.Fatalf("Failed to read version info: %v", err)
	}

	if found == nil {
		t.Fatalf("%v not found", VersionInfoYamlFilename)
	}

	if !found.IsValid() {
		t.Fatal("Invalid version information")
	}
}
