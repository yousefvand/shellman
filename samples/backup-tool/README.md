# Backup tool

In this example we are gonna write an script to backup some folders as a `tar.gz` archive to specified directory (default directory `~/backups`). Handle `CTRL+C` and cleanup if script interrupted. Optionally play animation after backup.

First we create our script structure.

## Script structure

1. shebang
2. summary & error codes
3. event handlers (related functions before them)
4. animation frames
5. functions
6. argument parsing
7. main code (entry point)

Create an empty file named `backup.sh`. Continue with following steps:

### shebang

Type `shebang` and press `ENTER`. From dropdown menu press `ENTER` on default option (`bash`):

```bash
#!/usr/bin/env bash

```

### summary

Type `summary` and press `ENTER`. Use `TAB` and fill out fields.

```bash
#!/usr/bin/env bash

# Title:         Backup tool
# Description:   Backups defined folders into backups folder in home directory
# Author:        Remisa Yousefvand <remisa.yousefvand@gmail.com>
# Date:          2020-07-07
# Version:       1.0.0

# Exit codes
# ==========
# 0   no error
# TODO: more later

```

### event handlers

We want to capture `CTRL+C` if script is interrupted by user so we can run our cleanup code.

Use `region` snippet to create `event handling` region. Inside that region:

Type `event` and with arrow up/down keys select `event CTRL+C` and hit `ENTER`. We are using code snippets from Shellman. We already have used `shebang`, `summary` snippets.

Call `cleanup` function in `on_ctrl_c` function and use `color` snippet (type `color`) and select `red` to print an error message. Above `on_ctrl_c` function use `func` snippet to define a function. Type `cleanup` as function name and add a `TODO` comment as a place holder. define `1` as error code for interruption in `Exit codes` section.

```bash
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
# TODO: more later

# >>>>>>>>>>>>>>>>>>>>>>>> event handling >>>>>>>>>>>>>>>>>>>>>>>>
function cleanup () {
  echo "TODO: cleanup code"
}


# CTRL+C event handler
function on_ctrl_c() {
  echo # Set cursor to the next line of '^C'
  tput cnorm # show cursor. You need this if animation is used.
  cleanup # Call cleanup function
  echo `tput setaf 1`Backup canceled by user!`tput sgr0`

  exit 1 # Don't remove. Use a number (1-255) for error code.
}
# Put this line at the beginning of your script entry point (after functions used by event handlers).
# Register CTRL+C event handler
trap on_ctrl_c SIGINT

# <<<<<<<<<<<<<<<<<<<<<<<< event handling <<<<<<<<<<<<<<<<<<<<<<<<

```

### animation frames

Create a `animation frames` region using `region` snippet. Add two animation frames using `animation frame` snippet. Remember to update frame number when adding new frames. Here we use [funny cat](../animation/funny-cat.sh).

```bash
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
# TODO: more later

# >>>>>>>>>>>>>>>>>>>>>>>> event handling >>>>>>>>>>>>>>>>>>>>>>>>
function cleanup () {
  echo "TODO: cleanup code"
}


# CTRL+C event handler
function on_ctrl_c() {
  echo # Set cursor to the next line of '^C'
  tput cnorm # show cursor. You need this if animation is used.
  cleanup # Call cleanup function
  echo `tput setaf 1`Backup canceled by user!`tput sgr0`

  exit 1 # Don't remove. Use a number (1-255) for error code.
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

```

### functions

Add a `functions` region. Put a `TODO` comments. Later we write functions here which are gonna do the real job.

```bash
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
# TODO: more later

# >>>>>>>>>>>>>>>>>>>>>>>> event handling >>>>>>>>>>>>>>>>>>>>>>>>
function cleanup () {
  echo "TODO: cleanup code"
}


# CTRL+C event handler
function on_ctrl_c() {
  echo # Set cursor to the next line of '^C'
  tput cnorm # show cursor. You need this if animation is used.
  cleanup # Call cleanup function
  echo `tput setaf 1`Backup canceled by user!`tput sgr0`

  exit 1 # Don't remove. Use a number (1-255) for error code.
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

# TODO: functions

# <<<<<<<<<<<<<<<<<<<<<<<< functions <<<<<<<<<<<<<<<<<<<<<<<<

```

