# Change Log

## 3.6.0

- `system...`
  - System information (CPU, Memory, Kernel...)

## 3.5.0

- `git...`
  - `git` commands.

## 3.4.0

- `ftp`
  - `list`: Get the list of files on the ftp server at specific path.
  - `download`: Download specified file from ftp server.
  - `upload`: Upload specified file to ftp server.
  - `delete`: Delete specified file from ftp server.
  - `rename`: Rename specified file/directory on ftp server.

## 3.3.0

- `http`
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

- Initial release