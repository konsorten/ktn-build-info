package ver

import (
	"io/ioutil"
	"regexp"

	log "github.com/Sirupsen/logrus"
)

func UpdateRegexFile(filename string, matchMap UpdateActions, usePosix bool, vi *VersionInformation) error {
	// read the file
	txt, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Errorf("Failed to read %v file: %v", filename, err)
		return err
	}

	// update the data
	for k, v := range matchMap {
		// compile the regular expression
		var r *regexp.Regexp

		if usePosix {
			r, err = regexp.CompilePOSIX(k)
		} else {
			r, err = regexp.Compile(k)
		}

		// render value
		log.Debugf("Applying regex: %v", v)

		newValue, err := RenderTemplate(v, filename, vi)

		if err != nil {
			return err
		}

		// replace
		txt = r.ReplaceAll(txt, []byte(newValue))
	}

	// write back the file
	err = ioutil.WriteFile(filename, txt, 0644)

	if err != nil {
		return err
	}

	// done
	log.Infof("Updated regex %v file", filename)

	return nil
}
