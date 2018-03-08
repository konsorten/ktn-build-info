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
	app.Version = "1.0.2"
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
			b.WriteString("\t")
			b.WriteString(i.Description)
			b.WriteString("\n")

			if i.Parameters != nil {
				for _, p := range i.Parameters {
					b.WriteString("    ")
					b.WriteString(p)
					b.WriteString("\n")
				}
			}
		}

		// add outputs
		b.WriteString("\nOUTPUTS:\n")

		for _, o := range ver.AllOutputs {
			b.WriteString("  ")
			b.WriteString(o.Name)
			b.WriteString("\t")
			b.WriteString(o.Description)
			b.WriteString("\n")

			if o.Parameters != nil {
				for _, p := range o.Parameters {
					b.WriteString("    ")
					b.WriteString(p)
					b.WriteString("\n")
				}
			}
		}

		// add template file fields
		b.WriteString("\nTEMPLATE FILE FIELDS:\n")

		for _, d := range ver.GetTemplateFileFields() {
			b.WriteString("  {{\"")
			b.WriteString(fmt.Sprintf("%v.%v%v", ver.TemplateFileFieldPrefix, d.Name, ver.TemplateFileFieldSuffix))
			b.WriteString("\"}}\t")
			b.WriteString(d.Description)
			b.WriteString("\n")
		}

		b.WriteString("\nTEMPLATE FUNCTIONS:\n")

		for _, d := range ver.GetTemplateFileFunctions() {
			b.WriteString("  ")
			b.WriteString(d.Name)
			b.WriteString("\t")
			b.WriteString(d.Description)
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
