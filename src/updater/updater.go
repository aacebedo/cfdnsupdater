//This file is part of CFDNSUpdater.
//
//    CFDNSUpdater is free software: you can redistribute it and/or modify
//    it under the terms of the GNU General Public License as published by
//    the Free Software Foundation, either version 3 of the License, or
//    (at your option) any later version.
//
//    CFDNSUpdater is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    GNU General Public License for more details.
//
//    You should have received a copy of the GNU General Public License
//    along with CFDNSUpdater.  If not, see <http://www.gnu.org/licenses/>.

package updater

import (
	"github.com/aacebedo/cfdnsupdater/src/core"
	"github.com/aacebedo/cfdnsupdater/src/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/levigross/grequests"
	"net"
	"reflect"
	"sync"
	"time"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("cfdnsupdater.updater")

type DomainUpdater struct {
	domain      string
	recordTypes core.RecordTypeSlice
	recordNames []string
	authParams  map[string]string
	period      int
	restSession *grequests.Session
}

func NewDomainUpdater(domain string, email string, apiKey string, recordTypes core.RecordTypeSlice, recordNames []string, period int) (res *DomainUpdater) {
	res = &DomainUpdater{}
	res.restSession = grequests.NewSession(nil)
	res.authParams = map[string]string{
		"x-auth-email": email,
		"x-auth-key":   apiKey,
	}
	res.domain = domain
	res.recordTypes = recordTypes
	res.recordNames = recordNames
	res.period = period

	return
}

func (self *DomainUpdater) GetPublicIP() (res net.IP, err error) {

	url := "https://api.ipify.org"
	resp, err := self.restSession.Get(url, nil)
	if err == nil && resp.Ok {
		res = net.ParseIP(resp.String())
	} else {
		err = errors.New("")
	}
	return
}

func (self *DomainUpdater) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	logger.Infof("Starting to monitor domain '%s' every %v seconds.",
		self.domain, self.period)
	url := "https://api.cloudflare.com/client/v4/zones"
	var previousIP net.IP
SLEEP_PERIOD:
	for {
		publicIP, err := self.GetPublicIP()
		if err == nil {
			if !reflect.DeepEqual(previousIP, publicIP) {
				previousIP = publicIP
				logger.Infof("Public IPAddress is %v", publicIP)
				resp, err := self.restSession.Get(url, &grequests.RequestOptions{
					Headers: self.authParams,
					Params:  map[string]string{"name": self.domain},
				})

				if err != nil {
					logger.Errorf("%s: %s", self.domain, err.Error())
					continue SLEEP_PERIOD
				}
				if !resp.Ok {
					logger.Errorf("%s: Unable to get domain's zones: %v", self.domain, 
					  resp)
					continue SLEEP_PERIOD
				}

				var zoneDetails core.ZoneRequestResult
				err = json.Unmarshal([]byte(resp.String()), &zoneDetails)

				if err != nil {
					logger.Errorf("%s: %s", self.domain, err.Error())
					continue SLEEP_PERIOD
				}

				if zoneDetails.Result_info.Total_count != 1 {
					logger.Errorf("%s: Domain not found, Check your connection or credentials", self.domain)
					continue SLEEP_PERIOD
				}

				logger.Debugf("Domain '%s' found. Processing it...", self.domain)
				zoneDetail := zoneDetails.Result[0]
				logger.Debugf("Getting DNS records for the domain '%s'", zoneDetail.Name)
				url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zoneDetail.Id)
				resp, err = self.restSession.Get(url, &grequests.RequestOptions{Headers: self.authParams})

				if err != nil {
					logger.Errorf("%s: %s", self.domain, err.Error())
					continue SLEEP_PERIOD
				}

				if !resp.Ok {
					logger.Errorf("%s: Unable to get DNS records, %s.", self.domain,resp )
					continue SLEEP_PERIOD
				}

				var recordDetails core.RecordRequestResult
				err = json.Unmarshal([]byte(resp.String()), &recordDetails)
				if err != nil {
					logger.Errorf("%s: %s", self.domain, err.Error())
					return
				}
			RECORD_PROCESSING:
				for _, recordDetail := range recordDetails.Result {
					logger.Debugf("%s: Processing record %s:'%s'.",
						zoneDetail.Name, recordDetail.Type, recordDetail.Name)
					recordIp := net.ParseIP(recordDetail.Content)
					recordType, convertErr := core.FromString(recordDetail.Type)
							
					if !reflect.DeepEqual(recordIp, publicIP) && convertErr == nil &&
						(len(self.recordTypes) == 0 || self.recordTypes.Contains(recordType)) &&
						(len(self.recordNames) == 0 || utils.StringInSlice(recordDetail.Name,
							self.recordNames)) {
						logger.Debugf("%s: Record %s:'%s' needs to be updated",
							self.domain, recordDetail.Type, recordDetail.Name)

						recordDetail.Content = publicIP.String()
						url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s",
							zoneDetail.Id, recordDetail.Id)
						resp, err = self.restSession.Put(url, &grequests.RequestOptions{
							Headers: self.authParams,
							JSON:    recordDetail})
						if err != nil {
							logger.Warningf("%s: Record ignored (%s)", self.domain, err.Error())
							continue RECORD_PROCESSING
						}

						if !resp.Ok {
							logger.Warningf("%s: Unable to update DNS record '%s:%s'", self.domain, recordDetail.Type, recordDetail.Name)
						} else {
							logger.Infof("%s: Record %s:%s successfully updated with value %s",
								zoneDetail.Name,
								recordDetail.Type,
								recordDetail.Name,
								recordDetail.Content)
						}
					} else {
						logger.Debugf("%s: Record %s:'%s' ignored", zoneDetail.Name, recordDetail.Type, recordDetail.Name)
					}
				}
			} else {
				logger.Debugf("Public IP Address didn't change. Records will not be changed")
			}
		} else {
		  logger.Warningf("Unable to get Public IP Address. Check your connection.")
		}
		
		time.Sleep(time.Duration(self.period) * time.Second)
	}
}
