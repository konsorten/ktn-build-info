package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

const (
	packageJsonFilename = "package.json"
)

type packageJSON struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Author      struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Website string `json:"url"`
	} `json:"author"`
}

func tryReadPackageJSON(c *cli.Context, vi *versionInformation) (bool, error) {
	// check if the file exists
	if _, err := os.Stat(packageJsonFilename); err != nil {
		log.Debugf("No %v found", packageJsonFilename)
		return false, nil
	}

	// read the json
	data, err := ioutil.ReadFile(packageJsonFilename)

	if err != nil {
		log.Errorf("Failed to read %v file: %v", packageJsonFilename, err)
		return false, err
	}

	log.Debugf("Found version information: %v", packageJsonFilename)

	// parse the package json file
	pj := packageJSON{}

	err = json.Unmarshal(data, &pj)

	if err != nil {
		log.Errorf("Failed to parse %v file: %v", packageJsonFilename, err)
		return true, err
	}

	// copy data
	vi.Author = pj.Author.Name
	vi.Email = pj.Author.Email
	vi.Website = pj.Author.Website
	vi.Name = pj.Name
	vi.Description = pj.Description
	vi.SetSemVersion(pj.Version)

	return true, nil
}
