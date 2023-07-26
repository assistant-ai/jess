#!/bin/bash

# Function to display script usage
function display_usage() {
    echo "Usage: $0 [-f] <file>"
    echo "  -f: Flag to force deletion jess with config files"
    exit 1
}

# Default values
force=false

# Process command-line options
while getopts ":f" opt; do
    case "$opt" in
        f)
            force=true
            ;;
        \?)
            echo "Invalid option: -$OPTARG"
            display_usage
            ;;
    esac
done


# Function to check if a file exists
function check_file_exists() {
    if [ ! -e "$1" ]; then
        echo "Error: $1 does not exist."
        exit 1
    fi
}

# Function to prompt for confirmation
function deletion_confirm() {
    echo -e "\033[31mConfirm deleting jess with typing: 'delete jess'  \033[m"
    read -p "" confirmation

    if [ "$confirmation" != "delete jess" ]; then
        echo "Deletion canceled."
        exit 0
    fi
}

check_file_exists /usr/local/bin/jess


deletion_confirm


if [ "$force" = true ]; then
    echo -e "\033[31mDeleting jess with config files... \033[m"
    echo "  ~/.jess/*"
    rm -rf ~/.jess
else
    echo -e "\033[31mDeleting jess...  \033[m"
fi
sudo rm /usr/local/bin/jess
rm ~/.llmchat-client/messages.db


echo "This script will delete the following files:"
echo "  /usr/local/bin/jess"
echo "  ~/.llmchat-client/messages.db"
echo -e "\033[31mJess deleted successfully!  \033[m"