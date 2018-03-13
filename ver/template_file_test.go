package ver

import (
	"os"
	"testing"

	"github.com/Jeffail/gabs"
	log "github.com/Sirupsen/logrus"
)

func TestRenderTemplateFile(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	// change to test data dir
	currDir, _ := os.Getwd()
	os.Chdir("examples/template")
	defer os.Chdir(currDir)

	found := createTestVersionInformationFromYAML(t)

	// render template
	err := found.WriteTemplateFile("test.json.template", 0644)

	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	// read the json
	_, err = gabs.ParseJSONFile("test.json")

	if err != nil {
		t.Fatalf("Failed to read rendered template: %v", err)
	}
}
