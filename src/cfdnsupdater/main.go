package main

import (
	"cfdnsupdater/core"
	"cfdnsupdater/utils"
	"os"
	"sync"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("cfdnsupdater")

func main() {
  var cmdLine core.CommandLine
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
	for _, domainConfig := range config.DomainConfigs {
		updater := core.NewDomainUpdater(domainConfig.Domain,
			domainConfig.Email,
			domainConfig.ApiKey,
			domainConfig.RecordTypes,
			domainConfig.RecordNames,
			domainConfig.Period)
		go updater.Run(&wg)
	}
	wg.Wait()

}
