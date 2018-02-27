package ver

import (
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"
)

func TestTryReadWorkingDir(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	// change to test data dir
	currDir, _ := os.Getwd()
	os.Chdir("examples/npm")
	defer os.Chdir(currDir)

	// read version information
	found, err := TryReadFromWorkingDirectory()

	if err != nil {
		t.Fatalf("Failed to read version info: %v", err)
	}

	if found == nil {
		t.Fatalf("%v not found", VersionInfoYamlFilename)
	}

	if found.Name != "npm" {
		t.Fatal("Wrong project name")
	}
}
