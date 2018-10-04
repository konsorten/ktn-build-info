package ver

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func IsGitlabCIAvailable() bool {
	return os.Getenv("GITLAB_CI") == "true"
}

func TryReadFromGitlabCI(ignoreRevision bool, ignoreProjectName bool) (*VersionInformation, error) {
	if !IsGitlabCIAvailable() {
		return nil, nil
	}

	vi := MakeVersionInformation()

	// read project name
	if !ignoreProjectName && vi.Name == "" {
		bn := os.Getenv("CI_PROJECT_NAME")

		if bn != "" {
			vi.Name = bn

			log.Debugf("Found Gitlab CI project name: %v", vi.Name)
		}
	}

	// read revision
	if !ignoreRevision && vi.Revision == "" {
		bn := os.Getenv("CI_COMMIT_SHA")

		if bn != "" {
			vi.Revision = bn

			log.Debugf("Found Gitlab CI VCS revision: %v", vi.Revision)
		}
	}

	return vi, nil
}
