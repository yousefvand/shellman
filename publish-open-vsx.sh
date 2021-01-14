#!/usr/bin/env bash

# Title:         Publish open-vsx
# Description:   Publish vscode extension
# Author:        Remisa <remisa.yousefvand@gmail.com>
# Date:          2021-01-14
# Version:       1.0.0
# Usage:         bash publish.sh <Personal Access Token>

echo "$1" | npx ovsx publish -p $1