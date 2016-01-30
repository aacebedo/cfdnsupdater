package core

import (
	"encoding/json"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"math"
	"os"
)

const (
	MIN_UPDATE_PERIOD = 10
	MAX_UPDATE_PERIOD = 60 * 60
)

type DomainConfig struct {
	Email       string           `json:"email"`
	ApiKey      string           `json:"apikey"`
	Period      int              `json:"period"`
	Domain      string           `json:"domain"`
	RecordNames *[]string        `json:"record_names"`
	RecordTypes *RecordTypeSlice `json:"record_types,RecordType"`
}

type CFDNSUpdaterConfiguration struct {
	Verbose       bool           `json:"verbose"`
	Quiet         bool           `json:"quiet"`
	Syslog        bool           `json:"syslog"`
	DomainConfigs []DomainConfig `json:"domain_configs"`
}

func ParseParameters(rawParams []string) *CFDNSUpdaterConfiguration {
	var app = kingpin.New("cfdnsupdater", "Automtatic DNS updater for cloudflare.")
	var verbose = app.Flag("verbose", "Verbose mode.").Short('v').Bool()
	var quiet = app.Flag("quiet", "Quiet mode.").Short('q').Bool()
	var syslog = app.Flag("syslog", "Output logs to syslog.").Short('s').Bool()

	var inlineconfig = app.Command("inlineconfig", "Pass configuration as inline parameters.")
	var email = inlineconfig.Flag("email", "Email used to login to cloudflare account.").Short('e').Required().String()
	var apiKey = inlineconfig.Flag("apikey", "Api key used for authentication.").Short('a').Required().String()
	var domain = inlineconfig.Flag("domain", "Domain to update.").Short('d').Required().String()
	var period = inlineconfig.Flag("period", fmt.Sprintf("Public IP Refresh period in seconds (must be >= %v and <= %v).", MIN_UPDATE_PERIOD, MAX_UPDATE_PERIOD)).Short('p').Int()
	var recordtypes = inlineconfig.Flag("recordtype", "Types of records to update.").Short('t').Strings()
	var recordnames = inlineconfig.Flag("recordname", "Names of records to update.").Short('n').Strings()

	var fileconfig = app.Command("fileconfig", "Pass configuration as file.")
	var fileconfig_path = fileconfig.Arg("config", "Configuration file.").Required().String()

	var result CFDNSUpdaterConfiguration
	switch kingpin.MustParse(app.Parse(rawParams)) {
	case inlineconfig.FullCommand():

		var config = DomainConfig{Email: *email, ApiKey: *apiKey, Domain: *domain, Period: *period, RecordNames: recordnames}
		if len(*recordtypes) > 0 {
			config.RecordTypes = &RecordTypeSlice{}
			for _, v := range *recordtypes {
				var recordType, err = FromString(v)
				if err != nil {
					RootLogger.Fatalf("Invalid configuration: '%s'", err)
				}
				*config.RecordTypes = append(*config.RecordTypes,recordType)
			}
		}
		result.DomainConfigs = []DomainConfig{config}
		result.Verbose = *verbose
		result.Quiet = *quiet
		result.Syslog = *syslog

	case fileconfig.FullCommand():

		var file, _ = os.Open(*fileconfig_path)
		var decoder = json.NewDecoder(file)
		var err = decoder.Decode(&result)
		if err != nil {
			RootLogger.Fatalf("Invalid configuration file '%s': %s", fileconfig_path, err)
		}

	}

	for _, domain := range result.DomainConfigs {
		domain.Period = int(math.Min(MAX_UPDATE_PERIOD, float64(domain.Period)))
		domain.Period = int(math.Max(MIN_UPDATE_PERIOD, float64(domain.Period)))
	}

	return &result

}
