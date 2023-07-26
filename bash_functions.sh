#!/bin/bash


# Function to build the app
build_app() {
  local GOOS=$1
  local GOARCH=$2
  local filename=$3

  CGO_ENABLED=1 GOOS="$GOOS" GOARCH="$GOARCH" go build -ldflags="-X main.version=$VERSION" -o "$filename"
}
# Determine the architecture based on the machine's uname
get_architecture() {
  case $(uname -m) in
    i386|i686) architecture="386" ;;
    x86_64)
      if [ "$os" = "darwin" ]; then
        architecture="arm64"
      else
        architecture="amd64"
      fi
      ;;
    arm64) architecture="arm64" ;; # handle ARM macs
  esac
  echo $architecture
}

# Determine the operating system based on the machine's OSTYPE
get_os() {
  case "$OSTYPE" in
    darwin*) os="darwin" ;;
    linux*)  os="LINUX" ;;
    msys*|cygwin*) os="windows" ;;
    *) os="unknown: $OSTYPE" ;;
  esac
  echo $os
}
# Check if OS is windows and exit
check_windows_os() {
  if [ "$1" == "windows" ]; then
    echo -e "\033[31mERROR: Windows is not supported by this builder \033[m"
    echo "try to use build.ps1 instead of build.sh"
    exit 1
  fi
}
# Build the app
build_jessica() {
  VERSION=$1
  os=$2
  architecture=$3
  echo "---------------------------------------------"
  echo "building JESSICA:"
  echo "version: $VERSION"
  echo "OS: $os"
  echo "architecture: $architecture"

  # Execute with error handling
  (
    go mod tidy &&
    build_app "$os" "$architecture" "jess-${os}-${architecture}" &&
    # Print success message in green color
    echo -e "\033[32mSUCCESS building for $os on $architecture \033[m"
  ) || {
    # Print error message in red color
    echo -e "\033[31mERROR while building for $os on $architecture \033[m"
    echo ""
    exit 1
  }
}

upload_and_remove_binary() {
  local path=$1
  gh release upload "v$VERSION" "${path}"
  rm -rf "${path}"
}
# Function to install jq if it is not already installed
install_jq() {
  echo "jq not found, installing..."
  case "$(uname)" in
    "Linux")
      sudo apt-get update
      sudo apt-get install -y jq
      ;;
    "Darwin")
      brew install jq
      ;;
    *)
      echo "Unsupported OS. Please manually install jq"
      exit 1
      ;;
  esac
  echo "jq installed successfully!"
}

check_or_install_gh() {
  echo "gh not found, installing..."
  case "$(uname)" in
    "Linux")
      sudo apt-get update
      sudo apt-get install -y gh
      ;;
    "Darwin")
      brew install gh
      ;;
    *)
      echo "Unsupported OS. Please manually install gh - github command line tool"
      exit 1
      ;;
  esac
  echo "gh installed successfully!"
}


# Function to download and install the binary
install_binary() {
  # Set the binary name and repository URL here
  local BINARY_NAME="jess"
  local REPOSITORY="https://github.com/assistant-ai/jess/releases"

  # Obtain the latest pushed tag
  local LATEST_TAG=$(curl --silent "https://api.github.com/repos/assistant-ai/jess/releases/latest" | jq -r .tag_name)

  # Check the OS and architecture to download the correct binary
  local OS="$(uname | tr '[:upper:]' '[:lower:]')"
  local ARCH="$(uname -m)"
  [ "$ARCH" == "x86_64" ] && ARCH="amd64"

  # Download the binary
  curl -L -o /tmp/${BINARY_NAME} "${REPOSITORY}/download/${LATEST_TAG}/${BINARY_NAME}-${OS}-${ARCH}"

  # Make it executable and move it to /usr/local/bin
  chmod +x /tmp/${BINARY_NAME}
  sudo mv /tmp/${BINARY_NAME} /usr/local/bin/${BINARY_NAME}

  echo "${BINARY_NAME} installed successfully!"
}


# Function to find the file by the first part of its name in the current directory is going to be use in the future build and install
find_file_by_first_part_of_name() {
    local search_name="$1"

    # Use 'find' command to search for files with the specified name prefix in the current directory
    found_file=$(find . -maxdepth 1 -type f -name "$search_name*" -print -quit)

    if [ -n "$found_file" ]; then
        echo "$(realpath "$found_file")"
        return 0
    else
        return 1
    fi
}