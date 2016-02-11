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

package configuration

import (
	"errors"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math"
	"cfdnsupdater/utils"
)

const (
	MIN_UPDATE_PERIOD = 10
	MAX_UPDATE_PERIOD = 60 * 60
	VERSION           = "0.8"
)

type CommandLine struct{}

func (self *CommandLine) ParseParameters(rawParams []string) (res *CFDNSUpdaterConfiguration, err error) {
	app := kingpin.New("cfdnsupdater", "Automtatic DNS updater for cloudflare.")
	app.Version(VERSION)
	app.HelpFlag.Short('h')
	fileconfig_path := app.Flag("config", "Configuration file.").Short('c').Required().String()
	verbose := app.Flag("verbose", "Verbose mode.").Short('v').Bool()
	quiet := app.Flag("quiet", "Quiet mode.").Short('q').Bool()
	syslog := app.Flag("syslog", "Output logs to syslog.").Bool()
	filelog := app.Flag("filelog", "Output logs to file (/var/log/cfdnsupdater).").Bool()

	kingpin.MustParse(app.Parse(rawParams))
	
	res = &CFDNSUpdaterConfiguration{}
	res.Verbose = *verbose
	res.Quiet = *quiet
	if *syslog {
	  res.LoggingMode = utils.Syslog
	} else if *filelog {
	  res.LoggingMode = utils.Filelog
	} else {
	  res.LoggingMode = utils.Consolelog
	}

	fileContent, rerr := ioutil.ReadFile(*fileconfig_path)
	if rerr != nil {
		err = errors.New(fmt.Sprintf("Unable to read configuration file '%s'", *fileconfig_path))
		return
	}

	err = yaml.Unmarshal([]byte(fileContent), &res.DomainConfigs)
	if err != nil {
		err = errors.New(fmt.Sprintf("Unable to parse configuration file '%s'", *fileconfig_path))
		return
	}
 
	for _, domain := range res.DomainConfigs {
		domain.Period = int(math.Min(MAX_UPDATE_PERIOD, float64(domain.Period)))
		domain.Period = int(math.Max(MIN_UPDATE_PERIOD, float64(domain.Period)))
	}
	return
}
