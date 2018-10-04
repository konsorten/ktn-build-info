package ver

import (
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type versionInfoActionsYAML struct {
	Inputs  []map[string]string `yaml:"inputs"`
	Outputs []map[string]string `yaml:"outputs"`
}

type VersionInfoAction struct {
	Name       string
	Parameters map[string]string
}

type VersionInfoActions struct {
	Inputs  []VersionInfoAction
	Outputs []VersionInfoAction
}

func ReadVersionInfoActions() (*VersionInfoActions, error) {
	// check if the file exists
	if _, err := os.Stat(VersionInfoYamlFilename); err != nil {
		log.Debugf("No %v found", VersionInfoYamlFilename)
		return nil, nil
	}

	// read the file
	data, err := ioutil.ReadFile(VersionInfoYamlFilename)

	if err != nil {
		log.Errorf("Failed to read %v file: %v", VersionInfoYamlFilename, err)
		return nil, err
	}

	// parse the file
	viy := versionInfoActionsYAML{}

	err = yaml.Unmarshal(data, &viy)

	if err != nil {
		log.Errorf("Failed to parse %v file: %v", VersionInfoYamlFilename, err)
		return nil, err
	}

	// convert the information
	var ret VersionInfoActions

	for _, i := range viy.Inputs {
		action := VersionInfoAction{
			Parameters: make(map[string]string),
		}

		for n, p := range i {
			if n == "action" {
				action.Name = p
			} else {
				action.Parameters[n] = p
			}
		}

		// is action valid?
		if action.Name == "" {
			return nil, fmt.Errorf("Failed to get action name; missing 'action' property: %v", i)
		}

		// done
		ret.Inputs = append(ret.Inputs, action)
	}

	for _, o := range viy.Outputs {
		action := VersionInfoAction{
			Parameters: make(map[string]string),
		}

		for n, p := range o {
			if n == "action" {
				action.Name = p
			} else {
				action.Parameters[n] = p
			}
		}

		// is action valid?
		if action.Name == "" {
			return nil, fmt.Errorf("Failed to get action name; missing 'action' property: %v", o)
		}

		// done
		ret.Outputs = append(ret.Outputs, action)
	}

	return &ret, nil
}
