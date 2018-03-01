package ver

import (
	"io"
	"os"
	"testing"
)

func copyFile(t *testing.T, fromPath string, toPath string) {
	from, err := os.Open(fromPath)

	if err != nil {
		t.Fatal(err)
	}

	defer from.Close()

	to, err := os.Create(toPath)

	if err != nil {
		t.Fatal(err)
	}

	defer to.Close()

	_, err = io.Copy(to, from)

	if err != nil {
		t.Fatal(err)
	}
}

func createTestVersionInformationFromYAML(t *testing.T) *VersionInformation {
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
	found.Build = 4
	found.Revision = "abcdef0"

	if ok, _ := found.IsValid(); !ok {
		t.Fatal("Invalid version information")
	}

	return found
}
