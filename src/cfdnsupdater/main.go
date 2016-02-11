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

package main

import (
	"cfdnsupdater/utils"
	"cfdnsupdater/configuration"
	"cfdnsupdater/updater"
	"os"
	"sync"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("cfdnsupdater")

func main() {
  var cmdLine configuration.CommandLine
	config,err := cmdLine.ParseParameters(os.Args[1:])
	if(err != nil) {
	  logger.Fatalf(err.Error())
	}
	err = utils.InitLoggers(config.Verbose, config.Quiet, config.LoggingMode)
	if(err != nil) {
	  logger.Fatalf(err.Error())
	}
	var wg sync.WaitGroup
	wg.Add(len(config.DomainConfigs))
	for domain, domainConfig := range config.DomainConfigs {
		updater := updater.NewDomainUpdater(domain,
                                			domainConfig.Email,
                                			domainConfig.ApiKey,
                                			domainConfig.RecordTypes,
                                			domainConfig.RecordNames,
                                			domainConfig.Period)
		go updater.Run(&wg)
	}
	wg.Wait()

}
