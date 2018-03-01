package ver

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/konsorten/go-xmldom"
)

func isParamAction(input string, action string, param *string) bool {
	prefix := fmt.Sprintf("$%v:", action)

	if !strings.HasPrefix(input, prefix) {
		return false
	}

	if !strings.HasSuffix(input, "$") {
		return false
	}

	*param = string(input[len(prefix) : len(input)-1])

	return true
}

func extractXPathDestination(xpath *string, elemName *string, attrName *string) bool {
	*elemName = ""
	*attrName = ""

	// perform split
	pathParts := strings.Split(*xpath, "/")

	if len(pathParts) <= 0 {
		return false
	}

	// check for valid entries
	dest := pathParts[len(pathParts)-1]

	if dest == "." {
		return false
	}

	// rebuild path
	*xpath = strings.Join(pathParts[:len(pathParts)-1], "/")

	// handle root-only path
	if *xpath == "" {
		*xpath = "/"
	}

	// handle destination
	if strings.HasPrefix(dest, "@") {
		*attrName = dest[1:]
	} else {
		*elemName = dest
	}

	return true
}

func UpdateXmlFile(filePath string, updates UpdateActions, vi *VersionInformation, indent string) error {
	// parse the xml file
	doc, err := xmldom.ParseFile(filePath)

	if err != nil {
		log.Errorf("Failed to parse %v file: %v", filePath, err)
		return err
	}

	root := doc.Root

	// apply updates
	for _, u := range orderUpdates(updates) {
		update := u.Update
		path := u.Path
		var elemName, attrName string

		switch {
		case update == "$delete$":
			if extractXPathDestination(&path, &elemName, &attrName) {
				// handle simple xpath
				if elemName != "" {
					path = fmt.Sprintf("%v/%v", path, elemName)

					for _, elem := range root.Query(path) {
						elem.Parent.RemoveChild(elem)

						log.Debugf("Deleted element: %v", path)
					}
				} else {
					for _, elem := range root.Query(path) {
						elem.RemoveAttribute(attrName)

						log.Debugf("Deleted attribute: %v/@%v", path, attrName)
					}
				}
			} else {
				// delete complex xpath
				for _, elem := range root.Query(path) {
					elem.Parent.RemoveChild(elem)

					log.Debugf("Deleted element: %v", path)
				}
			}
		case update == "$create$":
			if !extractXPathDestination(&path, &elemName, &attrName) || elemName == "" {
				return fmt.Errorf("Invalid path syntax: %v; expected to end with '/element'", path)
			}

			for _, elem := range root.Query(path) {
				elem.CreateNode(elemName)

				if strings.HasSuffix(path, "/") {
					log.Debugf("Created element: %v%v", path, elemName)
				} else {
					log.Debugf("Created element: %v/%v", path, elemName)
				}
			}
		case update == "$ensure$":
			if !extractXPathDestination(&path, &elemName, &attrName) || elemName == "" {
				return fmt.Errorf("Invalid path syntax: %v; expected to end with '/element'", path)
			}

			for _, elem := range root.Query(path) {
				if elem.GetChild(elemName) == nil {
					elem.CreateNode(elemName)

					if strings.HasSuffix(path, "/") {
						log.Debugf("Created element: %v%v", path, elemName)
					} else {
						log.Debugf("Created element: %v/%v", path, elemName)
					}
				}
			}
		default:
			// render value
			newValue, err := RenderTemplate(update, path, vi)

			if err != nil {
				return err
			}

			// extract destination
			if !extractXPathDestination(&path, &elemName, &attrName) {
				// update based on xpath
				found := root.Query(path)

				if len(found) <= 0 {
					log.Warnf("No element found to update: %v (missing $create$ action?) in %v", path, filePath)
					break
				}

				// update the value
				for _, elem := range found {
					elem.Text = newValue

					log.Debugf("Updating element %v: %v", path, newValue)
				}
			} else {
				// handle update
				if elemName != "" {
					// find elements to update
					path = fmt.Sprintf("%v/%v", path, elemName)

					found := root.Query(path)

					if len(found) <= 0 {
						log.Warnf("No element found to update: %v (missing $create$ action?) in %v", path, filePath)
						break
					}

					// update the value
					for _, elem := range found {
						elem.Text = newValue

						log.Debugf("Updating element %v: %v", path, newValue)
					}
				} else {
					// find attributes to update
					found := root.Query(path)

					if len(found) <= 0 {
						log.Warnf("No parent element found to update: %v/@%v (missing $create$ action?) in %v", path, attrName, filePath)
						break
					}

					// update the value
					for _, elem := range found {
						elem.SetAttributeValue(attrName, newValue)

						log.Debugf("Updating attribute %v/@%v: %v", path, attrName, newValue)
					}
				}
			}
		}
	}

	// write back the file
	var newContent string

	if indent == "" {
		newContent = doc.XML()
	} else {
		newContent = doc.XMLPrettyEx(indent)
	}

	return ioutil.WriteFile(filePath, []byte(newContent), os.FileMode(644))
}
