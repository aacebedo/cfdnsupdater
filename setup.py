#!/usr/bin/env python3
#
# cfdnsupdater
# Copyright (c) 2015, Alexandre ACEBEDO, All rights reserved.
#
# This library is free software; you can redistribute it and/or
# modify it under the terms of the GNU Lesser General Public
# License as published by the Free Software Foundation; either
# version 3.0 of the License, or (at your option) any later version.
#
# This library is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
# Lesser General Public License for more details.
#
# You should have received a copy of the GNU Lesser General Public
# License along with this library.
#
"""
Setup script for CFDNUpdater
"""
import sys
from setuptools import setup, find_packages
import argparse
import platform
import os


def process_setup():
    """
    Setup function
    """
    if sys.version_info < (3,0):
        sys.exit("CloudFlare DNS updater only supports python3. Please run setup.py with python3.")
  
    parser = argparse.ArgumentParser(add_help=False)
    parser.add_argument("--prefix", type=str)
    args, _ = parser.parse_known_args(sys.argv[1:])

    prefix = sys.exec_prefix
    if args.prefix != None:
        prefix = args.prefix

    data_files = []
    if platform.system() == "Linux":
        for dname, _, filelist in os.walk("resources/etc"):
            for fname in filelist:
                data_files.append((os.path.join(prefix, "{}".format(dname)),\
                                   ["{}/{}".format(dname, fname)]))

    setup(
        name="cfdnsupdater",
        version="0.5",
        packages=find_packages("src"),
        package_dir ={'':'src'},
        data_files=data_files,
        install_requires=['argcomplete>=1.0.0','argparse>=1.0.0', 'requests>=2.9.1'],
        author="Alexandre ACEBEDO",
        author_email="Alexandre ACEBEDO",
        description="DNS updater for CloudFlare",
        license="LGPLv3",
        keywords="cloudflare dns",
        url="http://github.com/aacebedo/cfdnsupdater",
        entry_points={'console_scripts':
                      ['cfdnsupdater = cfdnsupdater.__main__:\
                      CloudFlareDNSUpdaterMain.main']}
    )
    
if __name__ == "__main__":
    process_setup()
    
