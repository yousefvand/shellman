#!/usr/bin/env bash

# Title:         Publish
# Description:   Publish vscode extension
# Author:        Remisa <remisa.yousefvand@gmail.com>
# Date:          2021-03-14
# Version:       1.1.0
# Usage:         bash publish.sh <Personal Access Token>

echo "${1}" | npx vsce login Remisa
npx vsce publish -p $1