### argument parsing

Add a `argument parsing` region and use `argument parsing` snippet as a scaffold for your parsing process. We need to design our script interface now. As a good practice we add `help` and `version` flags to all scripts. Also we accept `output` folder as a switch (a switch accepts parameter(s) while a flag does not) and if it was not passed to our script we are gonna use `~/backups` as default. To add a flag as an example lets play `animation` if such a flag is passed to our script (default to not play animation). Here is our design summary:

| short | long        | Description                           | Type   |
| ----- | ----------- | ------------------------------------- | ------ |
| -h    | --help      | Prints help message                   | Flag   |
| -v    | --version   | Prints version info                   | Flag   |
| -o    | --output    | Backup directory                      | Switch |
| -a    | --animation | Show animation after backup completed | Flag   |

We define functions when needed and add a place holder for that function in `functions` region so we can write those functions later. It turns out we need a variable to store script version. Define it as top of the script as `VERSION` (in summary section). To write error messages use `color` snippet with `red` ink.

```bash
#!/usr/bin/env bash

# Title:         Backup tool
# Description:   Backups defined folders into backups folder in home directory
# Author:        Remisa Yousefvand <remisa.yousefvand@gmail.com>
# Date:          2020-07-07
# Version:       1.0.0
VERSION="1.0.0"

# Exit codes
# ==========
# 0   no error
# 1   script interrupted
# 2   unknown argument
# 3   backup folder doesn't exist
# 4   backup failed due to tar error

# >>>>>>>>>>>>>>>>>>>>>>>> event handling >>>>>>>>>>>>>>>>>>>>>>>>
function cleanup () {
  echo "TODO: cleanup code"
}


# CTRL+C event handler
function on_ctrl_c() {
  echo # Set cursor to the next line of '^C'
  tput cnorm # show cursor. You need this if animation is used.
  cleanup # Call cleanup function
  echo `tput setaf 1`Backup canceled by user!`tput sgr0`

  exit 1 # Don't remove. Use a number (1-255) for error code.
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

function help () {
  echo "
Usage: ${0} [options]

Options:
  -a, --animation   Show animation after
                    backup completed
  -h, --help        Show this help message
  -o, --output      Backup directory
  -v, --version     Show version info
"
}

function playAnimation () {
  echo "TODO: play animation"
}

function backup () {
  echo "TODO: backup"
}

# <<<<<<<<<<<<<<<<<<<<<<<< functions <<<<<<<<<<<<<<<<<<<<<<<<

# >>>>>>>>>>>>>>>>>>>>>>>> argument parsing >>>>>>>>>>>>>>>>>>>>>>>>

while [[ $# > 0 ]]; do
  case "${1}" in
    -h|--help)
    help # call "help" function
    exit 0 # we don't process any other argument
    ;;
    -v|--version)
    echo "Backup tool v. ${VERSION}"
    exit 0 # we don't process any other argument
    ;;
    -o|--output)
    backupDirectory="${2}"
    shift 2 # one for switch and another for its value
    echo "${backupDirectory}" # test. TODO: remove
    ;;
    -a|--animation)
    playAnimation="true"
    shift # shift one for flag itself (we know it is either -a or --animation)
    echo "${playAnimation}" # test. TODO: remove
    ;;
    *) # unknown flag/switch
    echo `tput setaf 1`Error! Unknown argument: "${1}"`tput sgr0`
    exit 2
    ;;
  esac
done

# <<<<<<<<<<<<<<<<<<<<<<<< argument parsing <<<<<<<<<<<<<<<<<<<<<<<<

```

Test your script's argument parsing then remove test echo commands.

### main code

Now we have parsed arguments and need to set default values if any variable is not set. To do so use `assign if empty` snippet to set default values for `backupDirectory` and `playAnimation` variables.

