package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/blang/semver"
)

type versionInformation struct {
	Author      string
	Email       string
	Website     string
	Name        string
	Description string

	Major  int
	Minor  int
	Hotfix int
	Build  int

	Revision string

	BuildTimestamp int64
	BuildHost      string
}

func makeVersionInformation() *versionInformation {
	hostname, _ := os.Hostname()

	return &versionInformation{
		Major:          -1,
		Minor:          -1,
		Hotfix:         -1,
		Build:          0,
		BuildTimestamp: time.Now().Unix(),
		BuildHost:      hostname,
	}
}

func (vi *versionInformation) String() string {
	return vi.VersionTitle()
}

func (vi *versionInformation) IsValid() bool {
	return vi.Author != "" &&
		vi.Email != "" &&
		vi.Website != "" &&
		vi.Name != "" &&
		vi.Description != "" &&
		vi.Major >= 0 &&
		vi.Minor >= 0 &&
		vi.Hotfix >= 0 &&
		vi.Build >= 0 &&
		vi.BuildTimestamp > 0 &&
		vi.BuildHost != ""
}

func (vi *versionInformation) VersionString() string {
	return fmt.Sprintf("%v.%v.%v.%v", vi.Major, vi.Minor, vi.Hotfix, vi.Build)
}

func (vi *versionInformation) VersionTitle() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("%v.%v.%v", vi.Major, vi.Minor, vi.Hotfix))

	// add additional parts
	var parts []string

	if vi.Build > 0 {
		parts = append(parts, fmt.Sprintf("build: %v", vi.Build))
	}

	if vi.Revision != "" {
		parts = append(parts, fmt.Sprintf("rev: %v", vi.Revision))
	}

	// append parts
	if len(parts) > 0 {
		b.WriteString(" (")
		b.WriteString(strings.Join(parts, ", "))
		b.WriteString(")")
	}

	return b.String()
}

func (vi *versionInformation) SemVerString() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("%v.%v.%v", vi.Major, vi.Minor, vi.Hotfix))

	// add additional parts
	var parts []string

	if vi.Build > 0 {
		parts = append(parts, fmt.Sprintf("build+%v", vi.Build))
	}

	if vi.Revision != "" {
		parts = append(parts, fmt.Sprintf("rev.%v", vi.Revision))
	}

	// append parts
	if len(parts) > 0 {
		b.WriteString(" (")
		b.WriteString(strings.Join(parts, ", "))
		b.WriteString(")")
	}

	return b.String()
}

func (vi *versionInformation) SetSemVersion(semVerString string) {
	ver := semver.MustParse(semVerString)

	vi.Major = int(ver.Major)
	vi.Minor = int(ver.Minor)
	vi.Hotfix = int(ver.Patch)
}

func (vi *versionInformation) CopyMissingFrom(copy *versionInformation) {
	if vi.Author == "" && copy.Author != "" {
		vi.Author = copy.Author
	}

	if vi.Build < 0 && copy.Build >= 0 {
		vi.Build = copy.Build
	}

	if vi.BuildHost == "" && copy.BuildHost != "" {
		vi.BuildHost = copy.BuildHost
	}

	if vi.BuildTimestamp <= 0 && copy.BuildTimestamp > 0 {
		vi.BuildTimestamp = copy.BuildTimestamp
	}

	if vi.Description == "" && copy.Description != "" {
		vi.Description = copy.Description
	}

	if vi.Email == "" && copy.Email != "" {
		vi.Email = copy.Email
	}

	if vi.Hotfix < 0 && copy.Hotfix >= 0 {
		vi.Hotfix = copy.Hotfix
	}

	if vi.Major < 0 && copy.Major >= 0 {
		vi.Major = copy.Major
	}

	if vi.Minor < 0 && copy.Minor >= 0 {
		vi.Minor = copy.Minor
	}

	if vi.Name == "" && copy.Name != "" {
		vi.Name = copy.Name
	}

	if vi.Revision == "" && copy.Revision != "" {
		vi.Revision = copy.Revision
	}

	if vi.Website == "" && copy.Website != "" {
		vi.Website = copy.Website
	}
}
