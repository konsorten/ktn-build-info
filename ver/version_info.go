package ver

import (
	"fmt"
	"strings"

	"github.com/blang/semver"
)

type VersionInformation struct {
	Author      string
	Email       string
	Website     string
	Name        string
	Description string
	License     string

	Major  int
	Minor  int
	Hotfix int
	Build  int

	Revision string

	BuildTimestamp int
	BuildHost      string
}

func MakeVersionInformation() *VersionInformation {
	return &VersionInformation{
		Major:  -1,
		Minor:  -1,
		Hotfix: -1,
		Build:  0,
	}
}

func (vi *VersionInformation) String() string {
	return vi.VersionTitle()
}

func (vi *VersionInformation) IsValid() (ok bool, errors []string) {
	ok = true

	if vi.Author == "" {
		ok = false
		errors = append(errors, "Author is empty")
	}

	if vi.Email == "" {
		ok = false
		errors = append(errors, "E-mail is empty")
	}

	if vi.Website == "" {
		ok = false
		errors = append(errors, "Website is empty")
	}

	if vi.Name == "" {
		ok = false
		errors = append(errors, "Name is empty")
	}

	if vi.Description == "" {
		ok = false
		errors = append(errors, "Description is empty")
	}

	if vi.License == "" {
		ok = false
		errors = append(errors, "License is empty")
	}

	if vi.Major < 0 {
		ok = false
		errors = append(errors, "Major version < 0")
	}

	if vi.Minor < 0 {
		ok = false
		errors = append(errors, "Minor version < 0")
	}

	if vi.Hotfix < 0 {
		ok = false
		errors = append(errors, "Hotfix version < 0")
	}

	if vi.Build < 0 {
		ok = false
		errors = append(errors, "Build version < 0")
	}

	if vi.BuildTimestamp <= 0 {
		ok = false
		errors = append(errors, "Build timestamp is invalid")
	}

	if vi.BuildHost == "" {
		ok = false
		errors = append(errors, "Build Host is empty")
	}

	return
}

func (vi *VersionInformation) VersionString() string {
	return fmt.Sprintf("%v.%v.%v.%v", vi.Major, vi.Minor, vi.Hotfix, vi.Build)
}

func (vi *VersionInformation) VersionTitle() string {
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

func (vi *VersionInformation) SemVerString() string {
	parts := []string{fmt.Sprintf("%v.%v.%v", vi.Major, vi.Minor, vi.Hotfix)}

	// add additional parts
	if vi.Build > 0 {
		parts = append(parts, fmt.Sprintf("build+%v", vi.Build))
	}

	if vi.Revision != "" {
		parts = append(parts, fmt.Sprintf("rev.%v", vi.Revision))
	}

	return strings.Join(parts, "-")
}

func (vi *VersionInformation) SetSemVersion(semVerString string) {
	ver := semver.MustParse(semVerString)

	vi.Major = int(ver.Major)
	vi.Minor = int(ver.Minor)
	vi.Hotfix = int(ver.Patch)
}

func (vi *VersionInformation) CopyMissingFrom(copy *VersionInformation) {
	if copy == nil {
		return
	}

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

	if vi.License == "" && copy.License != "" {
		vi.License = copy.License
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
