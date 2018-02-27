package ver

import (
	"fmt"
)

type InputAction = func(param string) (*VersionInformation, error)

type InputSpec struct {
	Name        string
	Description string
	Parameter   string
	Action      InputAction
}

var AllInputs = []InputSpec{
	InputSpec{
		Name:        "version-info",
		Description: fmt.Sprintf("Read project and version information from an existing %v file in the current or parent directories.", VersionInfoYamlFilename),
		Action: func(param string) (*VersionInformation, error) {
			return TryReadFromVersionInfoYAML()
		},
	},
	InputSpec{
		Name:        "directory-name",
		Description: "Use the working directory's name as project name.",
		Action: func(param string) (*VersionInformation, error) {
			return TryReadFromWorkingDirectory()
		},
	},
	InputSpec{
		Name:        "build-host",
		Description: "Use the current machine's name and time as build host and timestamp.",
		Action: func(param string) (*VersionInformation, error) {
			return TryReadFromBuildHost()
		},
	},
	InputSpec{
		Name:        "npm",
		Description: fmt.Sprintf("Read project and version information from an existing %v file in the current directory.", PackageJsonFilename),
		Action: func(param string) (*VersionInformation, error) {
			return TryReadFromPackageJSON()
		},
	},
	InputSpec{
		Name:        "konsorten",
		Description: "Use marvin + konsorten default author information.",
		Action: func(param string) (*VersionInformation, error) {
			return TryReadFromKonsortenDefaults()
		},
	},
	InputSpec{
		Name:        "teamcity",
		Description: "Read version and revision information from TeamCity environment variables.",
		Action: func(param string) (*VersionInformation, error) {
			return TryReadFromTeamCity()
		},
	},
	InputSpec{
		Name:        "git",
		Description: "Read revision information from current directory's Git repository.",
		Action: func(param string) (*VersionInformation, error) {
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
