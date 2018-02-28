package ver

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

func IsGitlabCIAvailable() bool {
	return os.Getenv("GITLAB_CI") == "true"
}

func TryReadFromGitlabCI() (*VersionInformation, error) {
	if !IsGitlabCIAvailable() {
		return nil, nil
	}

	vi := MakeVersionInformation()

	// read project name
	if vi.Name == "" {
		bn := os.Getenv("CI_PROJECT_NAME")

		if bn != "" {
			vi.Name = bn

			log.Debugf("Found Gitlab CI project name: %v", vi.Name)
		}
	}

	// read revision
	if vi.Revision == "" {
		bn := os.Getenv("CI_COMMIT_SHA")

		if bn != "" {
			vi.Revision = bn[0:6]

			log.Debugf("Found Gitlab CI VCS revision: %v", vi.Revision)
		}
	}

	return vi, nil
}