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

	// read host information
	host, err := TryReadFromBuildHost()

	if err != nil {
		t.Fatalf("Failed to read build host info: %v", err)
	}

	// read version information
	found, err := TryReadFromVersionInfoYAML()

	if err != nil {
		t.Fatalf("Failed to read version info: %v", err)
	}

	if found == nil {
		t.Fatalf("%v not found", VersionInfoYamlFilename)
	}

	found.CopyMissingFrom(host)

	if ok, _ := found.IsValid(); !ok {
		t.Fatal("Invalid version information")
	}
}

func TestTryReadVersionInfoYAML_complex(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	// change to test data dir
	currDir, _ := os.Getwd()
	os.Chdir("examples/complex/a/b")
	defer os.Chdir(currDir)

	// read host information
	host, err := TryReadFromBuildHost()

	if err != nil {
		t.Fatalf("Failed to read build host info: %v", err)
	}

	// read version information
	found, err := TryReadFromVersionInfoYAML()

	if err != nil {
		t.Fatalf("Failed to read version info: %v", err)
	}

	if found == nil {
		t.Fatalf("%v not found", VersionInfoYamlFilename)
	}

	found.CopyMissingFrom(host)

	if ok, _ := found.IsValid(); !ok {
		t.Fatal("Invalid version information")
	}
}
