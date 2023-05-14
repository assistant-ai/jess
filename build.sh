#!/bin/bash

VERSION=$(cat VERSION)

go mod tidy

gox -osarch="linux/amd64 darwin/amd64 darwin/arm64 windows/amd64" -ldflags="-X main.version=$VERSION" -output "jess-{{.OS}}-{{.Arch}}"
