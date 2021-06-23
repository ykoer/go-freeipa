#!/usr/bin/env python3

# This will download the errors.py file from GitHub.
# TODO: use argparse to get command line parameters
#   This should include the github url for errors.py, as well as the output file (errors.json)


import urllib.request
import re
import inspect
import json
import sys

ERRORS_PY_URL = "https://raw.githubusercontent.com/freeipa/freeipa/master/ipalib/errors.py"

import_regex = re.compile(r"^(from [\w\.]+ )?import \w+( as \w+)?$")

# Get the output file from the commandline.
if len(sys.argv) > 1:
    output_file = sys.argv[1]
else:
     output_file = '../data/errors.json'

def should_keep(l):
    return (import_regex.match(l) is None)


errors_py_str = urllib.request.urlopen(ERRORS_PY_URL).read().decode('utf-8')
errors_py_str = "\n".join(
    [l for l in errors_py_str.splitlines() if should_keep(l)])
errors_py_str = """
class Six:
    PY3 = True
six = Six()
ungettext = None
class Messages:
    def iter_messages(*args):
        return []
messages = Messages()

""" + errors_py_str


import types
errors_mod = types.ModuleType("errors")

exec(errors_py_str, errors_mod.__dict__)

error_codes = [
    {
        "name": k,
        "errno": v.errno
    } for k, v in inspect.getmembers(errors_mod)
    if hasattr(v, '__dict__') and type(v.__dict__.get('errno', None)) == int
]
error_codes.sort(key=lambda x: x["errno"])

with open(output_file, 'w') as f:
    json.dump(error_codes, f)