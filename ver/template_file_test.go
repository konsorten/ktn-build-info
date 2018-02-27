package ver

import (
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"
)

func TestRenderTemplateFile(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	// change to test data dir
	currDir, _ := os.Getwd()
	os.Chdir("examples/template")
	defer os.Chdir(currDir)

	// read host information
	host, err := TryReadFromBuildHost()

	if err != nil {
		t.Fatalf("Failed to read build host info: %v", err)
	}

	// read version information
	found, err := TryReadVersionInfoYAML()

	if err != nil {
		t.Fatalf("Failed to read version info: %v", err)
	}

	if found == nil {
		t.Fatalf("%v not found", VersionInfoYamlFilename)
	}

	found.CopyMissingFrom(host)
	found.Build = 4
	found.Revision = "abcdef0"

	if !found.IsValid() {
		t.Fatal("Invalid version information")
	}

	// render template
	err = found.WriteTemplateFile("test.json.template")

	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}
}
