#!/usr/bin/env python3
import argparse
import subprocess
import sys
import os
import tarfile
from subprocess import Popen, PIPE

def parseArguments(raw_args):
  parser = argparse.ArgumentParser(prog="build",
                               description='Project Builder')
  parser.add_argument('-project','-p',required=True,help="Github project", type=str)
  parser.add_argument('-binname','-n',required=True,help="Binary name project", type=str)
  
  parser.add_argument('-arch','-a',required=True, help='Architecture to build', type=str)
                               
  parser.add_argument('--branch', '-b', help='Git branch to build', default="master", type=str)
  
  return parser.parse_args(raw_args)

 
def build(project, branch, arch, bin_name):
  if not os.path.exists(os.path.join("/", "out", "project")) :
    os.makedirs(os.path.join("/", "out","project"), exist_ok=True)
    process = None
    process = subprocess.Popen(["git", "clone",  "-b", branch,
                                "https://github.com/{}".format(project), 
                                os.path.join("/","out","project")],
                                shell=False, stdout=PIPE)  
    stdout,err = process.communicate()
    if err != None:
      os.exit("Error while cloning project: {}".format(err))

  else:
    print("Project already exists, branch name is ignored")
  
  go_path = "{}:/out/project".format(os.environ["GOPATH"])
  process = subprocess.Popen(["go", "get", "-d", "./..."],
                   cwd=os.path.join("/", "out", "project"), shell=False,
                   env=dict(os.environ, GOARCH=arch, GOPATH=go_path), stdout=PIPE)
  stdout,err = process.communicate()
  if err != None:
    os.exit("Error while getting dependencies: {}".format(err))
  
  process = subprocess.Popen(["go", "build", '-o', os.path.join('bin', bin_name), bin_name],
                   cwd=os.path.join("/", "out", "project"), shell=False,
                   env=dict(os.environ, GOARCH=arch, GOPATH=go_path), stdout=PIPE)
  stdout,err = process.communicate()
  if err != None:
    os.exit("Error while building project: {}".format(err))
  
  with tarfile.open(os.path.join("/", "out","project", "{}.{}.{}.tar.gz".format(bin_name, arch, branch)), "w:gz") as tar:
    tar.add(os.path.join("/", "out", "project", "bin", bin_name), arcname=bin_name)


if __name__ == "__main__":
    parsed_args = parseArguments(sys.argv[1:])
    build(parsed_args.project,parsed_args.branch,parsed_args.arch,parsed_args.binname)

    
