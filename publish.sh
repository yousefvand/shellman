#!/usr/bin/env bash

# >>>>>>>>>>>>>>>>>>>>>>>> functions >>>>>>>>>>>>>>>>>>>>>>>>

# Usage: bannerColor "my title" "red" "*"
function bannerColor() {
  case ${2} in
    black) color=0
    ;;
    red) color=1
    ;;
    green) color=2
    ;;
    yellow) color=3
    ;;
    blue) color=4
    ;;
    magenta) color=5
    ;;
    cyan) color=6
    ;;
    white) color=7
    ;;
    *) echo "color is not set"; exit 1
    ;;
  esac

  local msg="${3} ${1} ${3}"
  local edge=$(echo "${msg}" | sed "s/./${3}/g")
  tput setaf ${color}
  tput bold
  echo "${edge}"
  echo "${msg}"
  echo "${edge}"
  tput sgr 0
  echo
}


# Usage: multiChoice "header message" resultArray "comma separated options" "comma separated default values"
# Credit: https://serverfault.com/a/949806
# TODO: 1) Refactoring to return result array 2) Get input options as array
function multiChoice {
  echo "${1}"; shift
  echo $(tput dim)-"Change Option: [up/down], Change Selection: [space], Done: [ENTER]" $(tput sgr0)
  # little helpers for terminal print control and key input
  ESC=$( printf "\033")
  cursor_blink_on()   { printf "$ESC[?25h"; }
  cursor_blink_off()  { printf "$ESC[?25l"; }
  cursor_to()         { printf "$ESC[$1;${2:-1}H"; }
  print_inactive()    { printf "$2   $1 "; }
  print_active()      { printf "$2  $ESC[7m $1 $ESC[27m"; }
  get_cursor_row()    { IFS=';' read -sdR -p $'\E[6n' ROW COL; echo ${ROW#*[}; }
  key_input()         {
    local key
    IFS= read -rsn1 key 2>/dev/null >&2
    if [[ $key = ""      ]]; then echo enter; fi;
    if [[ $key = $'\x20' ]]; then echo space; fi;
    if [[ $key = $'\x1b' ]]; then
      read -rsn2 key
      if [[ $key = [A ]]; then echo up;    fi;
      if [[ $key = [B ]]; then echo down;  fi;
    fi
  }
  toggle_option()    {
    local arr_name=$1
    eval "local arr=(\"\${${arr_name}[@]}\")"
    local option=$2
    if [[ ${arr[option]} == 1 ]]; then
      arr[option]=0
    else
      arr[option]=1
    fi
    eval $arr_name='("${arr[@]}")'
  }

  local retval=$1
  local options
  local defaults

  IFS=';' read -r -a options <<< "$2"
  if [[ -z $3 ]]; then
    defaults=()
  else
    IFS=';' read -r -a defaults <<< "$3"
  fi

  local selected=()

  for ((i=0; i<${#options[@]}; i++)); do
    selected+=("${defaults[i]}")
    printf "\n"
  done

  # determine current screen position for overwriting the options
  local lastrow=$(get_cursor_row)
  local startrow=$(($lastrow - ${#options[@]}))

  # ensure cursor and input echoing back on upon a ctrl+c during read -s
  trap "cursor_blink_on; stty echo; printf '\n'; exit" 2
  cursor_blink_off

  local active=0
  while true; do
    # print options by overwriting the last lines
    local idx=0
    for option in "${options[@]}"; do
      local prefix="[ ]"
      if [[ ${selected[idx]} == 1 ]]; then
        prefix="[x]"
      fi

      cursor_to $(($startrow + $idx))
      if [ $idx -eq $active ]; then
        print_active "$option" "$prefix"
      else
        print_inactive "$option" "$prefix"
      fi
      ((idx++))
    done

    # user key control
    case $(key_input) in
      space)  toggle_option selected $active;;
      enter)  break;;
      up)     ((active--));
        if [ $active -lt 0 ]; then active=$((${#options[@]} - 1)); fi;;
      down)   ((active++));
        if [ $active -ge ${#options[@]} ]; then active=0; fi;;
    esac
  done

  # cursor position back to normal
  cursor_to $lastrow
  printf "\n"
  cursor_blink_on

  indices=()
  for((i=0;i<${#selected[@]};i++)); do
    if ((${selected[i]} == 1)); then
      indices+=(${i})
    fi
  done

  # eval $retval='("${selected[@]}")'
  eval $retval='("${indices[@]}")'
}

# <<<<<<<<<<<<<<<<<<<<<<<< functions <<<<<<<<<<<<<<<<<<<<<<<<

# Usage: bannerColor "my title" "red" "*"
bannerColor "Publish vscode extension to market" "magenta" ":"


# Usage: multiChoice "header message" resultArray "comma separated options" "comma separated default values"
multiChoice "Which markets to publish:" markets "vscode;open-vsx" "1;1"

if [[ "${markets[@]}" =~ '0' ]]; then
  echo "Enter vscode token:"
  read -s vscodeToken
  echo $(tput setaf 6)"Publishing to vscode..."$(tput sgr0)
  
  npx vsce publish -p "${vscodeToken}"
fi

if [[ "${markets[@]}" =~ '1' ]]; then
  echo "Enter open-vsx token:"
  read -s openVsxToken
  echo $(tput setaf 6)"Publishing to open-vsx..."$(tput sgr0)

  npx ovsx publish -p "${openVsxToken}"
fi

echo "DONE!"