```bash
#!/usr/bin/env bash

# Title:         Backup tool
# Description:   Backups defined folders into backups folder in home directory
# Author:        Remisa Yousefvand <remisa.yousefvand@gmail.com>
# Date:          2020-07-07
# Version:       1.0.0
VERSION="1.0.0"

# Exit codes
# ==========
# 0   no error
# 1   script interrupted
# 2   unknown argument
# 3   backup folder doesn't exist
# TODO: more later

# >>>>>>>>>>>>>>>>>>>>>>>> event handling >>>>>>>>>>>>>>>>>>>>>>>>
function cleanup () {
  echo "cleanup code"
}


# CTRL+C event handler
function on_ctrl_c() {
  echo # Set cursor to the next line of '^C'
  tput cnorm # show cursor. You need this if animation is used.
  cleanup # Call cleanup function
  echo `tput setaf 1`Backup canceled by user!`tput sgr0`

  exit 1 # Don't remove. Use a number (1-255) for error code.
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

function help () {
  echo "
Usage: ${0} [options]

Options:
  -a, --animation   Show animation after
                    backup completed
  -h, --help        Show this help message
  -o, --output      Backup directory
  -v, --version     Show version info
"
}

function playAnimation () {
  echo "TODO: play animation"
}

function backup () {
  echo "TODO: backup"
}

# <<<<<<<<<<<<<<<<<<<<<<<< functions <<<<<<<<<<<<<<<<<<<<<<<<

# >>>>>>>>>>>>>>>>>>>>>>>> argument parsing >>>>>>>>>>>>>>>>>>>>>>>>

while [[ $# > 0 ]]; do
  case "${1}" in
    -h|--help)
    help # call "help" function
    exit 0 # we don't process any other argument
    ;;
    -v|--version)
    echo "Backup tool v. ${VERSION}"
    exit 0 # we don't process any other argument
    ;;
    -o|--output)
    backupDirectory="${2}"
    shift 2 # one for switch and another for its value
    ;;
    -a|--animation)
    playAnimation="true"
    shift # shift one for flag itself (we know it is either -a or --animation)
    ;;
    *) # unknown flag/switch
    echo `tput setaf 1`Error! Unknown argument: "${1}"`tput sgr0`
    exit 2
    ;;
  esac
done

# <<<<<<<<<<<<<<<<<<<<<<<< argument parsing <<<<<<<<<<<<<<<<<<<<<<<<

# Setting default values
: "${backupDirectory:=~/backups}"
: "${playAnimation:=false}"

# Test default values. TODO: remove
echo "${backupDirectory}"
echo "${playAnimation}"

```

Test script and remove test echos for default values.

First we check if backup folder exists and if not we exit with error code 3 as documented before. Use `if directory exists` and a negation (`!`) to run the code when backup directory doesn't exist:

```bash
if [ ! -d "${backupDirectory}" ]; then
  echo `tput setaf 1`"Error! Backup directory doesn't exist: ${backupDirectory}"`tput sgr0`
  exit 3
fi
```

Now we need to write the actual backup code inside `backup` function and call it but we should take care of few things before that. Since desired folders to backup are known and are not gonna change frequently we use an array at the top of script to hard-code those paths. Use `array declare` snippet to define this array. Also we need to chose a name for our backup file and store it inside a variable. I'm gonna use `yyyy-mm-dd.tar.gz` using `date now short` snippet.

