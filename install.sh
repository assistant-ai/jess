#!/bin/bash
set -e

# Set the binary name and repository URL here
BINARY_NAME="jess"
MAIN_REPOSITORY="https://github.com/assistant-ai/jess"
RELEASE_REPOSITORY="$MAIN_REPOSITORY/releases"

# Collecting extra fucntions that ma be useful for installation
curl -L -o /tmp/bash_functions.sh "$MAIN_REPOSITORY/raw/master/bash_functions.sh"
echo "bash_functions.sh downloaded successfully!"
# attach external functions
source /tmp/bash_functions.sh
echo "bash_functions.sh sourced successfully!"

# delete temporary functions file
rm /tmp/bash_functions.sh
# Install jq if it is not installed already
install_jq



# Obtain the latest pushed tag
LATEST_TAG=$(curl --silent "https://api.github.com/repos/assistant-ai/jess/releases/latest" | jq -r .tag_name)

# Check the OS and architecture to download the correct binary
OS="$(uname | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
[ "$ARCH" == "x86_64" ] && ARCH="amd64"

# Download the binary
curl -L -o /tmp/${BINARY_NAME} "${RELEASE_REPOSITORY}/download/${LATEST_TAG}/${BINARY_NAME}-${OS}-${ARCH}"

# Make it executable and move it to /usr/local/bin
chmod +x /tmp/${BINARY_NAME}
echo "it is required to enter your password to move the binary to /usr/local/bin"
sudo mv /tmp/${BINARY_NAME} /usr/local/bin/${BINARY_NAME}

jess -v
echo -e "\033[32m${BINARY_NAME} installed successfully!  \033[m"
echo -e "\033[32m${BINARY_NAME} was installed into /usr/local/bin/${BINARY_NAME} \033[m"
echo -e "\033[32mPlease run \033[33mjess -h \033[32mto see the available commands.  \033[m"