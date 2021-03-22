# Change Log

## 5.3.0

- `iterate directories` snippet added.
- `archive compress .zip` snippet added.
- `archive decompress .zip` snippet added.
- `system uptime seconds` snippet added. Seconds can be formatted by `fn time format seconds` snippet.
- `remove old/new files/directories` snippet for removing files (-type f) or directories (-type d) if they are older (-mtime n) than n days or newer (+mtime n) than n days.
- Comparison of floating point numbers: `float|double >`, `float|double >=`, `float|double <`, `float|double <=`, `float|double ==`, `float|double !=` added.

## 5.2.0

- `echo array` snippet added.
- `fn/fx math factorial` snippet added.
- `fn/fx math fibonacci` snippet added.
- `fn/fx math fibonacci series` snippet added.
- `am I root` & `am I not root` snippets added.
- `fn/fx input multichoice`: select more than one option.
- `array contains`: check if the array contains an element".
- `fn/fx time format seconds` snippet to format seconds into days/hours/minutes/seconds.

Now there is a `input` namespace for getting user input:

- `input text`: formerly `ask question`.
- `input password`: don't show input on screen.
- `fn/fx input choice`: formerly `fn/fx options`.

## 5.1.0

- Smooth percentage for `fn progress`.
- `echo text` and `echo variable` snippet added.
- assign variable value via `variable assign` | `variable set` snippet.
- `iterate files` snippet for iterating specific file extensions inside a path.
- read/expand variable value via `var` | `variable read` | `variable expand` snippet.

## 5.0.0

- Default command substitution to `$(command)`.
- Default variable expansion to `"${variable}"`.
- Unified <kbd>TAB</kbd> ordering.
- Some functions and snippets enhanced.
- Descriptive aliases added to some snippets.

## 4.10.0

- `fn/fx urlencode`: Encode URL.
- `fn/fx urldecode`: Decode URL.

## 4.9.0

- `string random`: Random string from provided characters with desired length.

## 4.8.0

- `loop infinite`: Infinite loop.

## 4.7.0

