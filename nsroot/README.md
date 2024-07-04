# Namespaces

`nsroot` is the root folder of namespaces and each namespace has its own folder inside it.

Only one depth level is allowed (to maintain readability).

Add new folder for a new namespace.

Put snippet inside appropriate folder and run go build tool from root of the project:

```bash
./build
```

or in windows:

```bash
/>build.exe

to build go build tool from source:

```bash
# Compile executable
go build build.go

# Run without build
go run build.go
```

or download the appropriate one based on your operating system from project latest release page. This tool doesn't accept arguments and must be executed from project root folder.

Bump version in `package.json` file according to [SemVer](https://semver.org).

## Caution

- Directory structure depth is fixed to one.
- Don't use dot (`.`) in file and directory names.

## Testing changes

First disable `Shellman` extension. Your compiled version if snippets is in file `snippets/snippets`. Copy all of contents of this file. In vscode from main menu go to:

```
File -> Preferences -> User Snippets
```

Select `shellscript.json` from drop down menu. Delete all file contents and paste `snippets/snippets` content there and save the file.

Now you can use `test.sh` at project root to test your changes.

If you are using `vscode` there is a task named `build-test` that will do all above tasks.