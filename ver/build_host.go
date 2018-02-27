package ver

import (
	"os"
	"path/filepath"
	"time"
)

func TryReadFromBuildHost() (*VersionInformation, error) {
	cwd, err := os.Hostname()

	if err != nil {
		return nil, err
	}

	// done
	vi := MakeVersionInformation()

	vi.BuildHost = filepath.Base(cwd)
	vi.BuildTimestamp = time.Now().Unix()

	return vi, nil
}
