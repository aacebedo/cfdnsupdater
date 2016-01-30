package main

import (
	. "cfdnsupdater/core"
	"os"
)

func main() {
	var parsedArguments = ParseParameters(os.Args[1:])
	InitLoggers(parsedArguments.Verbose, parsedArguments.Quiet, parsedArguments.Syslog)
	SpawnDomainUpdaters(parsedArguments.DomainConfigs)
}
