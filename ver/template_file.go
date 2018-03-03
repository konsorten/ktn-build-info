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
	Name          string    `doc:"The name/ID of the project."`
	Description   string    `doc:"The description of the project."`
	License       string    `doc:"The license which applies to the project."`
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

func GetTemplateFileFields() ItemDocumentation {
	var m ItemDocumentation

	s := reflect.TypeOf(templateFileData{})

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)

		m = append(m, ItemDocumentationEntry{f.Name, f.Tag.Get("doc")})
	}

	m.Sort()

	return m
}

type templateFileFunctions struct {
	AsInt      interface{} `doc:"Convert the value to an integer number. Use this to convert enumerations to their respective numeric value."`
	AsString   interface{} `doc:"Convert the value to a string. Use this to convert complex objects to their string representation."`
	AsBool     interface{} `doc:"Convert the value to a boolean value. Use this to convert strings or numbers to true/false."`
	AsQuoted   interface{} `doc:"Convert the value to a quoted string. Use this to convert values to a quoted string. Will escape existing quotes to '\\\"'."`
	EncodeHtml interface{} `doc:"Encode the value to be HTML compatible, e.g. '&' becomes '&quot;'."`
	EncodeXml  interface{} `doc:"Encode the value to be XML compatible, e.g. '&' becomes '&quot;'."`
	EncodeUrl  interface{} `doc:"Encode the value to be URL compatible, e.g. '&' becomes '%26'."`
	EnvVar     interface{} `doc:"Reads the specified value from an environment variable, e.g. {$envVar \"TEMP\"$}."`
}

var templateFileFunctionsImpl = templateFileFunctions{
	AsInt: func(value interface{}) (ret int) {
		switch v := value.(type) {
		case time.Month:
			ret = int(v)
		case string:
			ret, _ = strconv.Atoi(v)
		default:
			ret, _ = strconv.Atoi(fmt.Sprintf("%v", v))
		}
		return
	},

	AsString: func(value interface{}) (ret string) {
		switch v := value.(type) {
		case string:
			ret = v
		default:
			ret = fmt.Sprintf("%v", v)
		}
		return
	},

	AsBool: func(value interface{}) (ret bool) {
		switch v := value.(type) {
		case int:
			ret = v != 0
		case float32:
			ret = v != 0
		case float64:
			ret = v != 0
		case string:
			ret, _ = strconv.ParseBool(v)
		default:
			ret, _ = strconv.ParseBool(fmt.Sprintf("%v", v))
		}
		return
	},

	AsQuoted: func(value interface{}) (ret string) {
		switch v := value.(type) {
		case string:
			ret = v
		default:
			ret = fmt.Sprintf("%v", v)
		}

		ret = strings.Replace(ret, "\\", "\\\\", -1)
		ret = strings.Replace(ret, "\"", "\\\"", -1)

		ret = fmt.Sprintf("\"%v\"", ret)
		return
	},

	EncodeHtml: func(value interface{}) (ret string) {
		switch v := value.(type) {
		case string:
			ret = html.EscapeString(v)
		default:
			ret = html.EscapeString(fmt.Sprintf("%v", v))
		}
		return
	},

	EncodeXml: func(value interface{}) (ret string) {
		switch v := value.(type) {
		case string:
			ret = html.EscapeString(v)
		default:
			ret = html.EscapeString(fmt.Sprintf("%v", v))
		}
		return
	},

	EncodeUrl: func(value interface{}) (ret string) {
		switch v := value.(type) {
		case string:
			ret = url.QueryEscape(v)
		default:
			ret = url.QueryEscape(fmt.Sprintf("%v", v))
		}
		return
	},

	EnvVar: func(value string) (ret string) {
		return os.Getenv(value)
	},
}

func GetTemplateFileFunctionsMap() template.FuncMap {
	m := make(template.FuncMap)

	s := reflect.TypeOf(templateFileFunctionsImpl)
	v := reflect.ValueOf(templateFileFunctionsImpl)

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)

		m[strings.ToLower(f.Name[:1])+f.Name[1:]] = v.Field(i).Interface()
	}

	return m
}

func GetTemplateFileFunctions() ItemDocumentation {
	var m ItemDocumentation

	s := reflect.TypeOf(templateFileFunctionsImpl)

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)

		m = append(m, ItemDocumentationEntry{strings.ToLower(f.Name[:1]) + f.Name[1:], f.Tag.Get("doc")})
	}

	m.Sort()

	return m
}

func RenderTemplate(templateContent string, templateName string, vi *VersionInformation) (string, error) {
	// build data object
	timestamp := time.Unix(int64(vi.BuildTimestamp), 0).UTC()

	data := templateFileData{
		Name:          vi.Name,
		Description:   vi.Description,
		License:       vi.License,
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

	// parse the template
	templ, err := template.New(templateName).Delims(TemplateFileFieldPrefix, TemplateFileFieldSuffix).Funcs(GetTemplateFileFunctionsMap()).Parse(templateContent)

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

	return ioutil.WriteFile(outputPath, []byte(result), 0644)
}
