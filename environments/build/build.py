#!/usr/bin/env python3
import argparse
import subprocess
import sys
import os
from subprocess import Popen, PIPE

parser = argparse.ArgumentParser(prog="build",
                             description='Builder.')

parser.add_argument('arch', help='Architecture to build', type=str)
                             
parser.add_argument('--branch', '-b', help='Git branch to build', type=str)

args = parser.parse_args(sys.argv[1:])

if not os.path.exists(os.path.join("/", "out","cfdnsupdater")) :
  os.makedirs(os.path.join("/", "out"), exist_ok=True)
  process = None
  if args.branch == None:
    process = subprocess.Popen(["git", "clone", "https://github.com/aacebedo/cfdnsupdater"],
        cwd=os.path.join("/", "out"), shell=False, stdout=PIPE)
  else:
    process = subprocess.Popen(["git", "clone", "https://github.com/aacebedo/cfdnsupdater",
        "-b", args.branch], cwd=os.path.join("/", "out"),
        shell=False, stdout=PIPE)
    
  stdout = process.communicate()[0]
  print("{}".format(stdout))

process = subprocess.Popen(["go", "install", "./..."],
                 cwd=os.path.join("/", "out", "cfdnsupdater"), shell=False,
                 env=dict(os.environ,GOARCH=args.arch), stdout=PIPE)
stdout = process.communicate()[0]
print("{}".format(stdout))
