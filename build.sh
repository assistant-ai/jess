#!/bin/bash

# attach external functions
source ./bash_functions.sh
# Read version from the VERSION file
VERSION=$(cat VERSION)

# Main execution
ARCHITECTURE=$(get_architecture)
OS=$(get_os)
check_windows_os $OS

build_jessica $VERSION $OS $ARCHITECTURE