```bash
#!/usr/bin/env bash

# Title:         Backup tool
# Description:   Backups defined folders into backups folder in home directory
# Author:        Remisa Yousefvand <remisa.yousefvand@gmail.com>
# Date:          2020-07-07
# Version:       1.0.0
VERSION="1.0.0"

# Exit codes
# ==========
# 0   no error
# 1   script interrupted
# 2   unknown argument
# 3   backup folder doesn't exist
# TODO: more later

sourcePaths=(
  '~/Desktop'
  '~/.config'
  '.bash_history'
)


# >>>>>>>>>>>>>>>>>>>>>>>> event handling >>>>>>>>>>>>>>>>>>>>>>>>
function cleanup () {
  echo "cleanup code"
}


# CTRL+C event handler
function on_ctrl_c() {
  echo # Set cursor to the next line of '^C'
  tput cnorm # show cursor. You need this if animation is used.
  cleanup # Call cleanup function
  echo `tput setaf 1`Backup canceled by user!`tput sgr0`

  exit 1 # Don't remove. Use a number (1-255) for error code.
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

function help () {
  echo "
Usage: ${0} [options]

Options:
  -a, --animation   Show animation after
                    backup completed
  -h, --help        Show this help message
  -o, --output      Backup directory
  -v, --version     Show version info
"
}

function playAnimation () {
  echo "TODO: play animation"
}

function backup () {
  echo "TODO: backup"
}

# <<<<<<<<<<<<<<<<<<<<<<<< functions <<<<<<<<<<<<<<<<<<<<<<<<

# >>>>>>>>>>>>>>>>>>>>>>>> argument parsing >>>>>>>>>>>>>>>>>>>>>>>>

while [[ $# > 0 ]]; do
  case "${1}" in
    -h|--help)
    help # call "help" function
    exit 0 # we don't process any other argument
    ;;
    -v|--version)
    echo "Backup tool v. ${VERSION}"
    exit 0 # we don't process any other argument
    ;;
    -o|--output)
    backupDirectory="${2}"
    shift 2 # one for switch and another for its value
    ;;
    -a|--animation)
    playAnimation="true"
    shift # shift one for flag itself (we know it is either -a or --animation)
    ;;
    *) # unknown flag/switch
    echo `tput setaf 1`Error! Unknown argument: "${1}"`tput sgr0`
    exit 2
    ;;
  esac
done

# <<<<<<<<<<<<<<<<<<<<<<<< argument parsing <<<<<<<<<<<<<<<<<<<<<<<<

# Setting default values
: "${backupDirectory:=~/backups}"
: "${playAnimation:=false}"

# test array. TODO: remove
echo "${sourcePaths[@]}"

backupFile=`date -I`.tar.gz

# test backup file name. TODO: remove
echo "${backupFile}"

if [ ! -d "${backupDirectory}" ]; then
  echo `tput setaf 1`"Error! Backup directory doesn't exist: ${backupDirectory}"`tput sgr0`
  exit 3
fi

```

Test the script and remove test echos.

```bash
#!/usr/bin/env bash

# Title:         Backup tool
# Description:   Backups defined folders into backups folder in home directory
# Author:        Remisa Yousefvand <remisa.yousefvand@gmail.com>
# Date:          2020-07-07
# Version:       1.0.0
VERSION="1.0.0"

# Exit codes
# ==========
# 0   no error
# 1   script interrupted
# 2   unknown argument
# 3   backup folder doesn't exist
# TODO: more later

# assuming all in user home directory
sourcePaths=(
  'Desktop'
  '.config'
  '.bash_history'
)


# >>>>>>>>>>>>>>>>>>>>>>>> event handling >>>>>>>>>>>>>>>>>>>>>>>>
function cleanup () {
  echo "cleanup code"
}


# CTRL+C event handler
function on_ctrl_c() {
  echo # Set cursor to the next line of '^C'
  tput cnorm # show cursor. You need this if animation is used.
  cleanup # Call cleanup function
  echo `tput setaf 1`Backup canceled by user!`tput sgr0`

  exit 1 # Don't remove. Use a number (1-255) for error code.
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

function help () {
  echo "
Usage: ${0} [options]

Options:
  -a, --animation   Show animation after
                    backup completed
  -h, --help        Show this help message
  -o, --output      Backup directory
  -v, --version     Show version info
"
}

function playAnimation () {
  echo "TODO: play animation"
}

function backup () {
  echo `tput setaf 2`Backup started...`tput sgr0`
  echo # empty line
  currentPath=`pwd` # save current directory
  cd "${HOME}"
  tar -czvf "${backupDirectory}/${backupFile}" "${sourcePaths[@]}"
  cd "${currentPath}" # restore current directory
}

# <<<<<<<<<<<<<<<<<<<<<<<< functions <<<<<<<<<<<<<<<<<<<<<<<<

# >>>>>>>>>>>>>>>>>>>>>>>> argument parsing >>>>>>>>>>>>>>>>>>>>>>>>

while [[ $# > 0 ]]; do
  case "${1}" in
    -h|--help)
    help # call "help" function
    exit 0 # we don't process any other argument
    ;;
    -v|--version)
    echo "Backup tool v. ${VERSION}"
    exit 0 # we don't process any other argument
    ;;
    -o|--output)
    backupDirectory="${2}"
    shift 2 # one for switch and another for its value
    ;;
    -a|--animation)
    playAnimation="true"
    shift # shift one for flag itself (we know it is either -a or --animation)
    ;;
    *) # unknown flag/switch
    echo `tput setaf 1`Error! Unknown argument: "${1}"`tput sgr0`
    exit 2
    ;;
  esac
done

# <<<<<<<<<<<<<<<<<<<<<<<< argument parsing <<<<<<<<<<<<<<<<<<<<<<<<

# Setting default values
: "${backupDirectory:=$HOME/backups}"
: "${playAnimation:=false}"

# test array. TODO: remove
echo "${sourcePaths[@]}"

backupFile=`date -I`.tar.gz

# test backup file name. TODO: remove
echo "${backupFile}"

if [ ! -d "${backupDirectory}" ]; then
  echo `tput setaf 1`"Error! Backup directory doesn't exist: ${backupDirectory}"`tput sgr0`
  exit 3
fi

# call backup function. It has access to global variables.
backup

```

