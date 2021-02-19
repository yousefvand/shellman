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


def inplace_change(filename, old_line, new_line):
    # Safely read the input filename using 'with'
    with open(filename) as f:
        s = f.read()
        if old_line not in s:
            print(f'"{old_line}" not found in {filename}.')
            return

    # Safely write the changed content, if found in the file
    with open(filename, "w") as f:
        print(f'Changing "{old_line}" to "{new_line}" in {filename}')
        s = s.replace(old_line, new_line)
        f.write(s)


def fix_line(line):
    ns = line.count("`")
    if ns % 2:
        print(f"ERROR in line {line}")
    for n in range(0, ns, 2):
        line = line.replace("`", "$(", 1)
    for n in range(1, ns, 2):
        line = line.replace("`", ")")
    return line


cmd = r"""grep -vE '```|^#' -r --exclude=\*.{js,gif,png,py} --exclude-dir=".git" --exclude=CHANGELOG.md \
--exclude=README.md . | grep -v '`-_' | grep '`'"""

out_grep = _getoutput(cmd).splitlines()

fdict = defaultdict(list)

for row in out_grep:
    rfile, line = row.split(":", maxsplit=1)
    if line.strip().startswith('"#'):
        continue
    nline = fix_line(line)
    fdict[rfile].append(nline)
    inplace_change(rfile, line, nline)
