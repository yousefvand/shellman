#!/usr/bin/env bash

# Title:         Publish
# Description:   Publish vscode extension
# Author:        Remisa <remisa.yousefvand@gmail.com>
# Date:          2018-11-06
# Version:       1.0.0
# Usage:         bash publish.sh <Personal Access Token>

echo "$1" | vsce login Remisa
vsce publish -p $1
