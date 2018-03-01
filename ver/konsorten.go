package ver

import (
	log "github.com/Sirupsen/logrus"
)

func TryReadFromKonsortenDefaults() (*VersionInformation, error) {
	vi := MakeVersionInformation()

	vi.Author = "marvin + konsorten GmbH"
	vi.Email = "open-source@konsorten.de"
	vi.Website = "http://www.konsorten.de"
	vi.License = "commercial"

	log.Debugf("Found author information from defaults: %v", vi.Author)
	log.Debugf("Found license information from defaults: %v", vi.License)

	return vi, nil
}
