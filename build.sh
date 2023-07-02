#!/bin/bash

build_app() {
  local GOOS=$1
  local GOARCH=$2
  local filename=$3
  
  CGO_ENABLED=1 GOOS="$GOOS" GOARCH="$GOARCH" go build -ldflags="-X main.version=$VERSION" -o "$filename"
}

VERSION=$(cat VERSION)

architecture=""
case $(uname -m) in
    i386)   architecture="386" ;;
    i686)   architecture="386" ;;
    x86_64)
      if [ $os = "darwin" ]; then
        architecture="arm64"
      else
        architecture="amd64"
      fi
      ;;
    arm64)    architecture="arm64" ;; # handle ARM macs
esac

os=""
case "$OSTYPE" in
  darwin*)  os="darwin" ;;
  linux*)   os="LINUX" ;;
  msys*)    os="windows" ;;
  cygwin*)  os="windows" ;;
  *)        os="unknown: $OSTYPE" ;;
esac

## check if os is windows and exit
if [ "$os" == "windows" ]; then
  echo -e "\033[31mERROR: Windows is not supported by this builder \033[m"
  echo "try to use build.ps1 instead of build.sh"
  exit 1
fi

echo "---------------------------------------------"
echo "building JESSICA:"
echo "version: $VERSION"
echo "OS: $os"
echo "architecture: $architecture"




# executing with handling error
(
  go mod tidy &&
  build_app "$os" $architecture "jess-${os}-${architecture}" &&
  # print with green color
  echo -e "\033[32mSUCCESS building for $os on $architecture \033[m"
) || {
  #  print with red color
  echo -e "\033[31mERROR while building for $os on $architecture \033[m"
  echo ""
  exit 1
  }

