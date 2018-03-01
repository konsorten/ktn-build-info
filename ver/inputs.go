package ver

import (
	"fmt"
	"strconv"

	log "github.com/Sirupsen/logrus"
)

type InputAction = func(params map[string]string) (*VersionInformation, error)

type InputSpec struct {
	Name        string
	Description string
	Parameters  []string
	Action      InputAction
}

var AllInputs = []InputSpec{
	InputSpec{
		Name:        "version-info",
		Description: fmt.Sprintf("Read project and version information from an existing %v file in the current or parent directories.", VersionInfoYamlFilename),
		Parameters: []string{
			"depth:{levels}\tLimit depth of directories to scan. Set to 0 for current directory, only. Default is 10.",
		},
		Action: func(params map[string]string) (*VersionInformation, error) {
			// retrieve max depth
			maxDepth := 10

			if params["depth"] != "" {
				md, err := strconv.Atoi(params["depth"])

				if err != nil {
					log.Errorf("Failed to parse 'maxDepth' parameter: %v", params["depth"])
					return nil, err
				}

				maxDepth = md
			}

			return TryReadFromVersionInfoYAML(maxDepth)
		},
	},
	InputSpec{
		Name:        "directory-name",
		Description: "Use the working directory's name as project name.",
		Action: func(params map[string]string) (*VersionInformation, error) {
			return TryReadFromWorkingDirectory()
		},
	},
	InputSpec{
		Name:        "build-host",
		Description: "Use the current machine's name and time as build host and timestamp.",
		Action: func(params map[string]string) (*VersionInformation, error) {
			return TryReadFromBuildHost()
		},
	},
	InputSpec{
		Name:        "npm",
		Description: fmt.Sprintf("Read project and version information from an existing %v file in the current directory.", PackageJsonFilename),
		Parameters: []string{
			"name:false\tIgnore the project name.",
			"desc:false\tIgnore the project description.",
			"license:false\tIgnore the project license information.",
			"version:false\tIgnore the project version number.",
			"author:false\tIgnore the project author information.",
		},
		Action: func(params map[string]string) (*VersionInformation, error) {
			return TryReadFromPackageJSON(
				params["name"] == "false",
				params["version"] == "false",
				params["desc"] == "false",
				params["author"] == "false",
				params["license"] == "false",
			)
		},
	},
	InputSpec{
		Name:        "konsorten",
		Description: "Use marvin + konsorten default author information.",
		Action: func(params map[string]string) (*VersionInformation, error) {
			return TryReadFromKonsortenDefaults()
		},
	},
	InputSpec{
		Name:        "teamcity",
		Description: "Read project name, version, and revision information from TeamCity environment variables.",
		Parameters: []string{
			"build:false\tIgnore the build number.",
			"rev:false\tIgnore the revision number.",
			"name:false\tIgnore the project name.",
		},
		Action: func(params map[string]string) (*VersionInformation, error) {
			return TryReadFromTeamCity(
				params["build"] == "false",
				params["rev"] == "false",
				params["name"] == "false",
			)
		},
	},
	InputSpec{
		Name:        "gitlab-ci",
		Description: "Read project name and revision information from GitLab CI environment variables.",
		Parameters: []string{
			"rev:false\tIgnore the revision number.",
			"name:false\tIgnore the project name.",
		},
		Action: func(params map[string]string) (*VersionInformation, error) {
			return TryReadFromGitlabCI(
				params["rev"] == "false",
				params["name"] == "false",
			)
		},
	},
	InputSpec{
		Name:        "git",
		Description: "Read revision information from current directory's Git repository.",
		Action: func(params map[string]string) (*VersionInformation, error) {
			return TryReadFromGit()
		},
	},
}

func GetInputSpec(name string) *InputSpec {
	for _, s := range AllInputs {
		if s.Name == name {
			return &s
		}
	}

	return nil
}
