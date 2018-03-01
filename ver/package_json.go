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

func TryReadFromPackageJSON(ignoreName bool, ignoreVersion bool, ignoreDescription bool, ignoreAuthor bool, ignoreLicense bool) (*VersionInformation, error) {
	// check if the file exists
	if _, err := os.Stat(PackageJsonFilename); err != nil {
		log.Debugf("No %v found", PackageJsonFilename)
		return nil, nil
	}

	log.Debugf("Found version information: %v", PackageJsonFilename)

	// parse the package json file
	json, err := gabs.ParseJSONFile(PackageJsonFilename)

	if err != nil {
		log.Errorf("Failed to parse %v file: %v", PackageJsonFilename, err)
		return nil, err
	}

	// copy data
	vi := MakeVersionInformation()

	if c := json.Path("name"); !ignoreName && c != nil {
		vi.Name = c.Data().(string)
	}

	if c := json.Path("version"); !ignoreVersion && c != nil {
		vi.SetSemVersion(c.Data().(string))
	}

	if c := json.Path("description"); !ignoreDescription && c != nil {
		vi.Description = c.Data().(string)
	}

	if c := json.Path("license"); !ignoreLicense && c != nil {
		vi.License = c.Data().(string)
	}

	// parse author information
	if !ignoreAuthor {
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
	}

	// done
	return vi, nil
}

func UpdatePackageJSON(vi *VersionInformation) error {
	// parse the package json file
	json, err := gabs.ParseJSONFile(PackageJsonFilename)

	if err != nil {
		log.Errorf("Failed to parse %v file: %v", PackageJsonFilename, err)
		return err
	}

	// update the data
	_, err = json.SetP(vi.SemVerString(), "version")

	if err != nil {
		return err
	}

	_, err = json.SetP(vi.Name, "name")

	if err != nil {
		return err
	}

	_, err = json.SetP(vi.Description, "description")

	if err != nil {
		return err
	}

	_, err = json.SetP(vi.License, "license")

	if err != nil {
		return err
	}

	// drop old author information
	json.DeleteP("author")

	// create new author information
	_, err = json.SetP(vi.Author, "author.name")

	if err != nil {
		return err
	}

	_, err = json.SetP(vi.Email, "author.email")

	if err != nil {
		return err
	}

	_, err = json.SetP(vi.Website, "author.url")

	if err != nil {
		return err
	}

	// write back the file
	err = ioutil.WriteFile(PackageJsonFilename, json.BytesIndent("", "  "), os.FileMode(644))

	if err != nil {
		return err
	}

	// done
	log.Infof("Updated NPM %v file", PackageJsonFilename)

	return nil
}
