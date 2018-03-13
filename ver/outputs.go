package ver

import (
	"fmt"
	"os"
	"strconv"
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
			"mode:{filemode}\tFile mode to use for the file (Linux and macOS). Typical values are 644 for a regular file and 755 for an executable file.",
		},
		Action: func(vi *VersionInformation, params map[string]string) error {
			// determine filename
			f := params["file"]

			if f == "" {
				return fmt.Errorf("No template file specified; missing 'file' parameter")
			}

			// parse the file mode
			var fm os.FileMode = 0644

			fmStr := params["mode"]

			if fmStr != "" {
				fmParsed, err := strconv.ParseUint(fmStr, 8, 32)

				if err != nil {
					return fmt.Errorf("Failed to parse file mode: %v: %v", fmStr, err)
				}

				fm = os.FileMode(fmParsed)
			}

			return vi.WriteTemplateFile(f, fm)
		},
	},
	OutputSpec{
		Name:        "update-npm",
		Description: fmt.Sprintf("Updates an existing NPM %v file in the current directory.", PackageJsonFilename),
		Action: func(vi *VersionInformation, params map[string]string) error {
			return UpdatePackageJSON(vi)
		},
	},
	OutputSpec{
		Name:        "update-json",
		Description: "Updates an existing JSON document.",
		Parameters: []string{
			"file:{path}\tPath to the JSON file.",
			"indent:{chars}\tCharacters to use for indent, like four space characters.",
			"!{path}:{value}\tSet the value at {path} to {value}. All template file fields are supported. Strings have to be quoted, e.g. '\"{$.Author$}\"' or '{$.Author | asQuoted$}'. Missing objects are created automatically.",
			"!{path}:$null$\tSets the value at {path} to null.",
			"!{path}:$delete$\tDeletes the value at {path}.",
		},
		Action: func(vi *VersionInformation, params map[string]string) error {
			f := params["file"]

			if f == "" {
				return fmt.Errorf("No JSON file specified; missing 'file' parameter")
			}

			// gather actions
			actions := make(UpdateActions)

			for k, v := range params {
				if strings.HasPrefix(k, "!") {
					actions[k[1:]] = v
				}
			}

			return UpdateJsonFile(f, actions, vi, params["indent"])
		},
	},
	OutputSpec{
		Name:        "update-xml",
		Description: "Updates an existing XML document.",
		Parameters: []string{
			"file:{path}\tPath to the XML file.",
			"indent:{chars}\tCharacters to use for indent, like four space characters.",
			"!{xpath}:{value}\tSet the value at {xpath} to {value}. The target element/attribute must exist. All template file fields are supported.",
			"!{xpath}:$create$\tCreates a new element or attribute from {xpath}. The parent element must exist.",
			"!{xpath}:$ensure$\tIf missing, creates a new element or attribute from {xpath}. The parent element must exist.",
			"!{xpath}:$delete$\tDeletes the value at {xpath}.",
		},
		Action: func(vi *VersionInformation, params map[string]string) error {
			f := params["file"]

			if f == "" {
				return fmt.Errorf("No XML file specified; missing 'file' parameter")
			}

			// gather actions
			actions := make(UpdateActions)

			for k, v := range params {
				if strings.HasPrefix(k, "!") {
					actions[k[1:]] = v
				}
			}

			return UpdateXmlFile(f, actions, vi, params["indent"])
		},
	},
	OutputSpec{
		Name:        "update-regex",
		Description: "Updates an existing text document.",
		Parameters: []string{
			"file:{path}\tPath to the text file.",
			"posix:true\tRestricts the regular expression to POSIX ERE (egrep) syntax.",
			"!{regex}:{replace}\tReplaces all matches of {regex} with {replace}. All template file fields are supported.",
		},
		Action: func(vi *VersionInformation, params map[string]string) error {
			f := params["file"]

			if f == "" {
				return fmt.Errorf("No text file specified; missing 'file' parameter")
			}

			// gather actions
			actions := make(UpdateActions)

			for k, v := range params {
				if strings.HasPrefix(k, "!") {
					actions[k[1:]] = v
				}
			}

			return UpdateRegexFile(f, actions, params["posix"] == "true", vi)
		},
	},
	OutputSpec{
		Name:        "go-syso",
		Description: fmt.Sprintf("Writes the version information to be embedded into an Go executable (%v). Supported on Windows, only.", goSysoFilename),
		Action: func(vi *VersionInformation, params map[string]string) error {
			return WriteGoSysoFile(vi)
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
