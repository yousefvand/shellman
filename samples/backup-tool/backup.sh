#!/usr/bin/env bash

# Title:         Backup tool
# Description:   Backups defined folders into backups folder in home directory
# Author:        Remisa Yousefvand <remisa.yousefvand@gmail.com>
# Date:          2020-07-07
# Version:       1.0.0

# Exit codes
# ==========
# 0   no error
# 1   script interrupted
# 2   unknown argument
# 3   backup folder doesn't exist

# >>>>>>>>>>>>>>>>>>>>>>>> variables >>>>>>>>>>>>>>>>>>>>>>>>

VERSION="1.0.0"

# assuming all in user home directory
source_paths=(
  'Desktop'
  '.config'
  '.bash_history'
)

backup_success="false"

# <<<<<<<<<<<<<<<<<<<<<<<< variables <<<<<<<<<<<<<<<<<<<<<<<<

# >>>>>>>>>>>>>>>>>>>>>>>> event handling >>>>>>>>>>>>>>>>>>>>>>>>
function cleanup () {
  rm "$backup_dir"/"$backup_file"
}


# CTRL+C event handler
function on_ctrl_c() {
  echo # Set cursor to the next line of '^C'
  tput cnorm # show cursor. You need this if animation is used.

  if [ "$backup_success" = "true" ]; then
    exit 0
  else
    cleanup # Call cleanup function
    echo `tput setaf 1`Backup canceled by user!`tput sgr0`
    exit 1 # Don't remove. Use a number (1-255) for error code.
  fi
}
# Put this line at the beginning of your script entry point (after functions used by event handlers).
# Register CTRL+C event handler
trap on_ctrl_c SIGINT

# <<<<<<<<<<<<<<<<<<<<<<<< event handling <<<<<<<<<<<<<<<<<<<<<<<<

# >>>>>>>>>>>>>>>>>>>>>>>> animation frames >>>>>>>>>>>>>>>>>>>>>>>>

# Your frames need to have the exact same width and height.
# If they are different in size, fill unused space with `space`s (no `TAB`s).

IFS='' read -r -d '' frames[1] <<"EOF"
 .       .
 |\_---_/|
/   o_O   \
|    U    |
\  ._I_.  /
 `-_____-'
EOF

IFS='' read -r -d '' frames[2] <<"EOF"
 .       .
 |\_---_/|
/   O_o   \
|    U    |
\  ._I_.  /
 `-_____-'
EOF

# <<<<<<<<<<<<<<<<<<<<<<<< animation frames <<<<<<<<<<<<<<<<<<<<<<<<

# >>>>>>>>>>>>>>>>>>>>>>>> functions >>>>>>>>>>>>>>>>>>>>>>>>

function show_help () {
  echo "
Usage: $0 [options]

Options:
  -a, --animation   Show animation after
                    backup completed
  -h, --help        Show this help message
  -o, --output      Backup directory
  -v, --version     Show version info
"
}

function play_animation () {
  echo "TODO: play animation"
}

function backup () {
  echo `tput setaf 2`Backup started...`tput sgr0`
  echo # empty line
  current_path=`pwd` # save current directory
  cd "$HOME"
  tar -czvf "$backup_dir"/"$backup_file" "${source_paths[@]}"

  if [[ $? != 0 ]]; then
    echo `tput setaf 1`Unknown error. Backup failed!`tput sgr0`
    cd "$current_path" # restore current directory
    exit 4 # TODO: document error code
  fi

  cd "$current_path" # restore current directory
}

# Usage: animate frames_array interval
function animate () {
  local frames=("$@")
  ((last_index=${#frames[@]} - 1))
  local interval=${frames[last_index]}
  unset frames[last_index]

  # Comment out next two lines if you are using CTRL+C event handler.
	# trap 'tput cnorm; echo' EXIT
  # trap 'exit 127' HUP INT TERM

  tput civis # hide cursor
  tput sc # save cursor position

  while true; do
    for frame in "${frames[@]}"; do
      tput rc # restore cursor position
      echo "$frame"
      sleep "$interval"
    done
  done
}


# <<<<<<<<<<<<<<<<<<<<<<<< functions <<<<<<<<<<<<<<<<<<<<<<<<

# >>>>>>>>>>>>>>>>>>>>>>>> argument parsing >>>>>>>>>>>>>>>>>>>>>>>>

while [[ $# > 0 ]]; do
  case "$1" in
    -h|--help)
    show_help # call "show_help" function
    exit 0 # we don't process any other argument
    ;;
    -v|--version)
    echo "Backup tool v. $VERSION"
    exit 0 # we don't process any other argument
    ;;
    -o|--output)
    backup_dir="$2"
    shift 2 # one for switch and another for its value
    ;;
    -a|--animation)
    play_animation="true"
    shift # shift one for flag itself (we know it is either -a or --animation)
    ;;
    *) # unknown flag/switch
    echo `tput setaf 1`Error! Unknown argument: "$1"`tput sgr0`
    exit 2
    ;;
  esac
done

# <<<<<<<<<<<<<<<<<<<<<<<< argument parsing <<<<<<<<<<<<<<<<<<<<<<<<

# Setting default values
: "${backup_dir:=$HOME/backups}"
: "${play_animation:=false}"

backup_file=`date -I`.tar.gz

if [ ! -d "$backup_dir" ]; then
  echo `tput setaf 1`"Error! Backup directory doesn't exist: $backup_dir"`tput sgr0`
  exit 3
fi

# call backup function. It has access to global variables.
backup

echo `tput setaf 4`Backup complete.`tput sgr0`
backup_success="true"

if [ "$play_animation" = "true" ]; then
  animate "${frames[@]}" 0.5
fi
