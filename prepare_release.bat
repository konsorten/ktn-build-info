@echo off

cd /d "%~dp0"

REM ----- build the tool

go build

if %ERRORLEVEL% NEQ 0 (
    exit /b 1
)

REM ----- create its own build information

.\ktn-build-info.exe

if %ERRORLEVEL% NEQ 0 (
    exit /b 1
)

REM ----- rebuild the tool

go build

if %ERRORLEVEL% NEQ 0 (
    exit /b 1
)

REM ----- create help file

echo # Commandline help text: > HELP.md
echo ``` >> HELP.md
.\ktn-build-info.exe -h >> HELP.md 2>NUL
echo ``` >> HELP.md

exit /b 0
