package ver

import (
	"fmt"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func IsTeamCityAvailable() bool {
	return os.Getenv("TEAMCITY_VERSION") != ""
}

func TryReadFromTeamCity(ignoreBuildNumber bool, ignoreRevision bool, ignoreName bool) (*VersionInformation, error) {
	if !IsTeamCityAvailable() {
		return nil, nil
	}

	vi := MakeVersionInformation()

	// read build number (from meta project)
	if !ignoreBuildNumber && vi.Build <= 0 {
		bn := os.Getenv("BUILDMETA_BUILD_NUMBER")

		if bn != "" {
			vi.Build, _ = strconv.Atoi(bn)

			log.Debugf("Found TeamCity build number: %v", vi.Build)
		}
	}

	// read build number
	if !ignoreBuildNumber && vi.Build <= 0 {
		bn := os.Getenv("BUILD_NUMBER")

		if bn != "" {
			vi.Build, _ = strconv.Atoi(bn)

			log.Debugf("Found TeamCity build number: %v", vi.Build)
		}
	}

	// read revision
	if !ignoreRevision && vi.Revision == "" {
		bn := os.Getenv("BUILD_VCS_NUMBER")

		if bn != "" {
			vi.Revision = bn

			log.Debugf("Found TeamCity VCS revision: %v", vi.Revision)
		}
	}

	// read project name
	if !ignoreName && vi.Name == "" {
		bn := os.Getenv("TEAMCITY_PROJECT_NAME")

		if bn != "" {
			vi.Name = bn

			log.Debugf("Found TeamCity project name: %v", vi.Name)
		}
	}

	return vi, nil
}

func (vi *VersionInformation) WriteToTeamCity() error {
	if !IsTeamCityAvailable() {
		return nil
	}

	// don't use logger, but write to console directly
	fmt.Printf("##teamcity[buildNumber '%v']", vi.VersionString())

	return nil
}
