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
This module contains cmdline objects of the cfdnsupdater
"""

from cfdnsupdater.core import CloudFlareDNSUpdaterThread, RecordType
import argcomplete
import argparse

class CloudFlareDNSUpdaterCmdLineParser:
    """
    This classes builds the arguments parser and parses the given
    arguments
    """
    @staticmethod
    def add_common_args(parser):
        """
        Add common arguments to the command line
        """
        verbosity_group = parser.add_mutually_exclusive_group()

        verbosity_group.add_argument(
            '--verbose',
            '-v',
            help='Verbose mode',
            action='store_true',
            default=False)
        verbosity_group.add_argument(
            '--veryverbose',
            '-vv',
            help='Very verbose mode',
            action='store_true',
            default=False)
        verbosity_group.add_argument(
            '--quiet',
            '-q',
            help='Quiet mode',
            action='store_true',
            default=False)

        parser.add_argument('--syslog',
                            '-s',
                            help='Logging is done with syslog,\
 console otherwise',
                            action='store_true',
                            default=False)

    @staticmethod
    def parse_args(raw_args):
        """
        Function to parse the command line arguments
        """
        # Create main parser
                        # Create main parser
        parser = argparse.ArgumentParser(
            prog="CloudFlareDNSUpdater",
            description='Automtatic DNS updater for cloudflare.')
        subparsers = parser.add_subparsers(dest="configtype")
        inline_config_parser = subparsers.add_parser(
            "inlineconfig",
            help="process configuration as inline parameters")

        inline_config_parser.add_argument(
            '-e',
            '--email',
            required=True,
            help='email used to log to cloudflare',
            type=str)
        inline_config_parser.add_argument(
            '-a',
            '--apikey',
            required=True,
            help='Api key used for authentication',
            type=str)
        inline_config_parser.add_argument(
            '-p',
            '--period',
            default=CloudFlareDNSUpdaterThread.MINIMUM_PERIOD,
            help='Public IP Refresh period (must be >= 30)',
            type=int)

        inline_config_parser.add_argument(
            'domain',
            help='Domain to update',
            type=str)
        inline_config_parser.add_argument(
            '--recordtype',
            '-t',
            help='Types of records to update',
            action='append',
            type=str,
            choices=[str(s) for s in RecordType]
        )
        inline_config_parser.add_argument(
            '--recordname',
            '-n',
            help='Names of records to update',
            action='append',
            type=str)

        CloudFlareDNSUpdaterCmdLineParser.add_common_args(inline_config_parser)

        file_config_parser = subparsers.add_parser(
            "fileconfig",
            help="process configuration in a file")
        file_config_parser.add_argument(
            'config',
            help='Config file',
            type=str)

        CloudFlareDNSUpdaterCmdLineParser.add_common_args(file_config_parser)

        # Parse the command
        argcomplete.autocomplete(parser)
        parsed_args = parser.parse_args(raw_args)

        if parsed_args.configtype == None:
            parser.print_usage()
            print("CFDNSUpdater: error: too few arguments")
            raise ValueError()
        return parsed_args
