package ver

import (
	"fmt"
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
