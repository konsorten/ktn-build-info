package ver

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func TryReadFromWorkingDirectory() (*VersionInformation, error) {
	cwd, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	// done
	vi := MakeVersionInformation()

	vi.Name = filepath.Base(cwd)

	log.Debugf("Found project name (from working directory): %v", vi.Name)

	return vi, nil
}
