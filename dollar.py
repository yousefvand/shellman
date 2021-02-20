import sys
import subprocess as sp
from collections import defaultdict


def _getoutput(cmd):
    """to simulate commands.getoutput in order to work with python 2.6 up to 3.x"""
    out = sp.Popen(cmd, shell=True, stderr=sp.STDOUT, stdout=sp.PIPE).communicate()[0][:-1]
    try:
        o = str(out.decode(errors="ignore"))  # to force str in python 3
    except AttributeError:
        o = out
    return o


def inplace_change(filename, old_line, new_line, commit=False):
    # Safely read the input filename using 'with'
    with open(filename) as f:
        s = f.read()
        if old_line not in s:
            print(f'"{old_line}" not found in {filename}.')
            return

    # Safely write the changed content, if found in the file
    print(f'{filename}: Changing "{old_line}" to "{new_line}"')
    if commit:
        with open(filename, "w") as f:
            s = s.replace(old_line, new_line)
            f.write(s)


def fix_line(line, escaped=False):
    ns = line.count("`")
    dollar = r"$("
    if escaped:
        dollar = r"\\$("
    if ns % 2:
        print(f"ERROR in line {line}")
    for n in range(0, ns, 2):
        line = line.replace("`", dollar, 1)
    for n in range(1, ns, 2):
        line = line.replace("`", ")")
    return line


if __name__ == "__main__":
    """
    Run it by:
    python3 dollar.py # will show what will change but nothing is really changed
    python3 dollar.py commit > changes_applied.txt # Argument 'commit' will update the files
    """
    commit = False
    if "commit" in sys.argv:
        commit = True

    cmd = r"""grep -vE '```|^#' -r --exclude=\*.{js,gif,png,py,md,txt} --exclude-dir=".git"  . --exclude-dir="snippets" \
    | grep -v '`-_' | grep '`'"""

    out_grep = _getoutput(cmd).splitlines()

    fdict = defaultdict(list)

    for row in out_grep:
        rfile, line = row.split(":", maxsplit=1)
        if line.strip().startswith('"#'):
            continue
        flag = False
        if ".json" == rfile[-5:]:
            flag = True
        nline = fix_line(line, escaped=flag)
        fdict[rfile].append(nline)
        inplace_change(rfile, line, nline, commit=commit)
