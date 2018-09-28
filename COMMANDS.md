# Commands

- [bash](##bash)

- [summary](##summary)

- [func](##func)

- [func args count](##func&#32;args&#32;count)

- [func args](##func&#32;args)

- [func ret val](##func&#32;ret&#32;val)

- [cmd](##cmd)

- [ask question](##ask&#32;question)

- [if](##if)

- [iff](##iff)

- [iff not](##iff&#32;not)

- [if string empty](##if&#32;string&#32;empty)

- [if string not empty](##if&#32;string&#32;not&#32;empty)

- [if string =](##if&#32;string&#32;=)

- [if string !=](##if&#32;string&#32;!=)

- [if int =](##if&#32;int&#32;=)

- [if int !=](##if&#32;int&#32;!=)

- [if int >](##if&#32;int&#32;>)

- [if int >=](##if&#32;int&#32;>=)

- [if int <](##if&#32;int&#32;<)

- [if int <=](##if&#32;int&#32;<=)

- [if cmd exists](##if&#32;cmd&#32;exists)

- [if exists](##if&#32;exists)

- [if file exists](##if&#32;file&#32;exists)

- [if file not empty](##if&#32;file&#32;not&#32;empty)

- [if directory exists](##if&#32;directory&#32;exists)

- [if file readable](##if&#32;file&#32;readable)

- [if file writeable](##if&#32;file&#32;writeable)

- [if file executable](##if&#32;file&#32;executable)

- [if file =](##if&#32;file&#32;=)

- [if file link](##if&#32;file&#32;link)

- [if file newer](##if&#32;file&#32;newer)

- [if file older](##if&#32;file&#32;older)

- [for i](##for&#32;i)

- [for ij](##for&#32;ij)

- [for in](##for&#32;in)

- [for in column](##for&#32;in&#32;column)

- [while](##while)

- [until](##until)

- [switch case](##switch&#32;case)

- [expr](##expr)

- [math](##math)

- [file read](##file&#32;read)

- [file write](##file&#32;write)

- [file write multiline](##file&#32;write&#32;multiline)

- [file write multiline sudo](##file&#32;write&#32;multiline&#32;sudo)

- [remove files older](##remove&#32;files&#32;older)

- [color black](##color&#32;black)

- [color red](##color&#32;red)

- [color green](##color&#32;green)

- [color yellow](##color&#32;yellow)

- [color blue](##color&#32;blue)

- [color magenta](##color&#32;magenta)

- [color cyan](##color&#32;cyan)

- [color white](##color&#32;white)

- [format bold](##format&#32;bold)

- [format italic](##format&#32;italic)

- [format dim](##format&#32;dim)

- [format reverse](##format&#32;reverse)

## `bash`

bash shebang [&uarr;](#Commands)

```bash
#!/usr/bin/env bash
```

## `summary`

Script summary [&uarr;](#Commands)

```bash
# Title:         title
# Description:   description
# Author:        author
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

## `ask question`

Ask question with default answer [&uarr;](#Commands)

```bash
read -ep "Question here? " -i Default answer ANSWER
echo "$ANSWER"
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

## `if string =`

if strings are equal [&uarr;](#Commands)

```bash
if [ "$string1" = "$string2" ]; then
  # body
fi
```

## `if string !=`

if strings are not equal [&uarr;](#Commands)

```bash
if [ "$string1" != "$string2" ]; then
  # body
fi
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
\$((i++))
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
find "\$path" -mtime +days | xargs rm -f
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