If user interrupts script we need to delete incomplete `tar.gz` file. In this example `cleanup` function doesn't worth to be written as a separate function (delete the file inside `on_ctrl_c` function) but this is a toy example and usually cleanup code is more than just a single line of code.

```bash
#!/usr/bin/env bash

# Title:         Backup tool
# Description:   Backups defined folders into backups folder in home directory
# Author:        Remisa Yousefvand <remisa.yousefvand@gmail.com>
# Date:          2020-07-07
# Version:       1.0.0
VERSION="1.0.0"

# Exit codes
# ==========
# 0   no error
# 1   script interrupted
# 2   unknown argument
# 3   backup folder doesn't exist
# TODO: more later

# assuming all in user home directory
sourcePaths=(
  'Desktop'
  '.config'
  '.bash_history'
)


# >>>>>>>>>>>>>>>>>>>>>>>> event handling >>>>>>>>>>>>>>>>>>>>>>>>
function cleanup () {
  rm "${backupDirectory}/${backupFile}"
}


# CTRL+C event handler
function on_ctrl_c() {
  echo # Set cursor to the next line of '^C'
  tput cnorm # show cursor. You need this if animation is used.
  cleanup # Call cleanup function
  echo `tput setaf 1`Backup canceled by user!`tput sgr0`

  exit 1 # Don't remove. Use a number (1-255) for error code.
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

function help () {
  echo "
Usage: ${0} [options]

Options:
  -a, --animation   Show animation after
                    backup completed
  -h, --help        Show this help message
  -o, --output      Backup directory
  -v, --version     Show version info
"
}

function playAnimation () {
  echo "TODO: play animation"
}

function backup () {
  echo `tput setaf 2`Backup started...`tput sgr0`
  echo # empty line
  currentPath=`pwd` # save current directory
  cd "${HOME}"
  tar -czvf "${backupDirectory}/${backupFile}" "${sourcePaths[@]}"
  cd "${currentPath}" # restore current directory
}

# <<<<<<<<<<<<<<<<<<<<<<<< functions <<<<<<<<<<<<<<<<<<<<<<<<

# >>>>>>>>>>>>>>>>>>>>>>>> argument parsing >>>>>>>>>>>>>>>>>>>>>>>>

while [[ $# > 0 ]]; do
  case "${1}" in
    -h|--help)
    help # call "help" function
    exit 0 # we don't process any other argument
    ;;
    -v|--version)
    echo "Backup tool v. ${VERSION}"
    exit 0 # we don't process any other argument
    ;;
    -o|--output)
    backupDirectory="${2}"
    shift 2 # one for switch and another for its value
    ;;
    -a|--animation)
    playAnimation="true"
    shift # shift one for flag itself (we know it is either -a or --animation)
    ;;
    *) # unknown flag/switch
    echo `tput setaf 1`Error! Unknown argument: "${1}"`tput sgr0`
    exit 2
    ;;
  esac
done

# <<<<<<<<<<<<<<<<<<<<<<<< argument parsing <<<<<<<<<<<<<<<<<<<<<<<<

# Setting default values
: "${backupDirectory:=$HOME/backups}"
: "${playAnimation:=false}"

backupFile=`date -I`.tar.gz

if [ ! -d "${backupDirectory}" ]; then
  echo `tput setaf 1`"Error! Backup directory doesn't exist: ${backupDirectory}"`tput sgr0`
  exit 3
fi

# call backup function. It has access to global variables.
backup

```

