#!/bin/bash

VERSION=$(cat VERSION)
git tag -a "v$VERSION" -m "Release version $VERSION"
git push --tags