package ver

import (
	"os"
	"testing"

	"github.com/Jeffail/gabs"
	log "github.com/sirupsen/logrus"
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
	found, err := TryReadFromPackageJSON(false, false, false, false, false)

	if err != nil {
		t.Fatalf("Failed to read version info: %v", err)
	}

	if found == nil {
		t.Fatalf("%v not found", PackageJsonFilename)
	}

	found.CopyMissingFrom(host)

	if ok, _ := found.IsValid(); !ok {
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
	found, err := TryReadFromPackageJSON(false, false, false, false, false)

	if err != nil {
		t.Fatalf("Failed to read version info: %v", err)
	}

	if found == nil {
		t.Fatalf("%v not found", PackageJsonFilename)
	}

	found.CopyMissingFrom(host)
	found.CopyMissingFrom(defaults)

	if ok, _ := found.IsValid(); !ok {
		t.Fatal("Invalid version information")
	}
}

func TestUpdatePackageJson(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	// change to test data dir
	currDir, _ := os.Getwd()
	os.Chdir("examples/npm-update")
	defer os.Chdir(currDir)

	copyFile(t, "package.json.test", PackageJsonFilename)

	found := createTestVersionInformationFromYAML(t)

	// update package json
	err := UpdatePackageJSON(found)

	if err != nil {
		t.Fatalf("Failed to update NPM package.json: %v", err)
	}

	// parse the output json file
	_, err = gabs.ParseJSONFile(PackageJsonFilename)

	if err != nil {
		t.Fatalf("Failed to read updated NPM package.json: %v", err)
	}
}
