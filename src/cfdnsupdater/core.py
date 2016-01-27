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
CFDNSUpdater core module
"""
from enum import Enum
import json
import requests
import threading
import time
from cfdnsupdater.loggers import ROOTLOGGER

class RecordType(Enum):
    """
    Enumeration describing supported record types
    """
    A = 1 # pylint: disable=invalid-name
    AAAA = 2
    CNAME = 3
    MX = 4 # pylint: disable=invalid-name
    LOC = 5
    SRV = 6
    SPF = 7
    TXT = 8
    NS = 9 # pylint: disable=invalid-name

    def __str__(self):
        return str(self.name)

    def __repr__(self):
        return str(self.name)

    @staticmethod
    def from_string(val):
        """
        Turns a string into a record type
        """
        mapping_table = {str(RecordType.A): RecordType.A,
                         str(RecordType.AAAA): RecordType.AAAA,
                         str(RecordType.CNAME): RecordType.CNAME,
                         str(RecordType.MX): RecordType.MX,
                         str(RecordType.LOC): RecordType.LOC,
                         str(RecordType.SRV): RecordType.SRV,
                         str(RecordType.SPF): RecordType.SPF,
                         str(RecordType.TXT): RecordType.TXT,
                         str(RecordType.NS): RecordType.NS
                        }
        if val in mapping_table.keys():
            return mapping_table[val]
        else:
            raise RuntimeError(
                "Unable to convert string {} into a ReturnType".format(val))


class CloudFlareDNSUpdater(object):
    """
    Core class of the updater
    """
    def __init__(self, email, apikey, domain, instance):
        self.email = email
        self.apikey = apikey
        self.record_types_to_update = instance["record_types"]
        self.record_names_to_update = instance["record_names"]
        self.domain = domain

    @staticmethod
    def get_public_ip():
        """
        Get the public IP of your host
        It makes use of the ipify service
        """
        res = requests.get('https://api.ipify.org').text
        if res is None:
            ROOTLOGGER.error(
                "Unable to retrieve public IP, updater cannot start")
        return res

    def update_domain(self, new_record_content):
        """
        Update a domain
        """
        ROOTLOGGER.debug(
            "Getting cloudflare record for zone '{}'".format(self.domain))
        # First search for a corresponding domain
        response = requests.get(
            "https://api.cloudflare.com/client/v4/zones",
            params={"name": self.domain},
            headers={'x-auth-email': self.email,
                     'x-auth-key': self.apikey})
        response_json = json.loads(response.content.decode())
        ROOTLOGGER.log(5, "Received JSON response: '{}'".format(response_json))
        if response_json["result_info"]["total_count"] == 1:
            # If domain has been found, procesing it
            ROOTLOGGER.debug("One domain found, processing it")
            # Retrieving zone records for zone
            zone = response_json["result"][0]
            zone_id = zone["id"]
            ROOTLOGGER.debug(
                "Getting DNS records for the domain '{}'".format(self.domain))
            response = requests.get(
                "https://api.cloudflare.com/client/v4/zones/{}/dns_records"\
                .format(
                    zone_id),
                headers={'x-auth-email': self.email,
                         'x-auth-key': self.apikey})
            response_json = json.loads(response.content.decode())
            ROOTLOGGER.log(
                5,
                "Received JSON response: '{}'".format(response_json))
            dns_records = response_json["result"]
            for record in dns_records:
                # Each records of the zone will is process
                try:
                    ROOTLOGGER.debug(
                        "Processing record '{}'".format(record["name"]))
                    record_type = RecordType.from_string(record["type"])
                    record_id = record["id"]
                    record_name = record["name"]
                    record_content = record["content"]
                    # Record is updated only if its type is in the
                    # record types to be updated or its name is
                    # in the name of type to update
                    if (len(self.record_types_to_update) == 0 or \
                        record_type in self.record_types_to_update) and \
                        (len(self.record_names_to_update) == 0 or \
                         record_name in self.record_names_to_update) and \
                            new_record_content != record_content:
                        ROOTLOGGER.debug(
                            "Record {} needs to be updated"\
                            .format(record["name"]))
                        old_record_content = record["content"]
                        record["content"] = new_record_content
                        response = requests.put(
                            "https://api.cloudflare.com/client/v4/zones/{}/\
dns_records/{}".format(zone_id, record_id),
                            data=json.dumps(record),
                            headers={'x-auth-email': self.email,
                                     'x-auth-key': self.apikey})
                        ROOTLOGGER.debug("Executing update")
                        response_json = json.loads(response.content.decode())
                        ROOTLOGGER.log(5,
                                       "Received JSON response: '{}'"\
                                       .format(response_json))
                        if response_json["success"] == True:
                            ROOTLOGGER.info("Content of DNS record ({}) '{}'\
 of domain '{}' has been updated from {} to '{}'".format(record["type"],
                                                         record_name,
                                                         self.domain,
                                                         old_record_content,
                                                         new_record_content))
                        else:
                            ROOTLOGGER.warning("Unable to update content of\
 DNS record ({}) '{}' of domain '{}' from {} \
 to '{}' ({})".format(record["type"],
                      record_name,
                      self.domain,
                      old_record_content,
                      new_record_content,
                      response_json["errors"][0]["error_chain"][0]["message"]))
                    else:
                        ROOTLOGGER.info("Record '{}' ignored"\
                                      .format(record_name))
                except RuntimeError as raised_except:
                    ROOTLOGGER.warning("Record '{}' ignored ({})"\
                        .format(record["name"], raised_except))
        else:
            ROOTLOGGER.warning("Domain '{}' ignored (not found')")

class CloudFlareDNSUpdaterThread(threading.Thread):
    """
    Thread class handling periodic verification of public IP
    and executing the domain update
    """
    MINIMUM_PERIOD = 30

    def __init__(self, updater, period):
        threading.Thread.__init__(self)
        self.running = False
        self.updater = updater
        self.period = CloudFlareDNSUpdaterThread.MINIMUM_PERIOD
        if period < CloudFlareDNSUpdaterThread.MINIMUM_PERIOD:
            ROOTLOGGER.warning(
                "Given period {} seconds is too low, defaulting on {} \
                seconds".format(period,
                                CloudFlareDNSUpdaterThread.MINIMUM_PERIOD))
            self.period = CloudFlareDNSUpdaterThread.MINIMUM_PERIOD
        else:
            self.period = period
        ROOTLOGGER.debug(
            "Getting public IP every {} seconds".format(self.period))

    def run(self):
        previous_public_ip = None
        self.running = True
        while self.running:
            current_public_ip = CloudFlareDNSUpdater.get_public_ip()
            if previous_public_ip != current_public_ip:
                self.updater.update_domain(current_public_ip)
                previous_public_ip = current_public_ip
            time.sleep(self.period)
        ROOTLOGGER.info("CFDNSUpdater stopped")

    def stop(self):
        """
        Stops the thread
        """
        ROOTLOGGER.info("CFDNSUpdater is stopping")
        self.running = False
