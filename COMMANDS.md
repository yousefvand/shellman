# Commands

- [bash](#bash)

- [summary](#summary)

- [argument parsing](#argument-parsing)

- [func](#func)

- [func args count](#func-args-count)

- [func args](#func-args)

- [func ret val](#func-ret-val)

- [cmd](#cmd)

- [cmd nice](#cmd-nice)

- [cmd success check](#cmd-success-check)

- [cmd failure check](#cmd-failure-check)

- [archive compress tar.gz](#archive-compress-tar.gz)

- [archive decompress tar.gz](#archive-decompress-tar.gz)

- [array concat](#array-concat)

- [array declare](#array-declare)

- [array delete](#array-delete)

- [array delete at](#array-delete-at)

- [array all](#array-all)

- [array filter](#array-filter)

- [array push](#array-push)

- [array replace](#array-replace)

- [array at index](#array-at-index)

- [array iterate](#array-iterate)

- [array length](#array-length)

- [array slice](#array-slice)

- [array set element](#array-set-element)

- [assign if empty](#assign-if-empty)

- [ask question](#ask-question)

- [crypto base64 encode](#crypto-base64-encode)

- [crypto base64 decode](#crypto-base64-decode)

- [crypto hash](#crypto-hash)

- [date now short](#date-now-short)

- [date now UTC](#date-now-UTC)

- [date now year](#date-now-year)

- [date now monthNumber](#date-now-monthNumber)

- [date now monthName](#date-now-monthName)

- [date now dayOfMonth](#date-now-dayOfMonth)

- [date now dayOfWeek](#date-now-dayOfWeek)

- [date now dayOfYear](#date-now-dayOfYear)

- [time now local](#time-now-local)

- [time now UTC](#time-now-UTC)

- [time seconds epoch](#time-seconds-epoch)

- [directory create](#directory-create)

- [directory create nested](#directory-create-nested)

- [git branch list](#git-branch-list)

- [git branch create](#git-branch-create)

- [git branch delete local](#git-branch-delete-local)

- [git branch delete remote](#git-branch-delete-remote)

- [git branch push](#git-branch-push)

- [git branch rename](#git-branch-rename)

- [git changes revert](#git-changes-revert)

- [git clone](#git-clone)

- [git clone branch](#git-clone-branch)

- [git commit](#git-commit)

- [git commit list notPushed](#git-commit-list-notPushed)

- [git commit search](#git-commit-search)

- [git commit undo](#git-commit-undo)

- [git config list](#git-config-list)

- [git config set](#git-config-set)

- [git patch apply](#git-patch-apply)

- [git patch create](#git-patch-create)

- [git remote list](#git-remote-list)

- [git remote urlAdd](#git-remote-urlAdd)

- [git remote urlChange](#git-remote-urlChange)

- [git tag list](#git-tag-list)

- [git tag commit](#git-tag-commit)

- [git tag remote delete](#git-tag-remote-delete)

- [git tag remote push](#git-tag-remote-push)

- [http GET](#http-GET)

- [http POST](#http-POST)

- [http POST file](#http-POST-file)

- [http header](#http-header)

- [http cookie](#http-cookie)

- [http download](#http-download)

- [ftp list](#ftp-list)

- [ftp download](#ftp-download)

- [ftp upload](#ftp-upload)

- [ftp delete file](#ftp-delete-file)

- [ftp rename](#ftp-rename)

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

- [let](#let)

- [expr](#expr)

- [math ++](#math-++)

- [math --](#math---)

- [math +](#math-+)

- [math -](#math--)

- [math *](#math-*)

- [math /](#math-/)

- [math %](#math-%)

- [math ^](#math-^)

- [math random](#math-random)

- [math √](#math-√)

- [math 0.00](#math-0.00)

- [math const π](#math-const-π)

- [math const e](#math-const-e)

- [math const 𝛾](#math-const-𝛾)

- [math const ϕ](#math-const-ϕ)

- [math const Ω](#math-const-Ω)

- [file delete](#file-delete)

- [file find](#file-find)

- [file read](#file-read)

- [file search](#file-search)

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

- [region](#region)

- [service manage](#service-manage)

- [stopwatch start](#stopwatch-start)

- [stopwatch stop](#stopwatch-stop)

- [stopwatch elapsed](#stopwatch-elapsed)

- [sleep](#sleep)

- [system memory info](#system-memory-info)

- [system distro name](#system-distro-name)

- [system distro version](#system-distro-version)

- [system distro codename](#system-distro-codename)

- [system kernel name](#system-kernel-name)

- [system kernel release](#system-kernel-release)

- [system processor type](#system-processor-type)

- [system processor count](#system-processor-count)

- [system processor architecture](#system-processor-architecture)

- [system processor model](#system-processor-model)

- [system uptime](#system-uptime)

- [timeout](#timeout)

- [fn banner simple](#fn-banner-simple)

- [fx banner simple](#fx-banner-simple)

- [fn banner color](#fn-banner-color)

- [fx banner color](#fx-banner-color)

- [fn import](#fn-import)

- [fx import](#fx-import)

- [fn options](#fn-options)

- [fx options](#fx-options)

- [fn progress](#fn-progress)

- [fx progress](#fx-progress)

- [fn scan](#fn-scan)

- [fx scan](#fx-scan)

## `bash`

bash shebang [&uarr;](#Commands)

```bash
#!/usr/bin/env ${1|bash,node,perl,php,python,python3,ruby|}
```

## `summary`

script summary [&uarr;](#Commands)

```bash
# Title:         title
# Description:   description
# Author:        author <email>
# Date:          yyyy-mm-dd
# Version:       1.0.0
```

## `argument parsing,parse args`

parse command line arguments (flags/switches) [&uarr;](#Commands)

```bash
POSITIONAL=()
while [[ $# > 0 ]]; do
  case "$1" in
    -f|--flag)
    echo flag: $1
    shift # shift once since flags have no values
    ;;
    -s|--switch)
    echo switch $1 with value: $2
    shift 2 # shift twice to bypass switch and its value
    ;;
    *) # unknown flag/switch
    POSITIONAL+=("$1")
    shift
    ;;
  esac
done
set -- "${POSITIONAL[@]}" # restore positional params
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

## `cmd nice`

run command with desired privilege. n: -20 (most favorable scheduling) to 19 (least favorable) [&uarr;](#Commands)

```bash
sudo nice -n ${1|-20,-15,-10,-5,0,5,10,15,19|} command
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

## `archive compress tar.gz`

compress file/folder to a .tar.gz file [&uarr;](#Commands)

```bash
tar -czvf /path/to/archive.tar.gz /path/to/directory-or-file
```

## `archive decompress tar.gz`

compress file/folder to a .tar.gz file [&uarr;](#Commands)

```bash
tar -C /extract/to/path -xzvf /path/to/archive.tar.gz
```

## `array concat`

concatenate two arrays [&uarr;](#Commands)

```bash
newArray=("${array1[@]}" "${array2[@]}")
```

## `array declare`

declare an array [&uarr;](#Commands)

```bash
myArray=('one' 'two' 'three')
```

## `array delete`

delete entire array [&uarr;](#Commands)

```bash
unset myArray
```

## `array delete at`

delete element at index from array [&uarr;](#Commands)

```bash
unset myArray[index]
```

## `array all`

all array elements [&uarr;](#Commands)

```bash
${myArray[@]}
```

## `array filter`

filter elements of an array based on given pattern [&uarr;](#Commands)

```bash
filtered=(`for i in ${myArray[@]} ; do echo $i; done | grep pattern`)
```

## `array push,array add`

push new item to the end of array [&uarr;](#Commands)

```bash
myArray+=('newItem')
```

## `array replace`

find and replace elements in array using regex [&uarr;](#Commands)

```bash
${myArray[@]//find/replace}
```

## `array at index`

retrieve element at specified index (zero based) [&uarr;](#Commands)

```bash
${myArray[index]}
```

## `array iterate,array forEach`

iterate array elements [&uarr;](#Commands)

```bash
for item in ${myArray[@]}; do
  echo "$item"
done
```

## `array length`

length of array [&uarr;](#Commands)

```bash
${#myArray[@]}
```

## `array slice,array range`

elements of array from index, equal to count numbers [&uarr;](#Commands)

```bash
${myArray[@]:from:count}
```

## `array set element`

set array element at index [&uarr;](#Commands)

```bash
myArray[index]=value
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

## `crypto base64 encode`

encode variable to base64 [&uarr;](#Commands)

```bash
base64Encoded=`echo -n "$variableToEncode" | base64`
```

## `crypto base64 decode`

decode variable from base64 [&uarr;](#Commands)

```bash
base64Decoded=`echo -n "$variableToDecode" | base64 -d`
```

## `crypto hash`

compute hash of variable (md5, sha, sha1, sha224, sha256, sha384, sha512) [&uarr;](#Commands)

```bash
hash=`echo -n "$variableToHash" | ${3|md5sum,shasum,sha1sum,sha224sum,sha256sum,sha384sum,sha512sum|} | cut -f1 -d ' '`
```

## `date now short`

yyyy/mm/dd [&uarr;](#Commands)

```bash
dateShort=`date -I`
```

## `date now UTC`

coordinated Universal Time [&uarr;](#Commands)

```bash
dateUTC=`date -u`
```

## `date now year`

current Year [&uarr;](#Commands)

```bash
year=`date +%Y`
```

## `date now monthNumber`

current month number (1..12) [&uarr;](#Commands)

```bash
monthNumber=`date +%m`
```

## `date now monthName`

current month name (full/abbreviated B/b) [&uarr;](#Commands)

```bash
monthName=`date +%${1|B,b|}`
```

## `date now dayOfMonth`

current day of month (1..31) [&uarr;](#Commands)

```bash
dayOfMonth=`date +%d`
```

## `date now dayOfWeek`

current day of week name (full/abbreviated A/a) [&uarr;](#Commands)

```bash
dayOfWeek=`date +%${1|A,a|}`
```

## `date now dayOfYear`

current day of year (1..366) [&uarr;](#Commands)

```bash
dayOfYear=`date +%j`
```

## `time now local`

current local time (24hrs/12hrs R/r) [&uarr;](#Commands)

```bash
timeNowLocal=`date +%${1|R,r|}`
```

## `time now UTC`

current UTC time [&uarr;](#Commands)

```bash
timeNowUTC=`date -u +%R`
```

## `time seconds epoch`

seconds since epoch (1970-01-01 00:00:00) [&uarr;](#Commands)

```bash
timeNowSecondsEpoch=`date +%s`
```

## `directory create`

create directory [&uarr;](#Commands)

```bash
mkdir "dirname"
```

## `directory create nested`

create nested directories [&uarr;](#Commands)

```bash
mkdir -p "parent dir"/"child dir"
```

## `git branch list`

List all branches. [&uarr;](#Commands)

```bash
git branch
```

## `git branch create`

Create branch locally and switch into it. [&uarr;](#Commands)

```bash
git checkout -b branchName
```

## `git branch delete local`

Delete local branch. [&uarr;](#Commands)

```bash
git branch --delete localBranch
```

## `git branch delete remote`

Delete remote branch. [&uarr;](#Commands)

```bash
git push origin --delete remoteBranch
```

## `git branch push`

Push branch to remote. [&uarr;](#Commands)

```bash
git push origin branchName
```

## `git branch rename`

Rename current branch. [&uarr;](#Commands)

```bash
git branch -m newName
```

## `git changes revert`

Revert tracked changes [&uarr;](#Commands)

```bash
git checkout .
```

## `git clone`

Clone remote repository to local machine. [&uarr;](#Commands)

```bash
git clone https://github.com/user/repository.git
```

## `git clone branch`

Clone a remote branch to local machine. [&uarr;](#Commands)

```bash
git clone -b branchName https://github.com/user/repository.git
```

## `git commit`

Commit changes. [&uarr;](#Commands)

```bash
git commit -m "commit message"
```

## `git commit list notPushed`

List non pushed commits. [&uarr;](#Commands)

```bash
git log origin/master..HEAD
```

## `git commit search`

Search for a commit which contains searchCriteria. [&uarr;](#Commands)

```bash
git log --all --grep='searchCriteria'
```

## `git commit undo`

Undo last N commits (soft: preserve local changes | hard: delete local changes [&uarr;](#Commands)

```bash
git reset --${1|soft,hard|} HEAD~${2|1,2,3,4,5|}
```

## `git config list`

List git configurations. [&uarr;](#Commands)

```bash
git config --list
```

## `git config set`

Configure git. [&uarr;](#Commands)

```bash
git config --${1|local,global|} ${2|user.name,user.email|} "value"
```

## `git patch apply`

Apply a patch from file. [&uarr;](#Commands)

```bash
git apply < /path/to/patch1.patch
```

## `git patch create`

Create a patch from changes. [&uarr;](#Commands)

```bash
git diff > /path/to/patch1.patch
```

## `git remote list`

List all remotes. [&uarr;](#Commands)

```bash
git remote
```

## `git remote urlAdd`

Add remote url. [&uarr;](#Commands)

```bash
git remote add origin https://github.com/user/repository.git
```

## `git remote urlChange`

Change remote url. [&uarr;](#Commands)

```bash
git remote set-url origin https://github.com/user/repository.git
```

## `git tag list`

List all tags. [&uarr;](#Commands)

```bash
git tag
```

## `git tag commit,git commit tag`

Tag a commit. [&uarr;](#Commands)

```bash
git tag -a release/1.0.0 -m "1.0.0 release"
```

## `git tag remote delete`

Delete tag from remote. [&uarr;](#Commands)

```bash
git push --delete origin tagName && git push origin :tagName
```

## `git tag remote push`

Push tag to remote. [&uarr;](#Commands)

```bash
git push origin tagName
```

## `http GET,http DELETE`

Send http GET/DELETE request using curl [&uarr;](#Commands)

```bash
curl --request ${1|GET,DELETE|} -sL \
  --user-agent '${2|Shellman,Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML\, like Gecko) Chrome/58.0.3029.110 Safari/537.36,Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:53.0) Gecko/20100101 Firefox/53.0,Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML\, like Gecko) Chrome/51.0.2704.79 Safari/537.36 Edge/14.14393|}' \
  --url 'http://example.com'
```

## `http POST,http PUT`

Send data with http POST/PUT, using curl [&uarr;](#Commands)

```bash
curl --request ${1|POST,PUT|} -sL \
  --user-agent '${2|Shellman,Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML\, like Gecko) Chrome/58.0.3029.110 Safari/537.36,Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:53.0) Gecko/20100101 Firefox/53.0,Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML\, like Gecko) Chrome/51.0.2704.79 Safari/537.36 Edge/14.14393|}' \
  --url 'http://example.com' \
  --data 'key=${5:value}' \
  --data 'key=${7:value}'
```

## `http POST file`

Send file with http POST, using curl [&uarr;](#Commands)

```bash
curl --request POST -sL \
  --user-agent '${1|Shellman,Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML\, like Gecko) Chrome/58.0.3029.110 Safari/537.36,Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:53.0) Gecko/20100101 Firefox/53.0,Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML\, like Gecko) Chrome/51.0.2704.79 Safari/537.36 Edge/14.14393|}' \
  --url 'http://example.com' \
  --form 'key=${4:value}' \
  --form 'file=@/path/to/file'
```

## `http header`

Send http request with custom header, using curl [&uarr;](#Commands)

```bash
curl --request ${1|GET,POST,DELETE,PUT|} -sL \
  --user-agent '${2|Shellman,Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML\, like Gecko) Chrome/58.0.3029.110 Safari/537.36,Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:53.0) Gecko/20100101 Firefox/53.0,Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML\, like Gecko) Chrome/51.0.2704.79 Safari/537.36 Edge/14.14393|}' \
  --header 'key: ${4:value}' \
  --url 'http://example.com'
```

## `http cookie`

Send http request with cookies, using curl [&uarr;](#Commands)

```bash
curl --request ${1|GET,POST,DELETE,PUT|} -sL \
  --user-agent '${2|Shellman,Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML\, like Gecko) Chrome/58.0.3029.110 Safari/537.36,Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:53.0) Gecko/20100101 Firefox/53.0,Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML\, like Gecko) Chrome/51.0.2704.79 Safari/537.36 Edge/14.14393|}' \
  --cookie 'key=${4:value}' \
  --url 'http://example.com'
```

## `http download`

Download from url and save to /path/to/file, using curl [&uarr;](#Commands)

```bash
curl --request ${1|GET,POST|} -sL \
  --user-agent '${2|Shellman,Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML\, like Gecko) Chrome/58.0.3029.110 Safari/537.36,Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:53.0) Gecko/20100101 Firefox/53.0,Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML\, like Gecko) Chrome/51.0.2704.79 Safari/537.36 Edge/14.14393|}' \
  --output '/path/to/file' \
  --url 'http://example.com/file.zip'
```

## `ftp list`

Get the list of files on the ftp server at specific path [&uarr;](#Commands)

```bash
curl ftp://user:password@ipOrDomain/directoryPathOnServer/
```

## `ftp download`

Download specified file from ftp server [&uarr;](#Commands)

```bash
curl ftp://user:password@ipOrDomain/filePathOnServer/
```

## `ftp upload`

Upload specified file to ftp server [&uarr;](#Commands)

```bash
curl -T fileToUpload ftp://user:password@ipOrDomain/directoryPathOnServer/
```

## `ftp delete file`

Delete specified file from ftp server [&uarr;](#Commands)

```bash
curl ftp://user:password@ipOrDomain/path/fileToDelete -Q "DELE fileToDelete"
```

## `ftp rename`

Rename specified file/directory on ftp server [&uarr;](#Commands)

```bash
curl ftp://user:password@ipOrDomain/path/ -Q "-RNFR path/renameFrom" -Q "-RNTO path/renameTo"
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
substring=`echo -e "${var:offset:length}"`
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
if (( int1 == ${2:int2} )); then
  echo equal
fi
```

## `if int !=`

if integers are not equal [&uarr;](#Commands)

```bash
if (( int1 != ${2:int2} )); then
  echo not equal
fi
```

## `if int >`

if integer greater than value [&uarr;](#Commands)

```bash
if (( int > ${2:val} )); then
  echo greater
fi
```

## `if int >=`

if integer greater than or equal value [&uarr;](#Commands)

```bash
if (( int >= ${2:val} )); then
  echo greater equal
fi
```

## `if int <`

if integer lesser than value [&uarr;](#Commands)

```bash
if (( int < ${2:val} )); then
  echo lesser
fi
```

## `if int <=`

if integer lesser than or equal value [&uarr;](#Commands)

```bash
if (( int <= ${2:val} )); then
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
if [ -e "$path" ]; then
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

## `let`

arithmetic operations [&uarr;](#Commands)

```bash
let "result = var1 ${3|+,-,*,/,%|} var2"
```

## `expr`

arithmetic operations [&uarr;](#Commands)

```bash
expr 2 ${2|+,-,\*,/,%|} 3
```

## `math ++`

increment variable [&uarr;](#Commands)

```bash
$((${1|var++,++var|}))
```

## `math --`

decrement variable [&uarr;](#Commands)

```bash
$((${1|var--,--var|}))
```

## `math +`

add two variables (+= assign result to var1) [&uarr;](#Commands)

```bash
$((var1 ${2|+,+=|} var2))
```

## `math -`

subtract var2 from var1 (-= assign result to var1) [&uarr;](#Commands)

```bash
$((var1 ${2|-,-=|} var2))
```

## `math *`

multiply var1 and var2 (*= assign result to var1) [&uarr;](#Commands)

```bash
$((var1 ${2|*,*=|} var2))
```

## `math /`

divide var1 by var2 (/= assign result to var1) [&uarr;](#Commands)

```bash
$((var1 ${2|/,/=|} var2))
```

## `math %`

reminder of dividing var1 by var2 (%= assign result to var1) [&uarr;](#Commands)

```bash
$((var1 ${2|%,%=|} var2))
```

## `math ^`

exponentiate base to power [&uarr;](#Commands)

```bash
$((base ** power))
```

## `math random`

generate random integer x such as min < x < max [&uarr;](#Commands)

```bash
echo $((min + RANDOM % $((max-min))))
```

## `math √`

square root of var up to scale decimal places [&uarr;](#Commands)

```bash
result=`echo "scale=${2|0,1,2,3,4,5,6,7,8,9|};sqrt($var)" | bc`
```

## `math 0.00`

math operations with up to scale decimal places precision [&uarr;](#Commands)

```bash
result=`echo "scale=${2|0,1,2,3,4,5,6,7,8,9|};($var1 ${4|+,-,*,/,^|} $var2)" | bc`
```

## `math const π`

math PI constant [&uarr;](#Commands)

```bash
MATH_PI='3.14159265358979323846264338327950288'
```

## `math const e`

math Napier's constant [&uarr;](#Commands)

```bash
MATH_NAPIER='2.71828182845904523536028747135266249'
```

## `math const 𝛾`

math Euler-Mascheroni constant [&uarr;](#Commands)

```bash
MATH_GAMMA='0.57721566490153286060651209008240243'
```

## `math const ϕ`

math golden ration constant [&uarr;](#Commands)

```bash
MATH_GOLDEN_RATIO='1.61803398874989484820458683436563811'
```

## `math const Ω`

math Omega constant [&uarr;](#Commands)

```bash
MATH_OMEGA='0.56714329040978387299996866221035554'
```

## `file delete,file remove`

delete file(s) [&uarr;](#Commands)

```bash
rm -f ./path/file
```

## `file find,directory find`

find files (-type f) or directories (-type d) by name [&uarr;](#Commands)

```bash
result=`find ./path -maxdepth ${2|0,1,2,3,4,5,6,7,8,9|} -type ${3|f,d|} -name "criteria"`
```

## `file read`

read a file [&uarr;](#Commands)

```bash
cat "$filepath" | while read line; do
  echo "$line"
done
```

## `file search,search in files,find in files`

find files which contain search criteria [&uarr;](#Commands)

```bash
result=`find ./path -maxdepth ${2|0,1,2,3,4,5,6,7,8,9|} -type f -exec grep "criteria" {} +`
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

## `region`

Comment out a special region (i.e. variable declarations [&uarr;](#Commands)

```bash
# >>>>>>>>>>>>>>>>>>>>>>>> name >>>>>>>>>>>>>>>>>>>>>>>>
$0
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

## `system memory info`

System memory information in kilobytes (KB) [&uarr;](#Commands)

```bash
sysMemory${1|MemTotal,MemFree,MemAvailable,Cached,Buffers,Active,Inactive,SwapTotal,SwapFree,SwapCached|}=`cat /proc/meminfo | grep '${1}' | awk '{print $2}' | head -n 1`
```

## `system distro name`

OS ID (i.e. Ubuntu) [&uarr;](#Commands)

```bash
lsb_release -i | awk '{print $3}'
```

## `system distro version`

OS Release (i.e. 16.04) [&uarr;](#Commands)

```bash
lsb_release -r | awk '{print $2}'
```

## `system distro codename`

OS codename (i.e. xenial) [&uarr;](#Commands)

```bash
lsb_release -c | awk '{print $2}'
```

## `system kernel name`

OS kernel name (i.e. Linux) [&uarr;](#Commands)

```bash
uname -s
```

## `system kernel release`

OS kernel release (i.e. 4.4.0-140-generic) [&uarr;](#Commands)

```bash
uname -r
```

## `system processor type`

OS processor type (i.e. x86_64) [&uarr;](#Commands)

```bash
uname -p
```

## `system processor count`

Processor count (cores) [&uarr;](#Commands)

```bash
lscpu | grep 'CPU(s)' |awk '{print $2}' | head -n 1
```

## `system processor architecture`

Processor architecture (i.e. x86_64) [&uarr;](#Commands)

```bash
lscpu | grep 'Architecture' |awk '{print $2}' | head -n 1
```

## `system processor model`

Processor model name (i.e. Intel(R) Core(TM) i5-5200U CPU @ 2.20GHz) [&uarr;](#Commands)

```bash
lscpu | grep 'Model name' |cut -d ' ' -f 3- | sed -e 's/^[[:space:]]*//'
```

## `system uptime`

System uptime (hh:mm:ss) [&uarr;](#Commands)

```bash
sys_uptime=`uptime | cut -d ' ' -f2`
```

## `timeout`

Run command within a time frame [&uarr;](#Commands)

```bash
timeout seconds command
```

## `fn banner simple`

function: print a banner with provided title [&uarr;](#Commands)

```bash
# Usage: banner_simple "my title"
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

print a color banner. [&uarr;](#Commands)

```bash
# Usage: banner_color green "my title"
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

## `fn import`

import functions from other shellscript files [&uarr;](#Commands)

```bash
# Usage: import "mylib"
function import() {
  local file="./lib/$1.sh"
  if [ -f "$file" ]; then
    source "$file"
  else
    echo "Error: Cannot find library at: $file"
    exit 1
  fi
}
```

## `fx import`

call import function, to import functions from other shellscript files located in a directory (default: lib) relative to current script file [&uarr;](#Commands)

```bash
import "libname"
```

## `fn options`

provide a list of options to user and return the index of selected option [&uarr;](#Commands)

```bash
# Usage: options=("one" "two" "three"); chooseOption "Choose:" 1 "${options[@]}"; choice=$?; echo "${options[$choice]}"
function chooseOption() {
  echo "$1"; shift
  echo `tput sitm``tput dim`-"Change selection: [up/down]  Select: [ENTER]" `tput sgr0`
  local selected="$1"; shift
  ESC=`echo -e "\033"`
  cursor_blink_on()  { tput cnorm; }
  cursor_blink_off() { tput civis; }
  cursor_to()        { tput cup $(($1-1)); }
  print_option()     { echo  `tput dim` "   $1" `tput sgr0`; }
  print_selected()   { echo `tput bold` "=> $1" `tput sgr0`; }
  get_cursor_row()   { IFS=';' read -sdR -p $'\E[6n' ROW COL; echo ${ROW#*[}; }
  key_input()        { read -s -n3 key 2>/dev/null >&2; [[ $key = $ESC[A ]] && echo up; [[ $key = $ESC[B ]] && echo down; [[ $key = "" ]] && echo enter; }
  for opt; do echo; done
  local lastrow=`get_cursor_row`
  local startrow=$(($lastrow - $#))
  trap "cursor_blink_on; echo; echo; exit" 2
  cursor_blink_off
  : selected:=0
  while true; do
    local idx=0
    for opt; do
      cursor_to $(($startrow + $idx))
      if [ $idx -eq $selected ]; then
        print_selected "$opt"
      else
        print_option "$opt"
      fi
      ((idx++))
    done
    case `key_input` in
      enter) break;;
      up)    ((selected--)); [ $selected -lt 0 ] && selected=$(($# - 1));;
      down)  ((selected++)); [ $selected -ge $# ] && selected=0;;
    esac
  done
  cursor_to $lastrow
  cursor_blink_on
  echo
  return $selected
}
```

## `fx options`

call options function [&uarr;](#Commands)

```bash
options=("one" "two" "three")
chooseOption "Choose:" 1 "${options[@]}"; choice=$?
echo "${options[$choice]}" selected
```

## `fn progress`

progress bar function [&uarr;](#Commands)

```bash
function progressBar() {
  local BAR='████████████████████'
  local SPACE='                    '
  for i in {1..20}; do
    echo -ne "\r|${BAR:0:$i}${SPACE:$i:20}| $(($i*5))% [ $2 ] "
    sleep $1
  done
  echo -ne ''
}
```

## `fx progress`

call progress bar function [&uarr;](#Commands)

```bash
progressBar ${1|.1,.2,.3,.4,.5,1,2,5|} "Installing foo..."
```

## `fn scan`

Scan host's port range (tcp/udp) [&uarr;](#Commands)

```bash
# Usage: scan proto host fromPort toPort
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

