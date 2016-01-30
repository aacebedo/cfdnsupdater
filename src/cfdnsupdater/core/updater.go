package core

import (
	"encoding/json"
	"fmt"
	"github.com/levigross/grequests"
	"net"
	"reflect"
	"sync"
	"time"
)

func GetPublicIP() net.IP {
	var restSession = grequests.NewSession(nil)
	var url = "https://api.ipify.org"
	var resp, err = restSession.Get(url, nil)
	if err == nil && resp.Ok {
		return net.ParseIP(resp.String())
	} else {
		return nil
	}
}

func updateDomain(wg *sync.WaitGroup, domain string,
	recordTypes *RecordTypeSlice, recordNames *[]string,
	authParams map[string]string, period int) {
	defer wg.Done()
	var restSession = grequests.NewSession(nil)
	var p = &grequests.RequestOptions{
		Headers: authParams,
		Params:  map[string]string{"name": domain},
	}
	RootLogger.Infof("Starting to monitor domain '%s' every %v seconds.",
		domain, period)
	var url = "https://api.cloudflare.com/client/v4/zones"

	for {
		var publicIP = GetPublicIP()
		var resp, err = restSession.Get(url, p)

		if err == nil && resp.Ok {
			var zoneDetails ZoneRequestResult
			var _ = json.Unmarshal([]byte(resp.String()), &zoneDetails)

			if zoneDetails.Result_info.Total_count == 1 {
				RootLogger.Debugf("Domain '%s' found. Processing it...", domain)

				for _, zoneDetail := range zoneDetails.Result {
					RootLogger.Debugf("Getting DNS records for the domain '%s'",
						zoneDetail.Name)
					var url = fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zoneDetail.Id)
					var resp, err = restSession.Get(url, &grequests.RequestOptions{Headers: authParams})

					if err == nil && resp.Ok {
						var recordDetails RecordRequestResult
						var _ = json.Unmarshal([]byte(resp.String()), &recordDetails)

						for _, recordDetail := range recordDetails.Result {
							RootLogger.Debugf("%s: Processing record %s:'%s'.",
								zoneDetail.Name, recordDetail.Type, recordDetail.Name)
							var recordIp = net.ParseIP(recordDetail.Content)
              var recordType, convertErr = FromString(recordDetail.Type)
							if !reflect.DeepEqual(recordIp, publicIP) && convertErr == nil &&
								(recordTypes == nil || ((*recordTypes).Contains(recordType))) &&
								(recordNames == nil || stringInSlice(recordDetail.Name, *recordNames)) {

								RootLogger.Debugf("%s: Record %s:'%s' needs to be updated",
									zoneDetail.Name,
									recordDetail.Type,
									recordDetail.Name)

								recordDetail.Content = publicIP.String()
								var url = fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s",
									zoneDetail.Id, recordDetail.Id)
								var resp, err = restSession.Put(url,
									&grequests.RequestOptions{
										Headers: authParams,
										JSON:    recordDetail,
									})

								if err == nil && resp.Ok {
									RootLogger.Infof("%s: Record %s:%s successfully updated with value %s",
										zoneDetail.Name,
										recordDetail.Type,
										recordDetail.Name,
										recordDetail.Content)
								} else {
									var errMsg string

									if err != nil {
										errMsg = err.Error()
									} else {
										errMsg = resp.String()
									}

									RootLogger.Errorf("%s: Unable to update Record %s:'%s: %s",
										zoneDetail.Name, recordDetail.Type, recordDetail.Name, errMsg)
								}
							} else {
								RootLogger.Infof("%s: Record %s:'%s': ignored",
									 zoneDetail.Name, recordDetail.Type, recordDetail.Name)
							}
						}
					} else {
						var errMsg string
						if err != nil {
							errMsg = err.Error()
						} else {
							errMsg = resp.String()
						}
						RootLogger.Errorf("%s: Unable to get records: %s",
							zoneDetail.Name, errMsg)
					}

				}
			} else {
				RootLogger.Errorf("%s: No records found", domain)
			}
		} else {
			var errMsg string
			if err != nil {
				errMsg = err.Error()
			} else {
				errMsg = resp.String()
			}
			RootLogger.Error("Unable to get zones for domain '%s': %s",
				domain, errMsg)
		}
		time.Sleep(time.Duration(period) * time.Second)
	}
}

func SpawnDomainUpdaters(configs []DomainConfig) {
	var wg sync.WaitGroup
	wg.Add(3)
	for _, config := range configs {
		var headers = map[string]string{
			"x-auth-email": config.Email,
			"x-auth-key":   config.ApiKey,
		}
		go updateDomain(&wg, 
		  config.Domain, 
		  config.RecordTypes,
		  config.RecordNames, 
		  headers, 
		  config.Period)
	}
	wg.Wait()
}