Inside `backup` function we need to make sure `tar` has completed successfully (if `tar` returned `0`). Use `cmd failure check` snippet for that. It turns out we need a new error code for when `tar` fails. Document error code `4` as when `tar` fails.

```bash
#!/usr/bin/env bash

# Title:         Backup tool
# Description:   Backups defined folders into backups folder in home directory
# Author:        Remisa Yousefvand <remisa.yousefvand@gmail.com>
# Date:          2020-07-07
# Version:       1.0.0
VERSION="1.0.0"

# Exit codes
# ==========
# 0   no error
# 1   script interrupted
# 2   unknown argument
# 3   backup folder doesn't exist
# 4   backup failed due to tar error

# assuming all in user home directory
sourcePaths=(
  'Desktop'
  '.config'
  '.bash_history'
)

# >>>>>>>>>>>>>>>>>>>>>>>> event handling >>>>>>>>>>>>>>>>>>>>>>>>
function cleanup () {
  rm "${backupDirectory}/${backupFile}"
}

# CTRL+C event handler
function on_ctrl_c() {
  echo # Set cursor to the next line of '^C'
  tput cnorm # show cursor. You need this if animation is used.
  cleanup # Call cleanup function
  echo `tput setaf 1`Backup canceled by user!`tput sgr0`

  exit 1 # Don't remove. Use a number (1-255) for error code.
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

function help () {
  echo "
Usage: ${0} [options]

Options:
  -a, --animation   Show animation after
                    backup completed
  -h, --help        Show this help message
  -o, --output      Backup directory
  -v, --version     Show version info
"
}

function playAnimation () {
  echo "TODO: play animation"
}

function backup () {
  echo `tput setaf 2`Backup started...`tput sgr0`
  echo # empty line
  currentPath=`pwd` # save current directory
  cd "${HOME}"
  tar -czvf "${backupDirectory}/${backupFile}" "${sourcePaths[@]}"

  if [[ $? != 0 ]]; then
    echo `tput setaf 1`Unknown error. Backup failed!`tput sgr0`
    cd "${currentPath}" # restore current directory
    exit 4
  fi

  cd "${currentPath}" # restore current directory
}

# <<<<<<<<<<<<<<<<<<<<<<<< functions <<<<<<<<<<<<<<<<<<<<<<<<

# >>>>>>>>>>>>>>>>>>>>>>>> argument parsing >>>>>>>>>>>>>>>>>>>>>>>>

while [[ $# > 0 ]]; do
  case "${1}" in
    -h|--help)
    help # call "help" function
    exit 0 # we don't process any other argument
    ;;
    -v|--version)
    echo "Backup tool v. ${VERSION}"
    exit 0 # we don't process any other argument
    ;;
    -o|--output)
    backupDirectory="${2}"
    shift 2 # one for switch and another for its value
    ;;
    -a|--animation)
    playAnimation="true"
    shift # shift one for flag itself (we know it is either -a or --animation)
    ;;
    *) # unknown flag/switch
    echo `tput setaf 1`Error! Unknown argument: "${1}"`tput sgr0`
    exit 2
    ;;
  esac
done

# <<<<<<<<<<<<<<<<<<<<<<<< argument parsing <<<<<<<<<<<<<<<<<<<<<<<<

# Setting default values
: "${backupDirectory:=$HOME/backups}"
: "${playAnimation:=false}"

backupFile=`date -I`.tar.gz

if [ ! -d "${backupDirectory}" ]; then
  echo `tput setaf 1`"Error! Backup directory doesn't exist: ${backupDirectory}"`tput sgr0`
  exit 3
fi

# call backup function. It has access to global variables.
backup

```

Finally play the animation if flag `-a` or `--animation` is passed to script. Use `if string =` snippet and check if `playAnimation` variable is equal to `true`.

