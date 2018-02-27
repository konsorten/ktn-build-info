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
}

type packageJSONAuthor struct {
	Author struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Website string `json:"url"`
	} `json:"author"`
}

type packageJSONAuthorOld struct {
	Author string `json:"author"`
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

	vi.Name = pj.Name
	vi.Description = pj.Description
	vi.SetSemVersion(pj.Version)

	// parse the package json file for author
	pja := packageJSONAuthor{}

	err = json.Unmarshal(data, &pja)

	if err == nil {
		vi.Author = pja.Author.Name
		vi.Email = pja.Author.Email
		vi.Website = pja.Author.Website
	} else {
		// try old author type
		pjo := packageJSONAuthorOld{}

		err = json.Unmarshal(data, &pjo)

		if err == nil {
			vi.Author = pjo.Author
		}
	}

	// done
	return vi, nil
}
