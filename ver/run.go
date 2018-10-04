package ver

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

func RunCurrentDirectory() error {
	// read actions
	actions, err := ReadVersionInfoActions()

	if err != nil {
		return err
	}

	if actions == nil {
		return fmt.Errorf("No %v found", VersionInfoYamlFilename)
	}

	// nothing to do?
	if len(actions.Inputs) <= 0 {
		return fmt.Errorf("No input actions found")
	}

	if len(actions.Outputs) <= 0 {
		return fmt.Errorf("No output actions found")
	}

	// processing inputs
	result := MakeVersionInformation()

	for _, i := range actions.Inputs {
		log.Debugf("Processing input: %v [%v]", i.Name, i.Parameters)

		// resolve spec
		spec := GetInputSpec(i.Name)

		if spec == nil {
			return fmt.Errorf("Failed to resolve input action: %v", i.Name)
		}

		// run the action
		err := spec.Action(result, i.Parameters)

		if err != nil {
			return fmt.Errorf("Failed to run input action %v: %v", i.Name, err)
		}
	}

	// show information
	log.Infof("Version: %v", result)

	// validate
	if ok, errors := result.IsValid(); !ok {
		return fmt.Errorf("Build information is invalid: %v", strings.Join(errors, "; "))
	}

	// processing outputs
	for _, o := range actions.Outputs {
		log.Debugf("Processing output: %v [%v]", o.Name, o.Parameters)

		// resolve spec
		spec := GetOutputSpec(o.Name)

		if spec == nil {
			return fmt.Errorf("Failed to resolve output action: %v", o.Name)
		}

		// run the action
		err := spec.Action(result, o.Parameters)

		if err != nil {
			return fmt.Errorf("Failed to run output action %v [%v]: %v", o.Name, o.Parameters, err)
		}
	}

	return nil
}
