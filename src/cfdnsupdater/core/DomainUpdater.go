package core

import (
  "cfdnsupdater/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/levigross/grequests"
	"net"
	"reflect"
	"sync"
	"time"
)

type DomainUpdater struct {
	domain      string
	recordTypes RecordTypeSlice
	recordNames []string
	authParams  map[string]string
	period      int
	restSession *grequests.Session
}

func GetPublicIP() (res net.IP, err error) {
	restSession := grequests.NewSession(nil)
	url := "https://api.ipify.org"
	resp, err := restSession.Get(url, nil)
	if err == nil && resp.Ok {
		res = net.ParseIP(resp.String())
	} else {
		err = errors.New("")
	}
	return
}

func NewDomainUpdater(domain string, email string, apiKey string, recordTypes RecordTypeSlice, recordNames []string, period int) (res *DomainUpdater) {
	res = &DomainUpdater{}
	res.restSession = grequests.NewSession(nil)
	res.authParams = map[string]string{
		"x-auth-email": email,
		"x-auth-key":   apiKey,
	}
	res.domain     = domain
	res.recordTypes = recordTypes
	res.recordNames = recordNames
	res.period = period
	
	return
}

func (self *DomainUpdater) Run(wg *sync.WaitGroup) {
	defer wg.Done()
	logger.Infof("Starting to monitor domain '%s' every %v seconds.",
		self.domain, self.period)
	url := "https://api.cloudflare.com/client/v4/zones"
	for {
		publicIP, _ := GetPublicIP()
		resp, err := self.restSession.Get(url, &grequests.RequestOptions{
			Headers: self.authParams,
			Params:  map[string]string{"name": self.domain},
		})

		if err != nil {
			logger.Errorf("%s: %s", self.domain, err.Error())
			return
		}
		if !resp.Ok {
			logger.Errorf("%s: Unable to get domain's zones", self.domain)
			return
		}

		var zoneDetails ZoneRequestResult
		err = json.Unmarshal([]byte(resp.String()), &zoneDetails)

		if err != nil {
			logger.Errorf("%s: %s", self.domain, err.Error())
			return
		}

		if zoneDetails.Result_info.Total_count != 1 {
			logger.Errorf("%s: Domain not found", self.domain)
			return
		}

		logger.Debugf("Domain '%s' found. Processing it...", self.domain)
		zoneDetail := zoneDetails.Result[0]
		logger.Debugf("Getting DNS records for the domain '%s'", zoneDetail.Name)
		url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zoneDetail.Id)
		resp, err = self.restSession.Get(url, &grequests.RequestOptions{Headers: self.authParams})

		if err != nil {
			logger.Errorf("%s: %s", self.domain, err.Error())
			return
		}

		if !resp.Ok {
			logger.Errorf("%s: Unable to get DNS records", self.domain)
			return
		}

		var recordDetails RecordRequestResult
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
			recordType, convertErr := FromString(recordDetail.Type)
			if !reflect.DeepEqual(recordIp, publicIP) && convertErr == nil &&
				(len(self.recordTypes) == 0 || (self.recordTypes.Contains(recordType))) &&
				(len(self.recordNames) == 0 || utils.StringInSlice(recordDetail.Name, 
				                                                self.recordNames)) {
				logger.Debugf("%s: Record %s:'%s' needs to be updated",
					self.domain, recordDetail.Type, recordDetail.Name)

				recordDetail.Content = publicIP.String()
				url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s",
					zoneDetail.Id, recordDetail.Id)
				resp, err = self.restSession.Put(url,&grequests.RequestOptions{
						                                          Headers: self.authParams,
						                                          JSON:    recordDetail,})
				if err != nil {
					logger.Warningf("%s: Record ignored (%s)", self.domain, err.Error())
					continue RECORD_PROCESSING
				}

				if !resp.Ok {
					logger.Warningf("%s: Unable to update DNS record '%s:%s'", self.domain,recordDetail.Type, recordDetail.Name)
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
		time.Sleep(time.Duration(self.period) * time.Second)
	}
}
