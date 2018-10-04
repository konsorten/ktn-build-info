package ver

import (
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestRunGitHere(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	// try read revision
	vi, err := TryReadFromGit()

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Revision: %v", vi.Revision)
}

func TestRunGitNoRepository(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	// change to test data dir
	currDir, _ := os.Getwd()
	os.Chdir(os.TempDir())
	defer os.Chdir(currDir)

	// try read revision
	vi, err := TryReadFromGit()

	if err != nil {
		t.Fatal(err)
	}

	if vi != nil {
		t.Fatal("No version information was expected")
	}
}
