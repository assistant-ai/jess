#!/bin/bash
source ./bash_functions.sh

. ./build.sh

# Call the function with the search file name prefix and store the result in a variable
search_name_prefix="jess"  # Replace this with the desired prefix of the file name
found_file=$(find_file_by_first_part_of_name "$search_name_prefix")

sudo cp $found_file /usr/local/bin/jess
jess -v
echo -e "\033[32mJess installed successfully!  \033[m"