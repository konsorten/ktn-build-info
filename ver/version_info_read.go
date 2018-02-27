package ver

import (
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

const (
	VersionInfoYamlFilename = "version-info.yml"
)

type versionInfoYAML struct {
	Version string `yaml:"version"`

	Product struct {
		Name        string `yaml:"name"`
		Description string `yaml:"desc"`
	} `yaml:"product"`

	Author struct {
		Name    string `yaml:"name"`
		Email   string `yaml:"email"`
		Website string `yaml:"url"`
	} `yaml:"author"`
}

func TryReadVersionInfoYAML() (*VersionInformation, error) {
	// gather version information
	path := VersionInfoYamlFilename
	var found []*VersionInformation

	for {
		vii, err := tryReadVersionInfoYAMLInternal(path)

		if err != nil {
			return nil, err
		}

		if vii == nil {
			break
		}

		log.Debugf("Found version information: %v", path)

		found = append([]*VersionInformation{vii}, found...)
		path = fmt.Sprintf("../%v", path)
	}

	// nothing found?
	if len(found) <= 0 {
		return nil, nil
	}

	// merge the results
	vi := MakeVersionInformation()

	for _, vii := range found {
		vi.CopyMissingFrom(vii)
	}

	return vi, nil
}

func tryReadVersionInfoYAMLInternal(filename string) (*VersionInformation, error) {
	// check if the file exists
	if _, err := os.Stat(filename); err != nil {
		log.Debugf("No %v found", filename)
		return nil, nil
	}

	// read the file
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Errorf("Failed to read %v file: %v", filename, err)
		return nil, err
	}

	// parse the file
	viy := versionInfoYAML{}

	err = yaml.Unmarshal(data, &viy)

	if err != nil {
		log.Errorf("Failed to parse %v file: %v", filename, err)
		return nil, err
	}

	// copy data
	vi := MakeVersionInformation()

	if viy.Version != "" {
		vi.SetSemVersion(viy.Version)
	}

	vi.Author = viy.Author.Name
	vi.Email = viy.Author.Email
	vi.Website = viy.Author.Website
	vi.Name = viy.Product.Name
	vi.Description = viy.Product.Description

	return vi, nil
}
