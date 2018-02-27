package ver

import (
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type versionInfoActionsYAML struct {
	Inputs  []map[string]string `yaml:"inputs"`
	Outputs []map[string]string `yaml:"outputs"`
}

type VersionInfoAction struct {
	Name      string
	Parameter string
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
		for n, p := range i {
			ret.Inputs = append(ret.Inputs, VersionInfoAction{
				Name:      n,
				Parameter: p,
			})
		}
	}

	for _, o := range viy.Outputs {
		for n, p := range o {
			ret.Outputs = append(ret.Outputs, VersionInfoAction{
				Name:      n,
				Parameter: p,
			})
		}
	}

	return &ret, nil
}
