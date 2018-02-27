package ver

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/Jeffail/gabs"
	log "github.com/Sirupsen/logrus"
)

func TestRun(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	// change to test data dir
	currDir, _ := os.Getwd()
	os.Chdir("examples/template")
	defer os.Chdir(currDir)

	// run
	err := RunCurrentDirectory()

	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	// read the json
	data, err := ioutil.ReadFile("test.json")

	if err != nil {
		t.Fatalf("Failed to read rendered template: %v", err)
	}

	// parse the output json file
	_, err = gabs.ParseJSON(data)

	if err != nil {
		t.Fatalf("Failed to parse rendered template: %v", err)
	}
}
