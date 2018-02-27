package main

import (
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	yaml "gopkg.in/yaml.v2"
)

const (
	versionInfoYamlFilename = "version-info.yml"
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

func tryReadVersionInfoYAML(c *cli.Context, vi *versionInformation) (bool, error) {
	// gather version information
	path := versionInfoYamlFilename
	var found []*versionInformation

	for {
		vii := makeVersionInformation()

		exists, err := tryReadVersionInfoYAMLInternal(c, vii, path)

		if err != nil {
			return exists, err
		}

		if !exists {
			break
		}

		log.Debugf("Found version information: %v", path)

		found = append([]*versionInformation{vii}, found...)
		path = fmt.Sprintf("../%v", path)
	}

	// nothing found?
	if len(found) <= 0 {
		return false, nil
	}

	// merge the results
	for _, vii := range found {
		vi.CopyMissingFrom(vii)
	}

	return true, nil
}

func tryReadVersionInfoYAMLInternal(c *cli.Context, vi *versionInformation, filename string) (bool, error) {
	// check if the file exists
	if _, err := os.Stat(filename); err != nil {
		log.Debugf("No %v found", filename)
		return false, nil
	}

	// read the file
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Errorf("Failed to read %v file: %v", filename, err)
		return false, err
	}

	// parse the file
	viy := versionInfoYAML{}

	err = yaml.Unmarshal(data, &viy)

	if err != nil {
		log.Errorf("Failed to parse %v file: %v", filename, err)
		return true, err
	}

	// copy data
	if viy.Version != "" {
		vi.SetSemVersion(viy.Version)
	}

	vi.Author = viy.Author.Name
	vi.Email = viy.Author.Email
	vi.Website = viy.Author.Website
	vi.Name = viy.Product.Name
	vi.Description = viy.Product.Description

	return true, nil
}
