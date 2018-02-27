package ver

import (
	log "github.com/Sirupsen/logrus"
)

func TryReadFromKonsortenDefaults() (*VersionInformation, error) {
	vi := MakeVersionInformation()

	vi.Author = "marvin + konsorten GmbH"
	vi.Email = "open-source@konsorten.de"
	vi.Website = "http://www.konsorten.de"

	log.Debugf("Found author information from defaults: %v", vi.Author)

	return vi, nil
}
