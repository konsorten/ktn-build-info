package ver

import (
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
)

func TryReadFromEnvironmentVariables(varMap map[string]string) (*VersionInformation, error) {
	vi := MakeVersionInformation()

	// project name
	if varMap["name"] != "" {
		v := os.Getenv(varMap["name"])

		if v != "" {
			log.Debugf("Found project name in environment variable %v: %v", varMap["rev"], v)

			vi.Name = v
		}
	}

	// project description
	if varMap["desc"] != "" {
		v := os.Getenv(varMap["desc"])

		if v != "" {
			log.Debugf("Found project description in environment variable %v: %v", varMap["rev"], v)

			vi.Description = v
		}
	}

	// revision
	if varMap["rev"] != "" {
		v := os.Getenv(varMap["rev"])

		if v != "" {
			log.Debugf("Found revision information in environment variable %v: %v", varMap["rev"], v)

			vi.Revision = v
		}
	}

	// build id
	if varMap["build"] != "" {
		v := os.Getenv(varMap["build"])

		if v != "" {
			i, err := strconv.Atoi(v)

			if err != nil {
				log.Fatalf("Failed to parse build ID in environment variable %v: %v", varMap["rev"], v)
				return nil, err
			}

			log.Debugf("Found build ID in environment variable %v: %v", varMap["rev"], v)

			vi.Build = i
		}
	}

	// build host
	if varMap["buildHost"] != "" {
		v := os.Getenv(varMap["buildHost"])

		if v != "" {
			log.Debugf("Found build host name in environment variable %v: %v", varMap["rev"], v)

			vi.BuildHost = v
		}
	}

	// done
	return vi, nil
}
