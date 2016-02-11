#!/usr/bin/env python3
# This file is part of CFDNSUpdater.
#
#    CFDNSUpdater is free software: you can redistribute it and/or modify
#    it under the terms of the GNU General Public License as published by
#    the Free Software Foundation, either version 3 of the License, or
#    (at your option) any later version.
#
#    CFDNSUpdater is distributed in the hope that it will be useful,
#    but WITHbuild ANY WARRANTY; withbuild even the implied warranty of
#    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
#    GNU General Public License for more details.
#
#    You should have received a copy of the GNU General Public License
#    along with CFDNSUpdater.  If not, see <http://www.gnu.org/licenses/>.

import json
import argparse
import subprocess
import sys
import os
import uuid
import time
import shutil
from subprocess import Popen, PIPE

def parseArguments(raw_args):
  parser = argparse.ArgumentParser(prog="build",
                               description='Project Builder')
  parser.add_argument('-project', '-p', required=True, help="Github project",
                      type=str)
  parser.add_argument('-arch', '-a', required=True, help='Architecture to build',
                      type=str)
  parser.add_argument('--branch', '-b', help='Git branch to build',
                      default="master", type=str)
  parser.add_argument('-binname', '-bn', required=True,
                      help='binname',
                      type=str)
  parser.add_argument('-distrib', '-d', required=True,
                      help='Distribution',
                      type=str)
  parser.add_argument('-user', '-u', required=True,
                      help='User',
                      type=str)
  parser.add_argument('-description', '-dc', required=True,
                      help='User',
                      type=str)
  parser.add_argument('--outputdirpath', '-o', help='Output direcotry',
                      required=True, type=str)
  parser.add_argument('--licenses', '-li', help='Software licences',
                      default=[], type=str, action='append')
  parser.add_argument('--labels', '-la', help='Package labels', action='append',
                      default=[], type=str)
  return parser.parse_args(raw_args)

def generateTmpDir():
  tmp_dir_path = None
  for x in range(0, 5):
    tmp_dir_path = os.path.join(os.path.abspath(os.sep), "tmp", str(uuid.uuid4()))
    if not os.path.exists(tmp_dir_path) :
      os.makedirs(tmp_dir_path, exist_ok=True)
      break
    else:
      tmp_dir_path = None
  if tmp_dir_path == None:
    raise Exception("Unable to generate a tmp direcctory")
  return tmp_dir_path
    
def generatePackage(build_dir_path,
                    package_type, package_name, version):
  process = subprocess.Popen(["fpm", "-t", package_type,
                                 "-n", package_name,
                                 "-p", build_dir_path,
                                 "-f",
                                 "-v", version.replace("/", "_"),
                                  "-C", os.path.join(build_dir_path,"packaging"),
                                 "-s", "dir", "."], shell=False, stdout=PIPE)
  stdout, err = process.communicate()
  if err != None:
    os.exit("Error while building project: {}".format(err))
      
def build(build_dir_path, project, branch, arch, bin_name):
  if len(os.listdir(build_dir_path)) != 0:
    raise Exception("Build error: {} is not empty.".format(build_dir_path))
  go_dir_path = os.path.join(generateTmpDir(), "go")
  src_dir_path = os.path.join(go_dir_path, 'src', "github.com", project)
  
  process = None
  process = subprocess.Popen(["git", "clone", "-b", branch,
                                "https://github.com/{}".format(project),
                                src_dir_path], shell=False, stdout=PIPE)
  stdout, err = process.communicate()
  if err != None:
    os.exit("Error while cloning project: {}".format(err))

  process = subprocess.Popen(["go", "get", "-d", "./..."],
                   cwd=src_dir_path, shell=False,
                   env=dict(os.environ,
                            GOARCH=arch,
                            GOPATH=go_dir_path,
                            CGO_ENABLED="0"),
                   stdout=PIPE)
  stdout, err = process.communicate()
  if err != None:
    os.exit("Error while getting dependencies: {}".format(err))
  
  process = subprocess.Popen(["go", "install", "./..."],
                   cwd=src_dir_path, shell=False,
                   env=dict(os.environ,
                            GOARCH=arch,
                            GOPATH=go_dir_path,
                            CGO_ENABLED="0",
                            GOBIN=os.path.join(build_dir_path,
                                               "packaging", "usr",
                                               "local", "bin")),
                             stdout=PIPE)
  stdout, err = process.communicate()
  if err != None:
    os.exit("Error while building project: {}".format(err))
  
  if os.path.exists(os.path.join(src_dir_path, "resources")) :
    for name in os.listdir(os.path.join(src_dir_path, "resources")):
      shutil.copytree(os.path.join(src_dir_path, "resources", name),
                      os.path.join(build_dir_path, "packaging", name))
  


def generateBintrayDescriptor(build_dir,
                              bin_name,
                              user,
                              desc,
                              project,
                              version,
                              distrib,
                              arch,
                              licenses=[],
                              labels=[]):
  github_addr = "https://github.com/{}".format(project)
  descriptor = {"package":{
                           "name":bin_name,
                           "repo":bin_name,
                           "subject":user,
                           "desc":desc,
                           "website_url":github_addr,
                           "issue_tracker_url":github_addr,
                           "vcs_url":github_addr,
                           "github_use_tag_release_notes":True,
                           "licenses":licenses,
                           "labels":labels,
                           "public_download_numebrs":False,
                           "public_stats":False
                           },
                "version":{
                           "name":version,
                           "desc":desc,
                           "released":time.strftime("%Y-%m-%d"),
                           "vcs_tag":version,
                           "gpgSign":False
                           },
                "files":[
                         {
                          "includePattern": os.path.join(build_dir, "(.*\.deb)"),
                          "uploadPattern": "$1",
                          "matrixParams":
                            {
                              "deb_distribution":distrib,
                              "deb_component":"main",
                              "deb_architecture":arch
                             }
                         }
                         ],
                "publish":True
                }
  file = open(os.path.join(build_dir, "bintray.desc"), 'w')
  json.dump(descriptor, file, ensure_ascii=False, indent=2)
  file.close()

if __name__ == "__main__":
    args = parseArguments(sys.argv[1:])
    
    if not os.path.exists(args.outputdirpath):
      os.makedirs(args.outputdirpath, exist_ok=True)
      
    build(args.outputdirpath,
          args.project, args.branch,
          args.arch, args.binname)
    generatePackage(args.outputdirpath, "deb", args.binname, args.branch)
    generatePackage(args.outputdirpath, "tar", args.binname, args.branch)

    generateBintrayDescriptor(args.outputdirpath,
                              args.binname,
                              args.user,
                              args.description,
                              args.project,
                              args.branch,
                              args.distrib,
                              args.arch,
                              args.licenses,
                              args.labels)
                                           

    