- New `Shellman ebook` release. [Download](https://github.com/yousefvand/shellman-ebook) free ebook (pdf, epub, mobi)
- Fixed issues [#20](https://github.com/yousefvand/shellman/issues/20) and [#21](https://github.com/yousefvand/shellman/issues/21).
- directory
  - `directory delete nested` | `directory remove nested`: delete directory and all contents.
- Updated `summary` snippet to include exit codes.
- Added `loop` to `while` and `until` prefix.
- Added samples:
  - Step by step [sample script](samples/backup-tool) tutorial.

## 4.6.0

- Fixed issue [#19](https://github.com/yousefvand/shellman/issues/19)
- Handle events `event ...`
  - `EXIT`: Run commands when script exits.
  - `CTRL+C`: Run commands when script is interrupted.
- Animation
  - Define animation frame.
  - Call `animate` function. Check [sample animations](samples/animation/funny-cat.sh).
  - pacman animation for texts. Check [sample animations](samples/animation/pacman.sh).
- String
  - Concat: Concatenate two strings.

## 4.5.0

- Renice process and sub processes. [PR](https://github.com/yousefvand/shellman/pull/17) by [@jakiro2017](https://github.com/jakiro2017).
- `archive...`
  - Compress/decompress file/directory to/from `.tar.xz`

## 4.4.0

- Fixed issue [#13](https://github.com/yousefvand/shellman/issues/13). Thanks to [@wangyuew](https://github.com/wangyuew).
- `fn/fx [version | semver] compare`
  - Compare two [semver](https://semver.org) versions. Function compares first argument with the second and returns `>`, `<` or `=`.
- `string substring [count | frequency]`: Frequency of a substring in a string. You may need character escaping, for example if string is `1.2.3` and substring is `.` then substring should be `\.` when passed to command.

## 4.3.1

- Fixed issue [#10](https://github.com/yousefvand/shellman/issues/10) in math namespace.
- Some enhancements in math namespace.

## 4.3.0

- Process commands `process ...`
  - `list`: List all processes.
  - `id`: Get process ID by name.
  - `kill`: Kill process by name.
- `cmd renice`: Change running process priority.

## 4.2.0

- Shellman ebook available for [download](https://github.com/yousefvand/shellman-ebook).
- `math fn/fx`
  - `sum`: Calculate sum of given integers.
  - `product`: Calculate product of given integers.
  - `average`: Calculate average of given integers.

## 4.1.1

- Minor bugs fixed.
- `let` added for arithmetic operations.

## 4.0.0

- fixed issues [#3](https://github.com/yousefvand/shellman/issues/3), [#4](https://github.com/yousefvand/shellman/issues/4), [#5](https://github.com/yousefvand/shellman/issues/5), [#6](https://github.com/yousefvand/shellman/issues/6), [#7](https://github.com/yousefvand/shellman/issues/7)

## 3.6.1

- Escape dollar sign

## 3.6.0

- `system...`
  - System information (CPU, Memory, Kernel...)

## 3.5.0

- `git...`
  - `git` commands.

## 3.4.0

- `ftp...`
  - `list`: Get the list of files on the ftp server at specific path.
  - `download`: Download specified file from ftp server.
  - `upload`: Upload specified file to ftp server.
  - `delete`: Delete specified file from ftp server.
  - `rename`: Rename specified file/directory on ftp server.

## 3.3.0

- `http...`
  - `GET`: Send http GET request using curl.
  - `POST`: Send data with http POST, using curl.
  - `POST file`: Send file via http POST, using curl.
  - `PUT`: Send data with http PUT, using curl
  - `DELETE`: Send http DELETE request using curl.
  - `cookie`: Send http request with cookies, using curl.
  - `header`: Send http request with protocol/custom header, using curl.
  - `download`: Download from url and save to `/path/to/file`, using curl.

## 3.2.0

- `fn/fx progress`
  - Dummy progress bar with custom message

## 3.1.0

- `nice`
  - Execute a command with desired privilege
- `archive...`
  - Compress/decompress file/directory (currently .tar.gz supported. More soon)
- `crypto...`
  - base64 encode / decode
  - Calculate hash of string using md5, sha, sha1, sha224, sha256, sha384 and sha512 algorithms

## 3.0.0

- `array...`
  - Array namespace

## 2.3.1

- `date...`, `time...`
  - Date and time commands
- `file find`
  - Find files and directories by name (wildcard supported)
- `file search`
  - Search inside all files for a text (i.e. word)

## 2.2.1

- Argument parsing
  - parse command line arguments (flags/switches)
- bug fixed (int less than)

## 2.1.0

- Documentation improved.
- `import` function (fn/fx) added.
  - Import functions from other shell script files.
- `options` function (fn/fx) added.
  - Provide a list of options to user.

## 2.0.0

- `math...`
  - math constants (π, e, Ω, ...)
  - math operations (with precision).

## 1.5.1

- `string substring` typo fixed.

## 1.5.0

- `region`: A region for specific purpose (functions, variables...)
- string manipulation
  - `string length`: length of string in characters.
  - `string trim`: remove leading and trailing white space(s).
  - `string trim left`: remove leading white space(s).
  - `string trim right`: remove trailing white space(s).
  - `string trim all`: remove all white space(s).
  - `string replace`: find all occurrences of a substrings and replace them.
  - `string reverse`: reverse string characters.
  - `string toLower`: convert string to lowercase.
  - `string toUpper`: convert string to uppercase.
  - `string substring`: part of the string from offset with [length] characters.
  - `string contains`: check whether string contains substring.
  - `string indexOf`: first index of substring in string.

## 1.4.0

- `assign array`: assign elements to an array.
- `ip ...`
  - `ips`: array of local IPs.
  - `ip info`: public ip information (country, city...).
  - `ip public`: public ip address.
- `random number`: generate random integer x such as min < x < max.
- `service manage`: manage service operations via `systemctl`.
- `timeout`: run command within a time frame. Cancel command if not finished within x seconds.
- `fn/fx`
  - `scan`: scan host's port range and find open ports (tcp/udp).

## 1.3.0

- stopwatch
- `fn...` / `fx...` whole function and usage insertion.
  - `fn banner simple`: insert function to print a banner.
  - `fn banner color`: insert function to print a color banner.

## 1.2.0

- Check last command success/failure via `cmd...`
- Directory operations via `directory...`

## 1.1.2

- Fixed vscode marketplace header color contrast

## 1.1.1

- Minor bugs fixed
- Write colorful text
- Write formatted text (bold, italic, dim, reverse)
- Snippet documentation added

## 0.0.1

Initial release
