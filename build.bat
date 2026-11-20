@echo off
setlocal EnableExtensions EnableDelayedExpansion

set "PACKAGE=.\cmd\cli"
set "OUTPUT_DIR=dist"
set "VERSION="

for /f "tokens=2 delims==" %%V in ('findstr /r /c:"version = " "%PACKAGE%\main.go"') do (
    set "VERSION=%%V"
)

if not defined VERSION (
    echo Unable to determine the application version from %PACKAGE%\main.go.
    exit /b 1
)

set "VERSION=!VERSION: =!"
set "VERSION=!VERSION:"=!"

if not exist "%OUTPUT_DIR%" mkdir "%OUTPUT_DIR%"

echo Building bert version !VERSION!

for %%O in (windows linux darwin) do (
    for %%A in (amd64 arm64) do (
        set "EXT="
        if "%%O"=="windows" set "EXT=.exe"

        echo Building %%O/%%A...
        set "GOOS=%%O"
        set "GOARCH=%%A"
        go build -o "%OUTPUT_DIR%\bert-!VERSION!-%%O-%%A!EXT!" "%PACKAGE%"
        if errorlevel 1 exit /b 1
    )
)

echo Build complete.
endlocal
