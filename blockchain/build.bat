@echo off
REM SkaffaCity Blockchain Build Script for Windows

echo Building SkaffaCity Blockchain...

REM Create bin directory
if not exist bin mkdir bin

REM Build for Linux (for deployment)
echo Building for Linux deployment...
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build -tags netgo -ldflags "-w -s" -o bin/skaffacityd ./cmd/skaffacityd

REM Build for Windows (for local development)
echo Building for Windows development...
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=0
go build -tags netgo -ldflags "-w -s" -o bin/skaffacityd.exe ./cmd/skaffacityd

echo Build complete!
dir bin\
