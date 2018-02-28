package ver

import (
	"bytes"
	"fmt"
	"html"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"text/template"
	"time"

	log "github.com/Sirupsen/logrus"
)

const (
	TemplateFileFieldPrefix = "{$"
	TemplateFileFieldSuffix = "$}"
)

type templateFileData struct {
	Author        string    `doc:"The name of the author/manufacturer, e.g. the company name."`
	Website       string    `doc:"The website of the author or the product."`
	Email         string    `doc:"The support e-mail address for the product."`
	Version       string    `doc:"The version number as A.B.C.D with four parts, example: 1.2.3.4"`
	VersionTitle  string    `doc:"The version number as descriptive text, example: 1.2.3 (build: 4, rev:b862de)"`
	VersionParts  [4]string `doc:"The version number as array with four elements, usage: {$index .VersionParts 0$} for major version number."`
	SemVer        string    `doc:"The semantic version number as A.B.C-build+D-rev.XXXXX with four parts, example: 1.2.3-build+3-rev.b862de"`
	Revision      string    `doc:"The source code revision, example: b862de"`
	BuildHost     string    `doc:"The name of the current machine, example: wks005"`
	BuildDateUnix string    `doc:"The current epoch unix timestamp, example: 1519752538"`
	BuildDateISO  string    `doc:"The current ISO timestamp, example: 2018-02-27T17:28:58Z"`
	BuildDate     time.Time `doc:"The current timestamp as Go time.Time object, usage: {$.BuildDate.Year$}"`
}

func GetTemplateFileFields() map[string]string {
	m := make(map[string]string)

	s := reflect.TypeOf(templateFileData{})

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)

		m[f.Name] = f.Tag.Get("doc")
	}

	return m
}

func RenderTemplate(templateContent string, templateName string, vi *VersionInformation) (string, error) {
	// build data object
	timestamp := time.Unix(int64(vi.BuildTimestamp), 0).UTC()

	data := templateFileData{
		Author:        vi.Author,
		Website:       vi.Website,
		Email:         vi.Email,
		Version:       vi.VersionString(),
		VersionTitle:  vi.VersionTitle(),
		VersionParts:  [4]string{strconv.Itoa(vi.Major), strconv.Itoa(vi.Minor), strconv.Itoa(vi.Hotfix), strconv.Itoa(vi.Build)},
		SemVer:        vi.SemVerString(),
		Revision:      vi.Revision,
		BuildHost:     vi.BuildHost,
		BuildDateUnix: strconv.Itoa(vi.BuildTimestamp),
		BuildDateISO:  timestamp.Format("2006-01-02T15:04:05Z"),
		BuildDate:     timestamp,
	}

	// template functions
	templateFuncs := template.FuncMap{
		"asInt": func(value interface{}) (ret string) {
			switch v := value.(type) {
			case time.Month:
				ret = strconv.Itoa(int(v))
			default:
				ret = fmt.Sprintf("%v", v)
			}
			return
		},
		"encodeHtml": func(value interface{}) (ret string) {
			switch v := value.(type) {
			case string:
				ret = html.EscapeString(v)
			default:
				ret = html.EscapeString(fmt.Sprintf("%v", v))
			}
			return
		},
		"encodeXml": func(value interface{}) (ret string) {
			switch v := value.(type) {
			case string:
				ret = html.EscapeString(v)
			default:
				ret = html.EscapeString(fmt.Sprintf("%v", v))
			}
			return
		},
		"encodeUrl": func(value interface{}) (ret string) {
			switch v := value.(type) {
			case string:
				ret = url.QueryEscape(v)
			default:
				ret = url.QueryEscape(fmt.Sprintf("%v", v))
			}
			return
		},
	}

	// parse the template
	templ, err := template.New(templateName).Delims(TemplateFileFieldPrefix, TemplateFileFieldSuffix).Funcs(templateFuncs).Parse(templateContent)

	if err != nil {
		return "", err
	}

	// render template
	b := bytes.NewBufferString("")

	err = templ.ExecuteTemplate(b, templateName, &data)

	if err != nil {
		return "", err
	}

	// done
	return b.String(), nil
}

func (vi *VersionInformation) WriteTemplateFile(templateFilePath string) error {
	log.Debugf("Rendering template: %v", templateFilePath)

	// read the template
	tc, err := ioutil.ReadFile(templateFilePath)

	if err != nil {
		return err
	}

	// render the template
	filename := filepath.Base(templateFilePath)

	result, err := RenderTemplate(string(tc), filename, vi)

	if err != nil {
		return err
	}

	// get target filename
	var outputFilename string

	if true {
		p := strings.Split(filename, ".")

		if len(p) < 2 {
			return fmt.Errorf("Filename has no extension to drop: %v", filename)
		}

		outputFilename = strings.Join(p[:len(p)-1], ".")
	}

	// write the file
	outputPath := filepath.Join(filepath.Dir(templateFilePath), outputFilename)

	return ioutil.WriteFile(outputPath, []byte(result), os.FileMode(644))
}
