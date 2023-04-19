#!/bin/bash
set -e

# Install jq if it is not installed already
if ! command -v jq &> /dev/null; then
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
fi

# Set the binary name and repository URL here
BINARY_NAME="jess"
REPOSITORY="https://github.com/assistant-ai/jess/releases"

# Obtain the latest pushed tag
LATEST_TAG=$(curl --silent "https://api.github.com/repos/assistant-ai/jess/releases/latest" | jq -r .tag_name)

# Check the OS and architecture to download the correct binary
OS="$(uname | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
[ "$ARCH" == "x86_64" ] && ARCH="amd64"

# Download the binary
curl -L -o /tmp/${BINARY_NAME} "${REPOSITORY}/download/${LATEST_TAG}/${BINARY_NAME}-${OS}-${ARCH}"

# Make it executable and move it to /usr/local/bin
chmod +x /tmp/${BINARY_NAME}
sudo mv /tmp/${BINARY_NAME} /usr/local/bin/${BINARY_NAME}

echo "${BINARY_NAME} installed successfully!"
