#!/bin/bash

VERSION=$(cat VERSION)
git tag -a "v$VERSION" -m "Release version $VERSION"
git push --tags

# Upload binary jess to the GitHub release
gh release create "v$VERSION" --title "Release version $VERSION" --notes "Added logic to upload binary jess."

upload_binary() {
  local path=$1
  gh release upload "v$VERSION" "${path}" --label "binary-jess"
}
upload_binary "./jess-darwin-arm64"
upload_binary "./jess-darwin-amd64"
upload_binary "./jess-linux-amd64"
upload_binary "./jess-windows-amd64.exe"
