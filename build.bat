@echo off

:: Variables
set "ParentDir=%~dp0"
set "BuildDir=%ParentDir%build"
set "SourceDir=%ParentDir%"

if not exist "%BuildDir%" (
    echo [*] Creating %BuildDir%
    mkdir "%BuildDir%"
)

pushd %BuildDir%

go build "%SourceDir%

go test -c "%SourceDir%

popd
