package ver

import (
	"fmt"
	"strconv"

	log "github.com/Sirupsen/logrus"
)

type InputAction = func(vi *VersionInformation, params map[string]string) error

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
		Action: func(vi *VersionInformation, params map[string]string) error {
			// retrieve max depth
			maxDepth := 10

			if params["depth"] != "" {
				md, err := strconv.Atoi(params["depth"])

				if err != nil {
					log.Errorf("Failed to parse 'maxDepth' parameter: %v", params["depth"])
					return err
				}

				maxDepth = md
			}

			// run
			ver, err := TryReadFromVersionInfoYAML(maxDepth)

			if err != nil {
				return err
			}

			vi.CopyMissingFrom(ver)

			return nil
		},
	},
	InputSpec{
		Name:        "directory-name",
		Description: "Use the working directory's name as project name.",
		Action: func(vi *VersionInformation, params map[string]string) error {
			// run
			ver, err := TryReadFromWorkingDirectory()

			if err != nil {
				return err
			}

			vi.CopyMissingFrom(ver)

			return nil
		},
	},
	InputSpec{
		Name:        "build-host",
		Description: "Use the current machine's name and time as build host and timestamp.",
		Action: func(vi *VersionInformation, params map[string]string) error {
			// run
			ver, err := TryReadFromBuildHost()

			if err != nil {
				return err
			}

			vi.CopyMissingFrom(ver)

			return nil
		},
	},
	InputSpec{
		Name:        "consul-build-id",
		Description: "Retrieves a build number based on the build revision. Use this for non-numeric revisions, like Git.",
		Parameters: []string{
			"url:{consulurl}\tThe connection URL to consul, e.g. http://consul:8500/dc1 or http://:token@consul:8500/dc1.",
			"root:{key}\tThe base KV path for the project, e.g. builds/myproject.",
		},
		Action: func(vi *VersionInformation, params map[string]string) error {
			return RetrieveBuildFromConsul(params["url"], params["root"], vi)
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
		Action: func(vi *VersionInformation, params map[string]string) error {
			// run
			ver, err := TryReadFromPackageJSON(
				params["name"] == "false",
				params["version"] == "false",
				params["desc"] == "false",
				params["author"] == "false",
				params["license"] == "false",
			)

			if err != nil {
				return err
			}

			vi.CopyMissingFrom(ver)

			return nil
		},
	},
	InputSpec{
		Name:        "konsorten",
		Description: "Use marvin + konsorten default author information.",
		Action: func(vi *VersionInformation, params map[string]string) error {
			// run
			ver, err := TryReadFromKonsortenDefaults()

			if err != nil {
				return err
			}

			vi.CopyMissingFrom(ver)

			return nil
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
		Action: func(vi *VersionInformation, params map[string]string) error {
			// run
			ver, err := TryReadFromTeamCity(
				params["build"] == "false",
				params["rev"] == "false",
				params["name"] == "false",
			)

			if err != nil {
				return err
			}

			vi.CopyMissingFrom(ver)

			return nil
		},
	},
	InputSpec{
		Name:        "gitlab-ci",
		Description: "Read project name and revision information from GitLab CI environment variables.",
		Parameters: []string{
			"rev:false\tIgnore the revision number.",
			"name:false\tIgnore the project name.",
		},
		Action: func(vi *VersionInformation, params map[string]string) error {
			// run
			ver, err := TryReadFromGitlabCI(
				params["rev"] == "false",
				params["name"] == "false",
			)

			if err != nil {
				return err
			}

			vi.CopyMissingFrom(ver)

			return nil
		},
	},
	InputSpec{
		Name:        "git",
		Description: "Read revision information from current directory's Git repository.",
		Action: func(vi *VersionInformation, params map[string]string) error {
			// run
			ver, err := TryReadFromGit()

			if err != nil {
				return err
			}

			vi.CopyMissingFrom(ver)

			return nil
		},
	},
	InputSpec{
		Name:        "env-vars",
		Description: "Read build information from environment variables.",
		Parameters: []string{
			"name:{envVarName}\tRead the project name from the specified environment variable {envVarName}. Ignored if missing or empty.",
			"desc:{envVarName}\tRead the project description from the specified environment variable {envVarName}. Ignored if missing or empty.",
			"rev:{envVarName}\tRead the build revision from the specified environment variable {envVarName}. Ignored if missing or empty.",
			"build:{envVarName}\tRead the build ID from the specified environment variable {envVarName}. Ignored if missing or empty.",
			"buildHost:{envVarName}\tRead the build host name from the specified environment variable {envVarName}. Ignored if missing or empty.",
		},
		Action: func(vi *VersionInformation, params map[string]string) error {
			// run
			ver, err := TryReadFromEnvironmentVariables(params)

			if err != nil {
				return err
			}

			vi.CopyMissingFrom(ver)

			return nil
		},
	},
	InputSpec{
		Name:        "limit-revision",
		Description: "Limits the length of the revision string.",
		Parameters: []string{
			"length:{int}\tThe number of characters to limit the revision string, e.g. 7 for Git short revision.",
		},
		Action: func(vi *VersionInformation, params map[string]string) error {
			l, err := strconv.Atoi(params["length"])

			if err != nil {
				return err
			}

			vi.LimitRevision(l)

			return nil
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
