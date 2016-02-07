package utils

import (
	"github.com/op/go-logging"
	"os"
)

var logger = logging.MustGetLogger("cfdnsupdater.utils")

func InitLoggers(verbose bool, quiet bool, syslog bool) (err error){
	var format logging.Formatter
	
  var backend logging.Backend
  if syslog {
    backend, err = logging.NewSyslogBackend("cfdnsupdater")
    if err != nil {
      return
    }
    format = logging.MustStringFormatter(`%{message}`)
    
  } else {
    var outstream *os.File
  	if quiet {
  		outstream = os.NewFile(uintptr(3), "/dev/null")
  	} else {
  		outstream = os.Stderr
  	}
    backend = logging.NewLogBackend(outstream, "", 0)
    format =  logging.MustStringFormatter(`%{color}%{time:15:04:05.000} | %{level:.10s} ▶%{color:reset} %{message}`)
  }
  
	formatter := logging.NewBackendFormatter(backend, format)
	leveledBackend := logging.AddModuleLevel(formatter)

	if verbose {
		leveledBackend.SetLevel(logging.DEBUG, "")
	} else {
		leveledBackend.SetLevel(logging.INFO, "")
	}

	logging.SetBackend(leveledBackend)
  return
}
