# Namespaces

`nsroot` is the root folder of namespaces and each namespace has its own folder inside it.

Only one depth level is allowed (due to readability).

Add new folder for a new namespace.

Put snippet inside appropriate folder and run npm build task:

```bash
$ npm run build
```

Bump version in `package.json` file according to [SemVer](https://semver.org).

## Caution:

- Directory structure depth is fixed to one.
- Don't use dot (`.`) in file and directory names.
