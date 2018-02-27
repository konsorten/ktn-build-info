package ver

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
)

const (
	PackageJsonFilename = "package.json"
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

func TryReadPackageJSON() (*VersionInformation, error) {
	// check if the file exists
	if _, err := os.Stat(PackageJsonFilename); err != nil {
		log.Debugf("No %v found", PackageJsonFilename)
		return nil, nil
	}

	// read the json
	data, err := ioutil.ReadFile(PackageJsonFilename)

	if err != nil {
		log.Errorf("Failed to read %v file: %v", PackageJsonFilename, err)
		return nil, err
	}

	log.Debugf("Found version information: %v", PackageJsonFilename)

	// parse the package json file
	pj := packageJSON{}

	err = json.Unmarshal(data, &pj)

	if err != nil {
		log.Errorf("Failed to parse %v file: %v", PackageJsonFilename, err)
		return nil, err
	}

	// copy data
	vi := MakeVersionInformation()

	vi.Author = pj.Author.Name
	vi.Email = pj.Author.Email
	vi.Website = pj.Author.Website
	vi.Name = pj.Name
	vi.Description = pj.Description
	vi.SetSemVersion(pj.Version)

	return vi, nil
}
