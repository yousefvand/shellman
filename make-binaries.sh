#!/usr/bin/env bash

# >>>>>>>>>>>>>>>>>>>>>>>> functions >>>>>>>>>>>>>>>>>>>>>>>>

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

# >>>>>>>>>>>>>>>>>>>>>>>> Operating Systems >>>>>>>>>>>>>>>>>>>>>>>>
# Windows
windows_amd64='GOOS=windows GOARCH=amd64 go build -o bin/build-windows-amd64.exe build.go'
windows_i386='GOOS=windows GOARCH=386 go build -o bin/build-windows-i386.exe build.go'

# macOS
mac_amd64='GOOS=darwin GOARCH=amd64 go build -o bin/build-darwin-amd64 build.go'
mac_arm64='GOOS=darwin GOARCH=arm64 go build -o bin/build-darwin-arm64 build.go'
#-# GOOS=darwin GOARCH=386 go build -o bin/build-darwin-386 build.go
#-# GOOS=darwin GOARCH=arm go build -o bin/build-darwin-arm build.go

# Linux
linux_amd64='GOOS=linux GOARCH=amd64 go build -o bin/build-linux-amd64 build.go'
linux_i386='GOOS=linux GOARCH=386 go build -o bin/build-linux-386 build.go'
linux_arm='GOOS=linux GOARCH=arm go build -o bin/build-linux-arm build.go'
linux_arm64='GOOS=linux GOARCH=arm64 go build -o bin/build-linux-arm64 build.go'
linux_ppc64='GOOS=linux GOARCH=ppc64 go build -o bin/build-linux-ppc64 build.go'
linux_ppc64le='GOOS=linux GOARCH=ppc64le go build -o bin/build-linux-ppc64le build.go'
linux_mips='GOOS=linux GOARCH=mips go build -o bin/build-linux-mips build.go'
linux_mipsle='GOOS=linux GOARCH=mipsle go build -o bin/build-linux-mipsle build.go'
linux_mips64='GOOS=linux GOARCH=mips64 go build -o bin/build-linux-mips64 build.go'
linux_mips64le='GOOS=linux GOARCH=mips64le go build -o bin/build-linux-mips64le build.go'

# Dragonfly
dragonfly_amd64='GOOS=dragonfly GOARCH=amd64 go build -o bin/build-dragonfly-amd64 build.go'

# FreeBSD
freebsd_amd64='GOOS=freebsd GOARCH=arm64 go build -o bin/build-freebsd-amd64 build.go'
freebsd_i386='GOOS=freebsd GOARCH=386 go build -o bin/build-freebsd-386 build.go'
freebsd_arm='GOOS=freebsd GOARCH=arm64 go build -o bin/build-freebsd-arm build.go'

# NetBSD
netbsd_amd64='GOOS=netbsd GOARCH=amd64 go build -o bin/build-netbsd-amd64 build.go'
netbsd_i386='GOOS=netbsd GOARCH=386 go build -o bin/build-netbsd-386 build.go'
netbsd_arm='GOOS=netbsd GOARCH=arm go build -o bin/build-netbsd-arm build.go'

# OpenBSD
openbsd_amd64='GOOS=openbsd GOARCH=amd64 go build -o bin/build-openbsd-amd64 build.go'
openbsd_i386='GOOS=openbsd GOARCH=386 go build -o bin/build-openbsd-386 build.go'
openbsd_arm='GOOS=openbsd GOARCH=arm go build -o bin/build-openbsd-arm build.go'

# Plan9
plan9_amd64='GOOS=plan9 GOARCH=amd64 go build -o bin/build-plan9-amd64 build.go'
plan9_i386='GOOS=plan9 GOARCH=386 go build -o bin/build-plan9-386 build.go'

# Solaris
# GOOS=solaris GOARCH=amd64 go build -o bin/build-openbsd-amd64 build.go

# <<<<<<<<<<<<<<<<<<<<<<<< Operating Systems <<<<<<<<<<<<<<<<<<<<<<<<

rm -rf ./bin 2>/dev/null
mkdir ./bin

# Usage: multiChoice "header message" resultArray "comma separated options" "comma separated default values"
multiChoice "Compile build tool for:" result "linux-amd64;windows-amd64;mac-amd64" "1;1;1"

if [[ "${result[@]}" =~ '0' ]]; then
  eval "${linux_amd64}"
  echo $(tput setaf 4)"Compiling build tool for Linux....."$(tput sgr0)
  
fi

if [[ "${result[@]}" =~ '1' ]]; then
  eval "${windows_amd64}"
  echo $(tput setaf 3)"Compiling build tool for Windows..."$(tput sgr0)
  
fi

if [[ "${result[@]}" =~ '2' ]]; then
  eval "${mac_amd64}"
  echo $(tput setaf 5)"Compiling build tool for MacOS....."$(tput sgr0)
  
fi
