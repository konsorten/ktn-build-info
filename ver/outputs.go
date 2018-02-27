package ver

type OutputAction = func(vi *VersionInformation, param string) error

type OutputSpec struct {
	Name        string
	Description string
	Parameter   string
	Action      OutputAction
}

var AllOutputs = []OutputSpec{
	OutputSpec{
		Name:        "template",
		Description: "Renders a template file and writes the result into a file dropping the last extension, e.g. myfile.c.template becomes myfile.c. Takes the relative file path as parameter.",
		Parameter:   "file",
		Action: func(vi *VersionInformation, param string) error {
			return vi.WriteTemplateFile(param)
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
