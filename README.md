# shellman

[![GitHub release](https://img.shields.io/github/release/yousefvand/shellman.svg?style=plastic)](https://github.com/yousefvand/shellman/releases)
[![GitHub license](https://img.shields.io/github/license/yousefvand/shellman.svg?style=plastic)](https://github.com/yousefvand/shellman/blob/master/LICENSE.md)
[![GitHub stars](https://img.shields.io/github/stars/yousefvand/shellman.svg?style=plastic)](https://github.com/yousefvand/shellman/stargazers)
[![GitHub issues](https://img.shields.io/github/issues/yousefvand/shellman.svg?style=plastic)](https://github.com/yousefvand/shellman/issues)
[![Gitter](https://img.shields.io/gitter/room/badges/shields.svg?style=plastic)](https://gitter.im/vscode-shellman/Lobby)
![Author](https://img.shields.io/badge/Author-Remisa-ff69b4.svg?style=plastic)

Shell script snippet

Learn Bash Scripting with Shellman, examples included. [Download](https://github.com/yousefvand/shellman-ebook) free ebook (pdf, epub, mobi).

Read [Shellman story on medium](https://medium.com/@remisa.yousefvand/shellman-reborn-f2cc948ce3fc) (3 min read).

![shellman](images/demo.gif)

## Math example

![shellman](images/math.gif)

## `fn... / fx...` example

![shellman](images/banner.gif)

## Requirements

- vscode
- bashdb (If you need to debug your scripts)

## Usage

Install extension in vscode by:

```bash
ext install Remisa.shellman
```

Start typing and Shellman will provide you available commands.

For more convenience similar commands are grouped into same prefixes. Here is an overview:

`bash`

Shebang should be used as the first line of your script. You can replace `bash` with any other installed scripting language like `node` or `python`.

`argument parsing` | `parse args`

Parse command-line arguments

`cmd...`

Run external commands and check if operation succeeded.

`color...`

Write colorful

`directory...`

Directory operations

`func...`

Snippets related to function.

`for...`

Iterate different collections/arrays...

`file...`

File operations

`format...`

Write in bold, italic, dim, reverse format.

`ftp...` and `http...`

Web methods and functionalities: GET, POST...

`git...`

git commands

`if...`

Wide range of logical conditions which are more common in bash scripts.

`math...`

Math operations

`string...`

String utilities

`stopwatch...`

Start and stop, stopwatch and read elapsed time.

## `fn` / `fx`

`fn...`

inserts a whole function into script. Function declaration should proceed its usage.

`fx...`

Call function which is declared by `fn...`

## Function usage examples

- banner simple
  - print a banner with provided title.
  - example: `banner_simple "my title"`
- banner color
  - print a color banner.
  - example: `banner_color red "my title"`
- import
  - Organize your project and reuse functions. Import functions from other shell script files. Default import directory is `lib`. This directory should be where the calling script exists and contain library files with `.sh` extension. For example if `libname.sh` contains some useful functions and exists in `lib` directory, you can import those functions into your script and call them.
  - example: `import "somefile"` will import all defined functions in `somefile.sh` from `lib` directory where calling script resides.
- animation
  - Create some frames with same size using `animation frame` snippet.
  - Insert `animate` function using `fn animation animate` snippet.
  - Call `animate` function.
  - Check [sample animations](samples/animation).

## List of [commands](COMMANDS.md)

## [Full release Notes](CHANGELOG.md)

## Latest release Notes

### 4.6.0

- Fixed issue [#19](https://github.com/yousefvand/shellman/issues/19)
- Handle events `event ...`
  - `exit`: Run commands when script exits.
  - `CTRL+C`: Run commands when script is interrupted.
- Animation
  - Define animation frame.
  - Call `animate` function. Check [sample animations](samples/animation/funny-cat.sh).
  - pacman animation for texts. Check [sample animations](samples/animation/pacman.sh).
- String
  - Concat: Concatenate two strings.
  