#!/bin/bash

. ./build.sh
sudo cp ./jess-darwin-amd64 /usr/local/bin/jess
jess -v
echo -e "\033[32mJess installed successfully!  \033[m"