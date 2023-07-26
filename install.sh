#!/bin/bash
set -e


# attach external functions
source ./bash_functions.sh
# Install jq if it is not installed already
install_jq

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

jess -v
echo -e "\033[32m${BINARY_NAME} installed successfully!  \033[m"