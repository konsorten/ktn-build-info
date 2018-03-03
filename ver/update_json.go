package ver

import (
	"bytes"
	jsonImpl "encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/Jeffail/gabs"
	log "github.com/Sirupsen/logrus"
)

type UpdateActions map[string]string

func UpdateJsonFile(filePath string, updates UpdateActions, vi *VersionInformation, indent string) error {
	// parse the json file
	json, err := gabs.ParseJSONFile(filePath)

	if err != nil {
		log.Errorf("Failed to parse %v file: %v", filePath, err)
		return err
	}

	// apply updates
	for _, u := range orderUpdates(updates) {
		update := u.Update
		path := u.Path

		switch update {
		case "$delete$":
			if json.ExistsP(path) {
				err = json.DeleteP(path)

				if err != nil {
					return err
				}

				log.Debugf("Deleted element: %v", path)
			}
		case "$null$":
			log.Debugf("Updating element %v: null", path)

			_, err = json.SetP(nil, path)

			if err != nil {
				return err
			}
		default:
			// render value
			newValue, err := RenderTemplate(update, path, vi)

			if err != nil {
				return err
			}

			// parse the value
			log.Debugf("Updating element %v: %v", path, newValue)

			var newValueObj interface{}

			if true {
				newValueBuf := bytes.NewBufferString(newValue)
				decoder := jsonImpl.NewDecoder(newValueBuf)

				token, err := decoder.Token()

				if err != nil {
					log.Errorf("Failed to parse JSON value '%v': %v", newValue, err)
					return err
				}

				switch v := token.(type) {
				case jsonImpl.Delim:
					return fmt.Errorf("Failed to parse JSON value '%v': got delimiter", newValue)
				default:
					newValueObj = v
				}
			}

			// update the value
			_, err = json.SetP(newValueObj, path)

			if err != nil {
				return err
			}
		}
	}

	// write back the file
	var newContent []byte

	if indent == "" {
		newContent = json.Bytes()
	} else {
		newContent = json.BytesIndent("", indent)
	}

	return ioutil.WriteFile(filePath, newContent, 0644)
}
