# Commandline help text: 
``` 
Incorrect Usage. flag: help requested

NAME:
   ktn-build-info.exe - Build Information Tool

    This tool creates build version information files
    based on the available information and writes
    them to different file formats.
    
    If available, it will read from existing files:
      * version-info.yml
      * package.json (NPM)

    It looks for the version-info.yml file in the specified
    and parent directories.

USAGE:
   ktn-build-info.exe [global options] [arguments...]

VERSION:
   1.0.5.19

AUTHOR:
   marvin + konsorten GmbH

GLOBAL OPTIONS:
   --directory value, -C value  The directory in which to run. Using the current as default. (default: ".")
   --debug, -d                  Enable debug output.
   --version, -v                print the version

INPUTS:
  version-info              Read project and version information from an existing version-info.yml file in the current or parent directories.
    depth:{levels}          Limit depth of directories to scan. Set to 0 for current directory, only. Default is 10.
  directory-name            Use the working directory's name as project name.
  build-host                Use the current machine's name and time as build host and timestamp.
  consul-build-id           Retrieves a build number based on the build revision. Use this for non-numeric revisions, like Git.
    url:{consulurl}         The connection URL to consul, e.g. http://consul:8500/dc1 or http://:token@consul:8500/dc1.
    root:{key}              The base KV path for the project, e.g. builds/myproject.
  npm                       Read project and version information from an existing package.json file in the current directory.
    name:false              Ignore the project name.
    desc:false              Ignore the project description.
    license:false           Ignore the project license information.
    version:false           Ignore the project version number.
    author:false            Ignore the project author information.
  konsorten                 Use marvin + konsorten default author information.
  teamcity                  Read project name, version, and revision information from TeamCity environment variables.
    build:false             Ignore the build number.
    rev:false               Ignore the revision number.
    name:false              Ignore the project name.
  gitlab-ci                 Read project name and revision information from GitLab CI environment variables.
    rev:false               Ignore the revision number.
    name:false              Ignore the project name.
  git                       Read revision information from current directory's Git repository.
  env-vars                  Read build information from environment variables.
    name:{envVarName}       Read the project name from the specified environment variable {envVarName}. Ignored if missing or empty.
    desc:{envVarName}       Read the project description from the specified environment variable {envVarName}. Ignored if missing or empty.
    rev:{envVarName}        Read the build revision from the specified environment variable {envVarName}. Ignored if missing or empty.
    build:{envVarName}      Read the build ID from the specified environment variable {envVarName}. Ignored if missing or empty.
    buildHost:{envVarName}  Read the build host name from the specified environment variable {envVarName}. Ignored if missing or empty.
  limit-revision            Limits the length of the revision string.
    length:{int}            The number of characters to limit the revision string, e.g. 7 for Git short revision.

OUTPUTS:
  template              Renders a template file and writes the result into a file dropping the last extension, e.g. myfile.c.template becomes myfile.c. Takes the relative file path as parameter.
    file:{path}         Path to the template file.
    mode:{filemode}     File mode to use for the file (Linux and macOS). Typical values are 644 for a regular file and 755 for an executable file.
  update-npm            Updates an existing NPM package.json file in the current directory.
  update-json           Updates an existing JSON document.
    file:{path}         Path to the JSON file.
    indent:{chars}      Characters to use for indent, like four space characters.
    !{path}:{value}     Set the value at {path} to {value}. All template file fields are supported. Strings have to be quoted, e.g. '"{$.Author$}"' or '{$.Author | asQuoted$}'. Missing objects are created automatically.
    !{path}:$null$      Sets the value at {path} to null.
    !{path}:$delete$    Deletes the value at {path}.
  update-xml            Updates an existing XML document.
    file:{path}         Path to the XML file.
    indent:{chars}      Characters to use for indent, like four space characters.
    !{xpath}:{value}    Set the value at {xpath} to {value}. The target element/attribute must exist. All template file fields are supported.
    !{xpath}:$create$   Creates a new element or attribute from {xpath}. The parent element must exist.
    !{xpath}:$ensure$   If missing, creates a new element or attribute from {xpath}. The parent element must exist.
    !{xpath}:$delete$   Deletes the value at {xpath}.
  update-regex          Updates an existing text document.
    file:{path}         Path to the text file.
    posix:true          Restricts the regular expression to POSIX ERE (egrep) syntax.
    !{regex}:{replace}  Replaces all matches of {regex} with {replace}. All template file fields are supported.
  go-syso               Writes the version information to be embedded into an Go executable (version_info.syso). Supported on Windows, only.
  teamcity              Writes the version number back to TeamCity.

TEMPLATE FILE FIELDS:
  {$.Author$}         The name of the author/manufacturer, e.g. the company name.
  {$.BuildDate$}      The current timestamp as Go time.Time object, usage: {$.BuildDate.Year$}
  {$.BuildDateISO$}   The current ISO timestamp, example: 2018-02-27T17:28:58Z
  {$.BuildDateUnix$}  The current epoch unix timestamp, example: 1519752538
  {$.BuildHost$}      The name of the current machine, example: wks005
  {$.Description$}    The description of the project.
  {$.Email$}          The support e-mail address for the product.
  {$.License$}        The license which applies to the project.
  {$.Name$}           The name/ID of the project.
  {$.Revision$}       The source code revision, example: b862de
  {$.SemVer$}         The semantic version number as A.B.C-build+D-rev.XXXXX with four parts, example: 1.2.3-build+3-rev.b862de
  {$.Version$}        The version number as A.B.C.D with four parts, example: 1.2.3.4
  {$.VersionParts$}   The version number as array with four elements, usage: {$index .VersionParts 0$} for major version number.
  {$.VersionTitle$}   The version number as descriptive text, example: 1.2.3 (build: 4, rev:b862de)
  {$.Website$}        The website of the author or the product.

TEMPLATE FUNCTIONS:
  asBool      Convert the value to a boolean value. Use this to convert strings or numbers to true/false.
  asInt       Convert the value to an integer number. Use this to convert enumerations to their respective numeric value.
  asQuoted    Convert the value to a quoted string. Use this to convert values to a quoted string. Will escape existing quotes to '\"'.
  asString    Convert the value to a string. Use this to convert complex objects to their string representation.
  encodeHtml  Encode the value to be HTML compatible, e.g. '&' becomes '&quot;'.
  encodeUrl   Encode the value to be URL compatible, e.g. '&' becomes '%26'.
  encodeXml   Encode the value to be XML compatible, e.g. '&' becomes '&quot;'.
  envVar      Reads the specified value from an environment variable, e.g. {$envVar "TEMP"$}.

  See for more details: https://golang.org/pkg/text/template/

``` 
