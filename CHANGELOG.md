# Change Log

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