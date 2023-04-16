#!/bin/bash

build_app() {
  local GOOS=$1
  local GOARCH=$2
  local filename=$3
  
  CGO_ENABLED=1 GOOS="$GOOS" GOARCH="$GOARCH" go build -ldflags="-X main.version=$VERSION" -o "$filename"
}

VERSION=$(cat VERSION)

build_app "linux"   "amd64" "jess-linux-amd64"
build_app "darwin"  "amd64" "jess-darwin-amd64"
build_app "darwin"  "arm64" "jess-darwin-arm64"
build_app "windows" "amd64" "jess-windows-amd64.exe"