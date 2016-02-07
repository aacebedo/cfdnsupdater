package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"math"
	"os"
)

const (
	MIN_UPDATE_PERIOD = 10
	MAX_UPDATE_PERIOD = 60 * 60
)

type CommandLine struct {}

func (self *CommandLine) ParseParameters(rawParams []string) (res *CFDNSUpdaterConfiguration, err error) {
	app := kingpin.New("cfdnsupdater", "Automtatic DNS updater for cloudflare.")
	verbose := app.Flag("verbose", "Verbose mode.").Short('v').Bool()
	quiet := app.Flag("quiet", "Quiet mode.").Short('q').Bool()
	syslog := app.Flag("syslog", "Output logs to syslog.").Short('s').Bool()

	inlineconfig := app.Command("inlineconfig", "Pass configuration as inline parameters.")
	email := inlineconfig.Flag("email", "Email used to login to cloudflare account.").Short('e').Required().String()
	apiKey := inlineconfig.Flag("apikey", "Api key used for authentication.").Short('a').Required().String()
	domain := inlineconfig.Flag("domain", "Domain to update.").Short('d').Required().String()
	period := inlineconfig.Flag("period", fmt.Sprintf("Public IP Refresh period in seconds (must be >= %v and <= %v).", MIN_UPDATE_PERIOD, MAX_UPDATE_PERIOD)).Short('p').Int()
	recordtypes := inlineconfig.Flag("recordtype", "Types of records to update.").Short('t').Strings()
	recordnames := inlineconfig.Flag("recordname", "Names of records to update.").Short('n').Strings()

	fileconfig := app.Command("fileconfig", "Pass configuration as file.")
	fileconfig_path := fileconfig.Arg("config", "Configuration file.").Required().String()
  
	res = &CFDNSUpdaterConfiguration{}
	switch kingpin.MustParse(app.Parse(rawParams)) {
	case inlineconfig.FullCommand():
		config := DomainConfiguration{Email: *email, ApiKey: *apiKey, Domain: *domain, Period: *period}
		if len(*recordtypes) > 0 {
			for _, v := range *recordtypes {
				var recordType, parseErr = FromString(v)
				if err != nil {
					err = errors.New(fmt.Sprintf("Invalid configuration: '%s'", parseErr))
					return
				}
				config.RecordTypes = append(config.RecordTypes, recordType)
			}
		}

		if len(*recordnames) == 0 {
			config.RecordNames = nil
		}

		res.DomainConfigs = []DomainConfiguration{config}
		res.Verbose = *verbose
		res.Quiet = *quiet
		res.Syslog = *syslog

	case fileconfig.FullCommand():

		var file, _ = os.Open(*fileconfig_path)
		var decoder = json.NewDecoder(file)
		var decodeErr = decoder.Decode(&res)
		if decodeErr != nil {
			err = errors.New(fmt.Sprintf("Invalid configuration file '%s': %s", fileconfig_path, decodeErr))
			return
		}
	}

	for _, domain := range res.DomainConfigs {
		domain.Period = int(math.Min(MAX_UPDATE_PERIOD, float64(domain.Period)))
		domain.Period = int(math.Max(MIN_UPDATE_PERIOD, float64(domain.Period)))
	}
  
	return

}
