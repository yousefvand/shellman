#!/usr/bin/env bash

# Usage: pac_man inputString interval pad
# Example: pacman "Hello World" 0.5 "*"
function pac_man () {
  local string="${1}"
  local interval="${2}"
  : "${interval:=0.2}"
  local pad="${3}"
  : "${pad:=.}"
  local length=${#string}
  local padding=""

  # Comment out next two lines if you are using CTRL+C event handler.
  trap 'tput cnorm; echo' EXIT
  trap 'exit 127' HUP INT TERM

  tput civis # hide cursor
  tput sc # save cursor position

  for((i=0;i<=length;i++)); do
    tput rc
    echo "${padding}c${string:i:length}"
    sleep "$interval"
    tput rc
    echo "${padding}C${string:i:length}"
    sleep "${interval}"
    padding+="${pad}"
  done

  tput cnorm
  tput rc
  echo "${padding}"
}

# Usage: pac_man inputString interval pad
pac_man "Hello World" 0.2 "."
