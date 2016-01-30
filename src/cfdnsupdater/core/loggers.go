package core

import (
	"github.com/op/go-logging"
	"os"
)

var RootLogger = logging.MustGetLogger("cfdnsupdater")

func InitLoggers(verbose bool, quiet bool, syslog bool) {
	var format = logging.MustStringFormatter(`%{color}%{time:15:04:05.000} | %{level:.10s} ▶%{color:reset} %{message}`)
	var outstream *os.File
	if quiet {
		outstream = os.NewFile(uintptr(3), "/dev/null")
	} else {
		outstream = os.Stderr
	}

	var backend = logging.NewLogBackend(outstream, "", 0)
	var formatter = logging.NewBackendFormatter(backend, format)
	var leveledBackend = logging.AddModuleLevel(formatter)

	if verbose {
		leveledBackend.SetLevel(logging.DEBUG, "")
	} else {
		leveledBackend.SetLevel(logging.INFO, "")
	}

	logging.SetBackend(leveledBackend)

}
