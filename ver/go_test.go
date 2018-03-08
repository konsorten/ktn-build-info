package ver

import (
	"os"
	"runtime"
	"testing"

	log "github.com/Sirupsen/logrus"
)

func TestWriteGoSysoFiles(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	// this feature works on windows, only
	if runtime.GOOS != "windows" {
		t.Logf("Skipping test, supported on Windows, only")
		return
	}

	// change to test data dir
	currDir, _ := os.Getwd()
	os.Chdir("examples/go-build-test")
	defer os.Chdir(currDir)

	// get version
	vi := createTestVersionInformationFromYAML(t)

	// write files
	err := WriteGoSysoFile(vi)

	if err != nil {
		t.Fatalf("Failed to write files: %v", err)
	}
}
