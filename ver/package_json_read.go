package ver

import (
	"io/ioutil"
	"os"

	"github.com/Jeffail/gabs"
	log "github.com/Sirupsen/logrus"
)

const (
	PackageJsonFilename = "package.json"
)

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
	json, err := gabs.ParseJSON(data)

	if err != nil {
		log.Errorf("Failed to parse %v file: %v", PackageJsonFilename, err)
		return nil, err
	}

	// copy data
	vi := MakeVersionInformation()

	if c := json.Path("name"); c != nil {
		vi.Name = c.Data().(string)
	}

	if c := json.Path("version"); c != nil {
		vi.SetSemVersion(c.Data().(string))
	}

	if c := json.Path("description"); c != nil {
		vi.Description = c.Data().(string)
	}

	// parse author information
	if json.Exists("author", "name") {
		if c := json.Path("author.name"); c != nil {
			vi.Author = c.Data().(string)
		}

		if c := json.Path("author.email"); c != nil {
			vi.Email = c.Data().(string)
		}

		if c := json.Path("author.url"); c != nil {
			vi.Website = c.Data().(string)
		}
	} else {
		// parse old author info
		if c := json.Path("author"); c != nil {
			vi.Author = c.Data().(string)
		}
	}

	// done
	return vi, nil
}