For playing animation we need the `animate` function. Inside `functions` region use `fn animation animate` snippet to add this function. As this function documentation says we need to comment out `CTRL+C` event handler because we already have `tput cnorm` in `on_ctrl_c` function which makes cursor visible. Now we can call `animate` function at the end of our script. Use `fx animation animate` snippet to call `animate` function and pass it arguments.

Here we face a new problem. To stop animation user needs to press `CTRL+C` but currently we consider this as an incomplete backup. We need a variable to ignore this if backup is already successfully completed (define `backupSuccess` variable). Here is the [final script](backup.sh).

```bash
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
# 4   backup failed due to tar error

# >>>>>>>>>>>>>>>>>>>>>>>> variables >>>>>>>>>>>>>>>>>>>>>>>>

VERSION="1.0.0"

# assuming all in user home directory
sourcePaths=(
  'Desktop'
  '.config'
  '.bash_history'
)

backupSuccess="false"

# <<<<<<<<<<<<<<<<<<<<<<<< variables <<<<<<<<<<<<<<<<<<<<<<<<

# >>>>>>>>>>>>>>>>>>>>>>>> event handling >>>>>>>>>>>>>>>>>>>>>>>>
function cleanup () {
  rm "${backupDirectory}/${backupFile}"
}


# CTRL+C event handler
function on_ctrl_c() {
  echo # Set cursor to the next line of '^C'
  tput cnorm # show cursor. You need this if animation is used.

  if [ "${backupSuccess}" = "true" ]; then
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

function help () {
  echo "
Usage: ${0} [options]

Options:
  -a, --animation   Show animation after
                    backup completed
  -h, --help        Show this help message
  -o, --output      Backup directory
  -v, --version     Show version info
"
}

function playAnimation () {
  echo "TODO: play animation"
}

function backup () {
  echo `tput setaf 2`Backup started...`tput sgr0`
  echo # empty line
  currentPath=`pwd` # save current directory
  cd "${HOME}"
  tar -czvf "${backupDirectory}/${backupFile}" "${sourcePaths[@]}"

  if [[ $? != 0 ]]; then
    echo `tput setaf 1`Unknown error. Backup failed!`tput sgr0`
    cd "${currentPath}" # restore current directory
    exit 4
  fi

  cd "${currentPath}" # restore current directory
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
      echo "${frame}"
      sleep "${interval}"
    done
  done
}

# <<<<<<<<<<<<<<<<<<<<<<<< functions <<<<<<<<<<<<<<<<<<<<<<<<

# >>>>>>>>>>>>>>>>>>>>>>>> argument parsing >>>>>>>>>>>>>>>>>>>>>>>>

while [[ $# > 0 ]]; do
  case "${1}" in
    -h|--help)
    help # call "help" function
    exit 0 # we don't process any other argument
    ;;
    -v|--version)
    echo "Backup tool v. ${VERSION}"
    exit 0 # we don't process any other argument
    ;;
    -o|--output)
    backupDirectory="${2}"
    shift 2 # one for switch and another for its value
    ;;
    -a|--animation)
    playAnimation="true"
    shift # shift one for flag itself (we know it is either -a or --animation)
    ;;
    *) # unknown flag/switch
    echo `tput setaf 1`Error! Unknown argument: "${1}"`tput sgr0`
    exit 2
    ;;
  esac
done

# <<<<<<<<<<<<<<<<<<<<<<<< argument parsing <<<<<<<<<<<<<<<<<<<<<<<<

# Setting default values
: "${backupDirectory:=$HOME/backups}"
: "${playAnimation:=false}"

backupFile=`date -I`.tar.gz

if [ ! -d "${backupDirectory}" ]; then
  echo `tput setaf 1`"Error! Backup directory doesn't exist: ${backupDirectory}"`tput sgr0`
  exit 3
fi

# call backup function. It has access to global variables.
backup

echo `tput setaf 4`Backup complete.`tput sgr0`
backupSuccess="true"

if [ "${playAnimation}" = "true" ]; then
  animate "${frames[@]}" 0.5
fi

```

Congratulations! You wrote a complete shell script to automate some task. This is an investment of your time. You consumed more time now to learn shell scripting so you can automate repetitive tasks and consume less time in the future.
