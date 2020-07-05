#!/usr/bin/env bash

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

# Usage: animate frames_array interval
function animate () {
  local frames=("$@")
  ((last_index=${#frames[@]} - 1))
  local interval=${frames[last_index]}
  unset frames[last_index]

  # Comment out next two lines if you are using CTRL+C event handler.
	trap 'tput cnorm; echo' EXIT
  trap 'exit 127' HUP INT TERM

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

# Usage: animate frames_array interval
animate "${frames[@]}" 0.5
