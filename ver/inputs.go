package ver

import (
	"fmt"
)

type InputSpec struct {
	Name         string
	Description  string
	HasParameter bool
}

var AllInputs = []InputSpec{
	InputSpec{
		Name:        "version-info",
		Description: fmt.Sprintf("Read project and version information from an existing %v file in the current or parent directories.", VersionInfoYamlFilename),
	},
	InputSpec{
		Name:        "directory-name",
		Description: "Use the working directory's name as project name.",
	},
	InputSpec{
		Name:        "npm",
		Description: fmt.Sprintf("Read project and version information from an existing %v file in the current directory.", PackageJsonFilename),
	},
	InputSpec{
		Name:        "konsorten",
		Description: "Use marvin + konsorten default author information.",
	},
	InputSpec{
		Name:        "teamcity",
		Description: "Read version and revision information from TeamCity environment variables.",
	},
	InputSpec{
		Name:        "git",
		Description: "Read revision information from current directory's Git repository.",
	},
}
