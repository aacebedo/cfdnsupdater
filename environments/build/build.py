#!/usr/bin/env python3
#This file is part of CFDNSUpdater.
#
#    CFDNSUpdater is free software: you can redistribute it and/or modify
#    it under the terms of the GNU General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.
#
#    CFDNSUpdater is distributed in the hope that it will be useful,
#    but WITHOUT ANY WARRANTY; without even the implied warranty of
#    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
#    GNU General Public License for more details.
#
#    You should have received a copy of the GNU General Public License
#    along with CFDNSUpdater.  If not, see <http://www.gnu.org/licenses/>.

import argparse
import subprocess
import sys
import os
import tarfile
from subprocess import Popen, PIPE

def parseArguments(raw_args):
  parser = argparse.ArgumentParser(prog="build",
                               description='Project Builder')
  parser.add_argument('-project','-p',required=True,help="Github project",
                      type=str)
  parser.add_argument('-binname','-n',required=True,help="Binary name project",
                      type=str)
  
  parser.add_argument('-arch','-a',required=True, help='Architecture to build',
                      type=str)
                               
  parser.add_argument('--branch', '-b', help='Git branch to build',
                      default="master", type=str)
  
  return parser.parse_args(raw_args)

 
def build(project, branch, arch, bin_name):
  os.makedirs(os.path.join("/", "out"), exist_ok=True)
  if len(os.listtdir("/", "out")) == 0 :
    process = None
    process = subprocess.Popen(["git", "clone",  "-b", branch,
                                "https://github.com/{}".format(project), 
                                os.path.join("/","go","src","github.com",project)],
                                shell=False, stdout=PIPE)  
    stdout,err = process.communicate()
    if err != None:
      os.exit("Error while cloning project: {}".format(err))

  else:
    print("Project already exists, branch name is ignored")
  
  
  process = subprocess.Popen(["go", "get", "-d", "./..."],
                   cwd=os.path.join("/", "go","src","github.com",project), shell=False,
                   env=dict(os.environ, GOARCH=arch),
                            stdout=PIPE)
  stdout,err = process.communicate()
  if err != None:
    os.exit("Error while getting dependencies: {}".format(err))
  
  process = subprocess.Popen(["go", "build", '-o', os.path.join('bin', bin_name), bin_name],
                   cwd=os.path.join("/", "out", "project"), shell=False,
                   env=dict(os.environ, GOARCH=arch, GOPATH=go_path), stdout=PIPE)
  stdout,err = process.communicate()
  if err != None:
    os.exit("Error while building project: {}".format(err))
  
   #with tarfile.open(os.path.join("/", "out","project", "{}.{}.{}.tar.gz".format(bin_name, arch, branch)), "w:gz") as tar:
   # tar.add(os.path.join("/", "out", "project", "bin", bin_name), arcname=bin_name)
  else :
     os.exit("/out directory is not empty")

if __name__ == "__main__":
    parsed_args = parseArguments(sys.argv[1:])
    build(parsed_args.project,parsed_args.branch,parsed_args.arch,parsed_args.binname)

    
