package main

import (
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/konsorten/ktn-build-info/ver"
)

func main() {
	exit := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	app := cli.NewApp()
	app.Author = "marvin + konsorten GmbH"
	app.Version = "0.1.0"
	app.HideHelp = true
	app.Usage = fmt.Sprintf(
		`Build Information Tool

		This tool creates build version information files
		based on the available information and writes
		them to different file formats.
		
		If available, it will read from existing files:
			* %v
			* %v (NPM)

		It looks for the %v file in the specified
		and parent directories.`,
		ver.VersionInfoYamlFilename,
		ver.PackageJsonFilename,
		ver.VersionInfoYamlFilename,
	)

	// generate app help templates
	if true {
		var b strings.Builder

		b.WriteString(cli.AppHelpTemplate)

		// add inputs
		b.WriteString("\nINPUTS:\n")

		for _, i := range ver.AllInputs {
			b.WriteString("  ")
			b.WriteString(i.Name)

			if i.Parameter != "" {
				b.WriteString(" {")
				b.WriteString(i.Parameter)
				b.WriteString("}")
			}

			b.WriteString("\t")
			b.WriteString(i.Description)
			b.WriteString("\n")
		}

		// add outputs
		b.WriteString("\nOUTPUTS:\n")

		for _, i := range ver.AllOutputs {
			b.WriteString("  ")
			b.WriteString(i.Name)

			if i.Parameter != "" {
				b.WriteString(" {")
				b.WriteString(i.Parameter)
				b.WriteString("}")
			}

			b.WriteString("\t")
			b.WriteString(i.Description)
			b.WriteString("\n")
		}

		// add template file fields
		b.WriteString("\nTEMPLATE FILE FIELDS:\n")

		for n, d := range ver.GetTemplateFileFields() {
			b.WriteString("  {{\"")
			b.WriteString(fmt.Sprintf("%v.%v%v", ver.TemplateFileFieldPrefix, n, ver.TemplateFileFieldSuffix))
			b.WriteString("\"}}\t")
			b.WriteString(d)
			b.WriteString("\n")
		}

		b.WriteString("\n  See for more details: https://golang.org/pkg/text/template/\n\n")

		cli.AppHelpTemplate = b.String()
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "directory, C",
			Value: ".",
			Usage: "The directory in which to run. Using the current as default.",
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "Enable debug output.",
		},
	}

	app.Action = func(c *cli.Context) {
		exit(run(c))
	}

	exit(app.Run(os.Args))
}

func run(c *cli.Context) error {
	// enable debug output
	if c.Bool("debug") {
		log.SetLevel(log.DebugLevel)
	}

	// change working directory
	dir := c.String("directory")

	if err := os.Chdir(dir); err != nil {
		return fmt.Errorf("Failed to change to directory %s: %v", dir, err)
	}

	log.Debugf("Working directory: %v", dir)

	// start scanning for version information
	return ver.RunCurrentDirectory()
}
