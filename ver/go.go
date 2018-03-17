package ver

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/josephspurrier/goversioninfo"
)

const (
	goSysoFilename = "version_info.syso"
)

func WriteGoSysoFile(vi *VersionInformation) error {
	// this feature works on windows, only
	if runtime.GOOS != "windows" {
		log.Warnf("Creating Go executable version information is supported on Windows, only")
		return nil
	}

	// create the manifest file
	manifestFile := filepath.Join(os.TempDir(), fmt.Sprintf("ktn-build-info-%v.manifest", time.Now().UnixNano()))
	manifest, err := getGoManifestXml(vi)

	if err != nil {
		return fmt.Errorf("Failed to generate manifest file: %v", err)
	}

	err = ioutil.WriteFile(manifestFile, []byte(manifest), 0644)

	if err != nil {
		return fmt.Errorf("Failed to write manifest file: %v", err)
	}

	defer os.Remove(manifestFile)

	// copy values
	gvi := goversioninfo.VersionInfo{
		FixedFileInfo: goversioninfo.FixedFileInfo{
			FileVersion: goversioninfo.FileVersion{
				Major: vi.Major,
				Minor: vi.Minor,
				Patch: vi.Hotfix,
				Build: vi.Build,
			},
			FileFlagsMask: "3f",
			FileFlags:     "00",
			FileOS:        "040004",
			FileType:      "01",
			FileSubType:   "00",
		},
		StringFileInfo: goversioninfo.StringFileInfo{
			CompanyName:     vi.Author,
			ProductName:     vi.Name,
			ProductVersion:  vi.VersionTitle(),
			FileDescription: vi.Description,
			LegalCopyright:  fmt.Sprintf("Copyright %v %v. License: %v", time.Now().Year(), vi.Author, vi.License),
		},
		VarFileInfo: goversioninfo.VarFileInfo{
			Translation: goversioninfo.Translation{
				LangID:    goversioninfo.LngUSEnglish,
				CharsetID: goversioninfo.CsUnicode,
			},
		},
		ManifestPath: manifestFile,
	}

	// prepare
	gvi.Build()
	gvi.Walk()

	// write the syso file
	err = gvi.WriteSyso(goSysoFilename, runtime.GOARCH)

	if err != nil {
		return fmt.Errorf("Failed to generate Go syso file: %v", err)
	}

	return nil
}

func getGoManifestXml(vi *VersionInformation) (string, error) {
	manifest :=
		`<?xml version="1.0" encoding="utf-8"?>
<assembly manifestVersion="1.0" xmlns="urn:schemas-microsoft-com:asm.v1">

	<trustInfo xmlns="urn:schemas-microsoft-com:asm.v2">
		<security>
			<requestedPrivileges xmlns="urn:schemas-microsoft-com:asm.v3">
				<requestedExecutionLevel level="asInvoker" uiAccess="false" />
			</requestedPrivileges>
		</security>
	</trustInfo>

	<compatibility xmlns="urn:schemas-microsoft-com:compatibility.v1">
		<application>
			<!-- Windows Vista -->
			<supportedOS Id="{e2011457-1546-43c5-a5fe-008deee3d3f0}" />

			<!-- Windows 7 -->
			<supportedOS Id="{35138b9a-5d96-4fbd-8e2d-a2440225f93a}" />

			<!-- Windows 8 -->
			<supportedOS Id="{4a2f28e3-53b9-4441-ba9c-d69d4a4a6e38}" />

			<!-- Windows 8.1 -->
			<supportedOS Id="{1f676c76-80e1-4239-95bb-83d0f6d0da78}" />

			<!-- Windows 10 -->
			<supportedOS Id="{8e0f7a12-bfb3-4fe8-b9a5-48fd50a15a9a}" />
		
			<windowsSettings xmlns:ws2="http://schemas.microsoft.com/SMI/2016/WindowsSettings">
				<ws2:longPathAware>true</ws2:longPathAware>
			</windowsSettings>
			
		</application>
	</compatibility>

</assembly>`

	return manifest, nil
}
