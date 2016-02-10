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
	err = utils.InitLoggers(config.Verbose, config.Quiet, config.Syslog)
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
