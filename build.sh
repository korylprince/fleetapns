#!/bin/bash

mkdir -p build

GOOS=linux GOARCH=amd64 go build -trimpath -o build/fleetapns_linux_amd64
GOOS=windows GOARCH=amd64 go build -trimpath -o build/fleetapns_windows_amd64
GOOS=darwin GOARCH=amd64 go build -trimpath -o build/fleetapns_darwin_amd64
GOOS=darwin GOARCH=arm64 go build -trimpath -o build/fleetapns_darwin_arm64
