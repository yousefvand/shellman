#!/usr/bin/env bash

# >>>>>>>>>>>>>>>>>>>>>>>> frames >>>>>>>>>>>>>>>>>>>>>>>>

IFS='' read -r -d '' frames[1] <<"EOF"
==~=========
| Shellman |
==~=========
EOF

IFS='' read -r -d '' frames[2] <<"EOF"
===~========
| sHellman |
===~========
EOF

IFS='' read -r -d '' frames[3] <<"EOF"
====~=======
| ShEllman |
====~=======
EOF

IFS='' read -r -d '' frames[4] <<"EOF"
=====~======
| SheLlman |
=====~======
EOF

IFS='' read -r -d '' frames[5] <<"EOF"
======~=====
| ShelLman |
======~=====
EOF

IFS='' read -r -d '' frames[6] <<"EOF"
=======~====
| ShellMan |
=======~====
EOF

IFS='' read -r -d '' frames[7] <<"EOF"
========~===
| ShellmAn |
========~===
EOF

IFS='' read -r -d '' frames[8] <<"EOF"
=========~==
| ShellmaN |
=========~==
EOF

# <<<<<<<<<<<<<<<<<<<<<<<< frames <<<<<<<<<<<<<<<<<<<<<<<<

# >>>>>>>>>>>>>>>>>>>>>>>> functions >>>>>>>>>>>>>>>>>>>>>>>>

# Usage: animate framesArray interval
function animate () {
  local frames=("$@")

  ((lastIndex=${#frames[@]} - 1))
  local mode=${frames[lastIndex]}
  unset frames[lastIndex]

  ((lastIndex=${#frames[@]} - 1))
  local interval=${frames[lastIndex]}
  unset frames[lastIndex]

  # Comment out next two lines if you are using CTRL+C event handler.
  trap 'tput cnorm; echo' EXIT
  trap 'exit 127' HUP INT TERM

  tput civis # hide cursor
  tput sc # save cursor position

  tput civis # hide cursor
  tput sc # save cursor position

  index=0
  max="${#frames[@]}"
  indices=()
  direction="forward"
  forwardIndices=( $(seq 0 1 "${max}") )
  backwardIndices=( $(seq "${max}" -1 0) )

  while true; do
    if [ "${mode}" = "circular" ]; then
      direction="forward"
    elif [ "${mode}" = "pendular" ]; then
      if (( index >= max )); then
        direction="backward"
      elif (( index <= 0 )); then
        direction="forward"
      fi
    else
      echo "Wrong mode! Valid modes: circular, pendular"
      exit 255
    fi

    if [ "${direction}" = "forward" ]; then
      indices=( "${forwardIndices[@]}" )
    else
      indices=( "${backwardIndices[@]}" )
    fi
    

    for index in "${indices[@]}"; do
      tput rc # restore cursor position
      echo "${frames[$index]}"
      sleep "${interval}"
    done
  done
}

# <<<<<<<<<<<<<<<<<<<<<<<< functions <<<<<<<<<<<<<<<<<<<<<<<<

# Usage: animate framesArray interval
animate "${frames[@]}" 0.3 pendular
