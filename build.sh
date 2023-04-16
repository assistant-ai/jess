#!/bin/bash
CGO_ENABLED=1

GOOS=windows GOARCH=amd64 go build -o jess-windows-amd64.exe
GOOS=linux   GOARCH=amd64 go build -o jess-linux-amd64
GOOS=darwin  GOARCH=amd64 go build -o jess-darwin-amd64
GOOS=darwin  GOARCH=arm64 go build -o jess-darwin-arm64