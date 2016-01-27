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
This is the main module of the cfdnsupdater
"""

import sys
from cfdnsupdater.core import CloudFlareDNSUpdater, \
CloudFlareDNSUpdaterThread, RecordType
import json
import logging
from cfdnsupdater.cmdline import CloudFlareDNSUpdaterCmdLineParser
from cfdnsupdater.loggers import ROOTLOGGER

class CloudFlareDNSUpdaterMain(object):
    """
    Main class
    """
    @staticmethod
    def process_inline_config(parsed_args):
        """
        Extract configuration from inline parameters
        """
        updaters = []
        email = parsed_args.email
        apikey = parsed_args.apikey
        record_types_to_update = set()
        record_names_to_update = set()
        domain = parsed_args.domain
        period = parsed_args.period
        if parsed_args.recordtype is not None:
            for recordtype in parsed_args.recordtype:
                record_types_to_update.add(
                    RecordType.from_string(recordtype))

        if parsed_args.recordname is not None:
            for name in parsed_args.recordname:
                record_names_to_update.add(
                    "{}.{}".format(name, parsed_args.domain))


        updaters.append(CloudFlareDNSUpdaterThread(
            CloudFlareDNSUpdater(email, apikey, domain,
                                 {"record_types": record_types_to_update,
                                  "record_names":record_names_to_update,
                                 }), period))
        return updaters

    @staticmethod
    def process_file_config(parsed_args): # pylint: disable=too-many-branches
        """
        Extract configuration from a config file
        """
        updaters = []
        config_file = None
        try:
            config_file = open(parsed_args.config, "r", encoding="utf8")
        except FileNotFoundError as excpt:
            ROOTLOGGER.error(
                "Unable to load configuration file '{}'\
 (file not found)".format(parsed_args.config))
            raise excpt

        config_file_json = json.loads(config_file.read())
        if not "email" in config_file_json:
            raise ValueError("Invalid configuration: Missing email")
        if not "apikey" in config_file_json:
            raise ValueError("Invalid configuration: Missing apikey")
        if not "period" in config_file_json:
            raise ValueError("Invalid configuration: Missing period")

        if not "instances" in config_file_json or \
           not isinstance(config_file_json["instances"], list):
            raise ValueError("Invalid configuration: Missing instances")

        email = config_file_json["email"]
        apikey = config_file_json["apikey"]
        period = config_file_json["period"]

        for instance in config_file_json["instances"]:
            try:
                if not "types" in instance or \
                   not isinstance(instance["types"], list):
                    raise ValueError("Invalid configuration:\
 Missing record types in instance")
                if not "names" in instance or \
                   not isinstance(instance["names"], list):
                    raise ValueError("Invalid configuration:\
 Missing record names  in instance")
                if not "domain" in instance:
                    raise ValueError("Invalid configuration:\
 Missing domain")
                domain = instance["domain"]
                record_types_to_update = set()
                record_names_to_update = set()
                for recordtype in instance["types"]:
                    record_types_to_update.add(
                        RecordType.from_string(recordtype))
                for name in instance["names"]:
                    record_names_to_update.add(
                        "{}.{}".format(name, domain))
                updaters.append(
                    CloudFlareDNSUpdaterThread(
                        CloudFlareDNSUpdater(
                            email,
                            apikey,
                            domain,
                            {"record_types":record_types_to_update,
                             "record_names":record_names_to_update
                            }
                            ),
                        period))
            except ValueError as raised_except:
                ROOTLOGGER.warning(raised_except)
        return updaters

    @staticmethod
    def init_loggers(parsed_args):
        """
        Initialize loggers of the program
        """
        formatter = logging.Formatter('%(message)s')
        hdlr = None
        if parsed_args.syslog:
            hdlr = logging.handlers.SysLogHandler(address='/dev/log',
                                                  facility=\
                                                  logging.handlers.\
                                                  SysLogHandler.LOG_DAEMON)
        else:
            hdlr = logging.StreamHandler(sys.stdout)

        hdlr.setLevel(1)
        hdlr.setFormatter(formatter)
        ROOTLOGGER.addHandler(hdlr)

    @staticmethod
    def main():
        """
        Main function
        """
        res = 0
        try:
            parsed_args = CloudFlareDNSUpdaterCmdLineParser.\
                          parse_args(sys.argv[1:])
            CloudFlareDNSUpdaterMain.init_loggers(parsed_args)
            ROOTLOGGER.setLevel(logging.INFO)
            if parsed_args.verbose:
                ROOTLOGGER.setLevel(logging.DEBUG)
            elif parsed_args.veryverbose:
                ROOTLOGGER.setLevel(4)
            elif parsed_args.quiet:
                ROOTLOGGER.setLevel(0)

            updaters = []
            if parsed_args.configtype == "inlineconfig":
                updaters = CloudFlareDNSUpdaterMain.\
                           process_inline_config(parsed_args)
            else:
                updaters = CloudFlareDNSUpdaterMain.\
                           process_file_config(parsed_args)

            for updater in updaters:
                updater.start()
            for updater in updaters:
                updater.join()
        except KeyboardInterrupt:
            for updater in updaters:
                updater.stop()
            for updater in updaters:
                updater.join()
        except (ValueError, FileNotFoundError):
            res = 1
        except SystemExit:
            pass
        return res

if __name__ == "__main__":
    CloudFlareDNSUpdaterMain.main()
