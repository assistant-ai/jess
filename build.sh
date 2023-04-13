#!/bin/bash
GOARCH=amd64
CGO_ENABLED=1

GOOS=windows go build -o jess-windows-amd64.exe
GOOS=linux go build -o jess-linux-amd64
GOOS=darwin go build -o jess-darwin-amd64