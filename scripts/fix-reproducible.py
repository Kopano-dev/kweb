#!/usr/bin/python
##
# Copyright 2019 Kopano and its licensors
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.
#

from __future__ import print_function


import os
import subprocess
import sys

version = "20190906-1"

config = {
    "go": os.environ.get("GO", "go"),
}


# Main below.


def main():
    args = sys.argv[1:]
    if len(args) != 1:
        usage()
        sys.exit(1)

    fn = args[0]
    go = config.get("go")

    buildid = subprocess.check_output([go, 'tool', 'buildid', fn]).strip()
    actionid = b'/'.join(buildid.split(b'/', 2)[:2])
    with open(fn, 'r+b') as f:
        data = f.read()
        idx = data.find(actionid)
        if idx == -1:
            raise ValueError('actionid not found in file')
        f.seek(idx)
        f.write(b'0'*(len(actionid)-2))
        f.write(b'/0')


def usage():
    print("usage: %s <file-to-fix>" % sys.argv[0])

if __name__ == "__main__":
    main()
