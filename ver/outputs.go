package ver

import (
	"fmt"
	"strings"
)

type OutputAction = func(vi *VersionInformation, params map[string]string) error

type OutputSpec struct {
	Name        string
	Description string
	Parameters  []string
	Action      OutputAction
}

var AllOutputs = []OutputSpec{
	OutputSpec{
		Name:        "template",
		Description: "Renders a template file and writes the result into a file dropping the last extension, e.g. myfile.c.template becomes myfile.c. Takes the relative file path as parameter.",
		Parameters: []string{
			"file:{path}\tPath to the template file.",
		},
		Action: func(vi *VersionInformation, params map[string]string) error {
			f := params["file"]

			if f == "" {
				return fmt.Errorf("No template file specified; missing 'file' parameter")
			}

			return vi.WriteTemplateFile(f)
		},
	},
	OutputSpec{
		Name:        "update-json",
		Description: "Updates an existing JSON document.",
		Parameters: []string{
			"file:{path}\tPath to the JSON file.",
			"indent:{chars}\tCharacters to use for indent, like four space characters.",
			"!{path}:{value}\tSet the value at {path} to {value}. All template file fields are supported. Strings have to be quoted, e.g. '\"{$.Author$}\"'",
			"!{path}:$null$\tSets the value at {path} to null.",
			"!{path}:$delete$\tDeletes the value at {path}.",
		},
		Action: func(vi *VersionInformation, params map[string]string) error {
			f := params["file"]

			if f == "" {
				return fmt.Errorf("No JSON file specified; missing 'file' parameter")
			}

			// gather actions
			actions := make(map[string]string)

			for k, v := range params {
				if strings.HasPrefix(k, "!") {
					actions[k[1:]] = v
				}
			}

			return UpdateJsonFile(f, actions, vi, params["indent"])
		},
	},
	OutputSpec{
		Name:        "teamcity",
		Description: "Writes the version number back to TeamCity.",
		Action: func(vi *VersionInformation, params map[string]string) error {
			return vi.WriteToTeamCity()
		},
	},
}

func GetOutputSpec(name string) *OutputSpec {
	for _, s := range AllOutputs {
		if s.Name == name {
			return &s
		}
	}

	return nil
}
