# Commands

- [bash](#bash)

- [summary](#summary)

- [func](#func)

- [func args count](#func-args-count)

- [func args](#func-args)

- [func ret val](#func-ret-val)

- [cmd](#cmd)

- [cmd success check](#cmd-success-check)

- [cmd failure check](#cmd-failure-check)

- [assign array](#assign-array)

- [assign if empty](#assign-if-empty)

- [ask question](#ask-question)

- [directory create](#directory-create)

- [directory create nested](#directory-create-nested)

- [if](#if)

- [iff](#iff)

- [iff not](#iff-not)

- [if string empty](#if-string-empty)

- [if string not empty](#if-string-not-empty)

- [if string =](#if-string-=)

- [if string !=](#if-string-!=)

- [string length](#string-length)

- [string trim](#string-trim)

- [string trim left](#string-trim-left)

- [string trim right](#string-trim-right)

- [string trim all](#string-trim-all)

- [string replace](#string-replace)

- [string reverse](#string-reverse)

- [string toLower](#string-toLower)

- [string toUpper](#string-toUpper)

- [string substring](#string-substring)

- [string contains](#string-contains)

- [string indexOf](#string-indexOf)

- [if int =](#if-int-=)

- [if int !=](#if-int-!=)

- [if int >](#if-int->)

- [if int >=](#if-int->=)

- [if int <](#if-int-<)

- [if int <=](#if-int-<=)

- [if cmd exists](#if-cmd-exists)

- [if exists](#if-exists)

- [if file exists](#if-file-exists)

- [if file not empty](#if-file-not-empty)

- [if directory exists](#if-directory-exists)

- [if file readable](#if-file-readable)

- [if file writeable](#if-file-writeable)

- [if file executable](#if-file-executable)

- [if file =](#if-file-=)

- [if file link](#if-file-link)

- [if file newer](#if-file-newer)

- [if file older](#if-file-older)

- [ips](#ips)

- [ip info](#ip-info)

- [ip public](#ip-public)

- [for i](#for-i)

- [for ij](#for-ij)

- [for in](#for-in)

- [for in column](#for-in-column)

- [while](#while)

- [until](#until)

- [switch case](#switch-case)

- [expr](#expr)

- [math](#math)

- [file read](#file-read)

- [file write](#file-write)

- [file write multiline](#file-write-multiline)

- [file write multiline sudo](#file-write-multiline-sudo)

- [remove files older](#remove-files-older)

- [color black](#color-black)

- [color red](#color-red)

- [color green](#color-green)

- [color yellow](#color-yellow)

- [color blue](#color-blue)

- [color magenta](#color-magenta)

- [color cyan](#color-cyan)

- [color white](#color-white)

- [format bold](#format-bold)

- [format italic](#format-italic)

- [format dim](#format-dim)

- [format reverse](#format-reverse)

- [random int](#random-int)

- [region](#region)

- [service manage](#service-manage)

- [stopwatch start](#stopwatch-start)

- [stopwatch stop](#stopwatch-stop)

- [stopwatch elapsed](#stopwatch-elapsed)

- [sleep](#sleep)

- [timeout](#timeout)

- [fn banner simple](#fn-banner-simple)

- [fx banner simple](#fx-banner-simple)

- [fn banner color](#fn-banner-color)

- [fx banner color](#fx-banner-color)

- [fn scan](#fn-scan)

- [fx scan](#fx-scan)

## `bash`

bash shebang [&uarr;](#Commands)

```bash
#!/usr/bin/env ${1|bash,node,perl,php,python,python3,ruby|}
```

## `summary`

Script summary [&uarr;](#Commands)

```bash
# Title:         title
# Description:   description
# Author:        author <email>
# Date:          yyyy-mm-dd
# Version:       1.0.0
```

## `func`

function [&uarr;](#Commands)

```bash
function name () {
  echo "$1" # arguments are accessible through $1, $2,...
}
```

## `func args count`

number of function arguments [&uarr;](#Commands)

```bash
echo $#
```

## `func args`

function arguments array [&uarr;](#Commands)

```bash
echo "$@"
```

## `func ret val`

last function/command return code [&uarr;](#Commands)

```bash
echo "$?"
```

## `cmd`

run command (command substitution) [&uarr;](#Commands)

```bash
`command`
```

## `cmd success check`

check if last command succeed [&uarr;](#Commands)

```bash
if [[ $? == 0 ]]; then
  echo command succeed
fi
```

## `cmd failure check`

check if last command failed [&uarr;](#Commands)

```bash
if [[ $? != 0 ]]; then
  echo command failed
fi
```

## `assign array`

assign elements to an array [&uarr;](#Commands)

```bash
newArray=("Element 1" "Element 2")
```

## `assign if empty`

assign default to variable if variable is empty or null [&uarr;](#Commands)

```bash
: "${variable:=default}"
```

## `ask question`

Ask question with default answer [&uarr;](#Commands)

```bash
read -ep "Question here? " -i Default answer ANSWER
echo "$ANSWER"
```

## `directory create`

create directory [&uarr;](#Commands)

```bash
mkdir dirname
```

## `directory create nested`

create nested directories [&uarr;](#Commands)

```bash
mkdir -p parent dir/child dir
```

## `if`

if [&uarr;](#Commands)

```bash
if [ condition ]; then
   # body
elif [ condition ]; then
   # body
else
   # body
fi
```

## `iff`

if condition is true then run command (short circuit) [&uarr;](#Commands)

```bash
[ condition ] && command
```

## `iff not`

if condition is false then run command (short circuit) [&uarr;](#Commands)

```bash
[ condition ] || command
```

## `if string empty`

if string is empty [&uarr;](#Commands)

```bash
if [ -z "$string" ]; then
  # body
fi
```

## `if string not empty`

if string is not empty [&uarr;](#Commands)

```bash
if [ -n "$string" ]; then
  # body
fi
```

## `if string =,string equal`

if strings are equal [&uarr;](#Commands)

```bash
if [ "$string1" = "$string2" ]; then
  # body
fi
```

## `if string !=,string not equal`

if strings are not equal [&uarr;](#Commands)

```bash
if [ "$string1" != "$string2" ]; then
  # body
fi
```

## `string length`

length of string in characters [&uarr;](#Commands)

```bash
length=${#variable}
```

## `string trim`

remove leading and trailing white space(s) [&uarr;](#Commands)

```bash
trimmed=`echo -e "${var}" |  sed -e 's/^[[:space:]]*//' | sed -e 's/[[:space:]]*$//'`
```

## `string trim left`

remove leading white space(s) [&uarr;](#Commands)

```bash
trimmed=`echo -e "${var}" | sed -e 's/^[[:space:]]*//'`
```

## `string trim right`

remove trailing white space(s) [&uarr;](#Commands)

```bash
trimmed=`echo -e "${var}" | sed -e 's/[[:space:]]*$//'`
```

## `string trim all`

remove all white space(s) [&uarr;](#Commands)

```bash
trimmed=`echo -e "${var}" | tr -d '[[:space:]]'`
```

## `string replace`

find all occurrences of a substrings and replace them [&uarr;](#Commands)

```bash
replaced=`echo -e "${var}" | sed -e 's/find/replace/g'`
```

## `string reverse`

reverse string characters [&uarr;](#Commands)

```bash
reversed=`echo -e "${var}" | rev`
```

## `string toLower`

convert string to lowercase [&uarr;](#Commands)

```bash
toLower=`echo -e "${var}" | tr '[:upper:]' '[:lower:]'`
```

## `string toUpper`

convert string to uppercase [&uarr;](#Commands)

```bash
toUpper=`echo -e "${var}" | tr '[:lower:]' '[:upper:]'`
```

## `string substring`

part of the string from offset with length characters [&uarr;](#Commands)

```bash
substring=`echo -e "${var:offset:length"`
```

## `string contains,if string contains`

check whether string contains substring [&uarr;](#Commands)

```bash
if [[ "$string" = *substring* ]]; then
  # body
fi
```

## `string indexOf`

first index of substring in string [&uarr;](#Commands)

```bash
temp=${string%%"substring"*} && indexOf=`echo ${string%%"substring"*} | echo ${#temp}`
# echo $indexOf
```

## `if int =`

if integers are equal [&uarr;](#Commands)

```bash
if [[ $int1 == $int2 ]]; then
  echo equal
fi
```

## `if int !=`

if integers are not equal [&uarr;](#Commands)

```bash
if [[ $int1 != $int2 ]]; then
  echo not equal
fi
```

## `if int >`

if integer greater than value [&uarr;](#Commands)

```bash
if [[ $int > $${2:val} ]]; then
  echo greater
fi
```

## `if int >=`

if integer greater than or equal value [&uarr;](#Commands)

```bash
if [[ $int >= $val ]]; then
  echo greater equal
fi
```

## `if int <`

if integer lesser than value [&uarr;](#Commands)

```bash
if [[ $int > $${2:val} ]]; then
  echo lesser
fi
```

## `if int <=`

if integer lesser than or equal value [&uarr;](#Commands)

```bash
if [[ $int <= $val ]]; then
  echo lesser equal
fi
```

## `if cmd exists`

check if command exists [&uarr;](#Commands)

```bash
if [ `command -v command` ]; then
  # body
fi
```

## `if exists`

if path exists (file, directory, link...) [&uarr;](#Commands)

```bash
if [ -e "$path"; then
  echo exists
fi
```

## `if file exists`

if file exists [&uarr;](#Commands)

```bash
if [ -f "$file" ]; then
  echo file exists
fi
```

## `if file not empty`

if file size is greater than zero [&uarr;](#Commands)

```bash
if [ -s "$file" ]; then
  echo file not empty
fi
```

## `if directory exists`

if directory exists [&uarr;](#Commands)

```bash
if [ -d "$directory" ]; then
  echo directory exists
fi
```

## `if file readable`

if file readable [&uarr;](#Commands)

```bash
if [ -r "$file" ]; then
  echo file is readable
fi
```

## `if file writeable`

if file writeable [&uarr;](#Commands)

```bash
if [ -w "$file" ]; then
  echo file is writeable
fi
```

## `if file executable`

if file executable [&uarr;](#Commands)

```bash
if [ -x "$file" ]; then
  echo file is executable
fi
```

## `if file =`

if files are equal [&uarr;](#Commands)

```bash
if [ "$file1" -ef "$file2" ]; then
  echo files are equal
fi
```

## `if file link`

if file exists and is a symbolic link [&uarr;](#Commands)

```bash
if [ -h "$file" ]; then
  echo symbolic link
fi
```

## `if file newer`

if file1 newer than file2 [&uarr;](#Commands)

```bash
if [ "$file1" -nt "$file2" ]; then
  echo file1 is newer than file2
fi
```

## `if file older`

if file1 older than file2 [&uarr;](#Commands)

```bash
if [ "$file1" -ot "$file2" ]; then
  echo file1 is older than file2
fi
```

## `ips`

Array of local IPs [&uarr;](#Commands)

```bash
IPS=`hostname -I`
```

## `ip info`

public ip information [&uarr;](#Commands)

```bash
echo `curl -s ipinfo.io/${1|ip,city,region,country,loc,postal,org|}`
```

## `ip public`

public ip address [&uarr;](#Commands)

```bash
PUBLIC_IP=`curl -s ${1|bot.whatismyipaddress.com,ident.me,ipecho.net/plain,icanhazip.com,ifconfig.me,api.ipify.org,ipinfo.io/ip|}`
```

## `for i`

for loop by index [&uarr;](#Commands)

```bash
for((i=0;i<n;i++)); do
  echo "$i"
done
```

## `for ij`

for loop by index [&uarr;](#Commands)

```bash
for((i=0;i<n;i++)); do
  for((j=0;j<m;j++)); do
    echo "$i, $j"
  done
done
```

## `for in`

for loop in collection [&uarr;](#Commands)

```bash
for item in {a..z}; do
  echo "$item"
done
```

## `for in column`

for loop in collection [&uarr;](#Commands)

```bash
for col in `docker images | awk '{ print $1":"$2 }'`; do
  echo "$col" | cut -d ':' -f 1
  echo "$col" | cut -d ':' -f 2
done
```

## `while`

while loop [&uarr;](#Commands)

```bash
while [ condition ]; do
  # body
done
```

## `until`

until loop [&uarr;](#Commands)

```bash
until [ condition ]; do
  # body
done
```

## `switch case`

switch case [&uarr;](#Commands)

```bash
case "$item" in
  1)
    echo "case 1"
  ;;
  2|3)
    echo "case 2 or 3"
  ;;
  *)
    echo "default"
  ;;
esac
```

## `expr`

arithmetic operations [&uarr;](#Commands)

```bash
expr 2 * 3
```

## `math`

Doing mathematical operations [&uarr;](#Commands)

```bash
$((i++))
```

## `file read`

read a file [&uarr;](#Commands)

```bash
cat "$filepath" | while read line; do
  echo "$line"
done
```

## `file write`

write a file [&uarr;](#Commands)

```bash
echo "file header" > ${2:"$filepath"}
for line in {lines}; do
  echo "$line" >> "$filepath"
done
```

## `file write multiline`

write multiple lines into file [&uarr;](#Commands)

```bash
cat >filepath <<EOL
# text here
EOL
```

## `file write multiline sudo`

write multiple lines into file with sudo permission [&uarr;](#Commands)

```bash
cat << EOL | sudo tee filepath
# text here
EOL
```

## `remove files older`

find and remove files older than x days [&uarr;](#Commands)

```bash
find "$path" -mtime +days | xargs rm -f
```

## `color black`

write in black [&uarr;](#Commands)

```bash
echo `tput setaf 0`black text`tput sgr0`
```

## `color red`

write in red [&uarr;](#Commands)

```bash
echo `tput setaf 1`red text`tput sgr0`
```

## `color green`

write in green [&uarr;](#Commands)

```bash
echo `tput setaf 2`green text`tput sgr0`
```

## `color yellow`

write in yellow [&uarr;](#Commands)

```bash
echo `tput setaf 3`yellow text`tput sgr0`
```

## `color blue`

write in blue [&uarr;](#Commands)

```bash
echo `tput setaf 4`blue text`tput sgr0`
```

## `color magenta`

write in magenta [&uarr;](#Commands)

```bash
echo `tput setaf 5`magenta text`tput sgr0`
```

## `color cyan`

write in cyan [&uarr;](#Commands)

```bash
echo `tput setaf 6`cyan text`tput sgr0`
```

## `color white`

write in white [&uarr;](#Commands)

```bash
echo `tput setaf 7`white text`tput sgr0`
```

## `format bold`

write in bold [&uarr;](#Commands)

```bash
echo `tput bold`bold text`tput sgr0`
```

## `format italic`

write in italic [&uarr;](#Commands)

```bash
echo `tput sitm`italic text`tput sgr0`
```

## `format dim`

write in dim [&uarr;](#Commands)

```bash
echo `tput dim`dimmed text`tput sgr0`
```

## `format reverse`

write in reverse [&uarr;](#Commands)

```bash
echo `tput rev`reversed text`tput sgr0`
```

## `random int`

generate random integer x such as min < x < max [&uarr;](#Commands)

```bash
echo $((min + RANDOM % $((max-min))))
```

## `region`

Comment out a special region (i.e. variable declrations [&uarr;](#Commands)

```bash
# >>>>>>>>>>>>>>>>>>>>>>>> name >>>>>>>>>>>>>>>>>>>>>>>>

# <<<<<<<<<<<<<<<<<<<<<<<< name <<<<<<<<<<<<<<<<<<<<<<<<
```

## `service manage`

Manage service operations [&uarr;](#Commands)

```bash
sudo systemctl ${1|enable,disable,start,stop,reload,restart,status|} service
```

## `stopwatch start`

start stopwatch [&uarr;](#Commands)

```bash
STOPWATCH_START_TIME=$(date +%s)
```

## `stopwatch stop`

stop stopwatch [&uarr;](#Commands)

```bash
STOPWATCH_END_TIME=$(date +%s)
```

## `stopwatch elapsed`

elapsed time [&uarr;](#Commands)

```bash
STOPWATCH_ELAPSED_TOTAL_SECONDS=$((STOPWATCH_END_TIME - STOPWATCH_START_TIME))
STOPWATCH_ELAPSED_MINUTES=$((STOPWATCH_ELAPSED_TOTAL_SECONDS / 60))
STOPWATCH_ELAPSED_SECONDS=$((STOPWATCH_ELAPSED_TOTAL_SECONDS % 60))
echo elapsed $STOPWATCH_ELAPSED_MINUTES minutes and $STOPWATCH_ELAPSED_SECONDS seconds
```

## `sleep`

sleep for a specified amount of time (s: second, m: minute, h: hour, d: day) [&uarr;](#Commands)

```bash
sleep 30${2|s,m,h,d|}
```

## `timeout`

Run command within a time frame [&uarr;](#Commands)

```bash
timeout seconds command
```

## `fn banner simple`

function: print a banner. example: banner_simple "my title" [&uarr;](#Commands)

```bash
function banner_simple() {
  local msg="* $* *"
  local edge=`echo "$msg" | sed 's/./*/g'`
  echo "$edge"
  echo "`tput bold`$msg`tput sgr0`"
  echo "$edge"
  echo
}
```

## `fx banner simple`

call banner_simple function [&uarr;](#Commands)

```bash
banner_simple "my title"
```

## `fn banner color`

function: print a color banner. example: banner_color red "my title"
 [&uarr;](#Commands)

```bash
function banner_color() {
  local color=$1
  shift
  case $color in
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
  local s=("$@") b w
  for l in "${s[@]}"; do
    ((w<${#l})) && { b="$l"; w="${#l}"; }
  done
  tput setaf $color
  echo " =${b//?/=}=
| ${b//?/ } |"
  for l in "${s[@]}"; do
    printf '| %s%*s%s |\n' "$(tput setaf $color)" "-$w" "$l" "$(tput setaf $color)"
  done
  echo "| ${b//?/ } |
 =${b//?/=}="
  tput sgr 0
}
```

## `fx banner color`

call banner_color function [&uarr;](#Commands)

```bash
banner_color ${1|black,red,green,yellow,blue,magenta,cyan,white|} "my title"
```

## `fn scan`

Scan host's port range (tcp/udp) [&uarr;](#Commands)

```bash
# scan proto host fromPort toPort
function scan () {
  for ((port=$3; port<=$4; port++)); do
    (echo >/dev/$1/$2/$port) >/dev/null 2>&1 && echo "$1 $port => open"
  done
}
```

## `fx scan`

call scan function to scan a host over a port range [&uarr;](#Commands)

```bash
scan ${1|tcp,udp|} host fromPort  toPort
```
