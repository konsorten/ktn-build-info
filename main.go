package main

import (
	"fmt"
	"os"

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
	app.Usage = fmt.Sprintf(`Build Information Tool

    This tool creates build version information files
    based on the available information and writes
    them to different file formats.

    It looks for the version-info.yml file in
    the specified and parent directories.
    
    If available, it will read from existing files:
        * %v
        * package.json (NPM)`, ver.VersionInfoYamlFilename)

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

	return nil
}
