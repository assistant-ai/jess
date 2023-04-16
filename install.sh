#!/bin/bash
set -e

# Set the binary name, version, and repository URL here
BINARY_NAME="jess"
VERSION="2"
REPOSITORY="https://github.com/assistant-ai/jess/releases/download/${VERSION}"

# Check the OS and architecture to download the correct binary
OS="$(uname | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
[ "$ARCH" == "x86_64" ] && ARCH="amd64"

# Download the binary
curl -L -o /tmp/${BINARY_NAME} "${REPOSITORY}/${BINARY_NAME}-${OS}-${ARCH}"

# Make it executable and move it to /usr/local/bin
chmod +x /tmp/${BINARY_NAME}
sudo mv /tmp/${BINARY_NAME} /usr/local/bin/${BINARY_NAME}

echo "${BINARY_NAME} installed successfully!"