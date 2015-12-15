import os
import sys
from waflib.Task import Task
from waflib import Logs
import subprocess
from distutils.version import LooseVersion, StrictVersion
import imp
import re
import stat
import distutils.spawn

top = "."
out = "build"
src = "src"


def options(ctx):
    ctx.add_option(
        '--prefix',
        action='store',
     default="/usr/local",
     help='install prefix')
    ctx.add_option(
        '--templates',
        action='store',
     default="*",
     help='templates to install')

def checkVersion(name,cmd,regex,requiredVersion,returnCode):
    res = False
    Logs.pprint('WHITE','{0: <40}'.format('Checking {0} version'.format(name)),sep=': ')
    try:
        p = subprocess.Popen(cmd,stdout=subprocess.PIPE,stderr=subprocess.PIPE)
        out, err = p.communicate()
        outlines = out + err
        res =  { "returncode": p.returncode, "out":''.join(outlines.splitlines()).decode("utf-8") }
        if res["returncode"] != returnCode:
            Logs.pprint('RED','{0} is not available. Cannot continue.'.format(name))
        else:
          v = re.match(regex, res["out"])
          
          version = v.group(1)
          requiredVersionObj = LooseVersion(requiredVersion)
          currentVersionObj = LooseVersion(version)
          if currentVersionObj >= requiredVersionObj:
              Logs.pprint('GREEN',version)
              res = True
          else:
              Logs.pprint('RED','Incorrect version {0} (must be equal or greater than {1}). Cannot continue.'.format(currentVersionObj,requiredVersionObj))
    except Exception as e:
        Logs.pprint('RED','Unable to get version ({0}).'.format(e))
    return res
   
def checkPythonModule(*moduleNames):
    res = False
    for mod in moduleNames:
        Logs.pprint(
            'WHITE',
            '{0: <40}'.format('Checking python module {0}'.format(mod)),
            sep=': ')
        try:
            imp.find_module(mod)
            res = True
            break
        except ImportError:
            Logs.pprint('RED', 'no')
            res = False
    if not res:
        Logs.pprint(
            'RED',
            'Required Python module is not available. Cannot continue')
    else:
        Logs.pprint('GREEN', "ok")
        res = True
    return res


def configure(ctx):
    ctx.env.PREFIX = ctx.options.prefix

    if checkVersion("Python", ["python3", "--version"], "Python ([0-9\.]*)", "3.0.0", 0) == False:
        ctx.fatal("Missing requirements. Installation will not continue.")
    if checkPythonModule("argcomplete") == False:
        ctx.fatal("Missing requirements. Installation will not continue.")


def build(ctx):
    srcPath = ctx.path.find_dir('src/')

    ctx.install_files(os.path.join(ctx.env.PREFIX,'bin'),
      srcPath.ant_glob('**/cfdnsupdater'),cwd=srcPath,relative_trick=True)
    ctx.install_files(os.path.join(ctx.env.PREFIX,'etc/default'),
      srcPath.ant_glob('**/cfdnsupdater.conf.sample'),cwd=srcPath,relative_trick=True)

