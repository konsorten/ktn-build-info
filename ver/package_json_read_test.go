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

	// read host information
	host, err := TryReadFromBuildHost()

	if err != nil {
		t.Fatalf("Failed to read build host info: %v", err)
	}

	// read version information
	found, err := TryReadPackageJSON()

	if err != nil {
		t.Fatalf("Failed to read version info: %v", err)
	}

	if found == nil {
		t.Fatalf("%v not found", PackageJsonFilename)
	}

	found.CopyMissingFrom(host)

	if !found.IsValid() {
		t.Fatal("Invalid version information")
	}
}

func TestTryReadPackageJSON_oldAuthor(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	// change to test data dir
	currDir, _ := os.Getwd()
	os.Chdir("examples/npm-old-author")
	defer os.Chdir(currDir)

	// read defaults
	defaults, err := TryReadFromKonsortenDefaults()

	if err != nil {
		t.Fatalf("Failed to read defaults: %v", err)
	}

	// read host information
	host, err := TryReadFromBuildHost()

	if err != nil {
		t.Fatalf("Failed to read build host info: %v", err)
	}

	// read version information
	found, err := TryReadPackageJSON()

	if err != nil {
		t.Fatalf("Failed to read version info: %v", err)
	}

	if found == nil {
		t.Fatalf("%v not found", PackageJsonFilename)
	}

	found.CopyMissingFrom(host)
	found.CopyMissingFrom(defaults)

	if !found.IsValid() {
		t.Fatal("Invalid version information")
	}
}
