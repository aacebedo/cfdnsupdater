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

package utils

import (
	"github.com/op/go-logging"
	"os"
	"path/filepath"
)

type LoggingMode int

const (
	Consolelog  LoggingMode = 0 << iota
	Filelog     LoggingMode = 1
	Syslog      LoggingMode = 2
)

var logger = logging.MustGetLogger("cfdnsupdater.utils")

func InitLoggers(verbose bool, quiet bool, loggingMode LoggingMode) (err error) {
	var format logging.Formatter

	var backend logging.Backend

	switch loggingMode {
	case Syslog:
		backend, err = logging.NewSyslogBackend("cfdnsupdater")
		if err != nil {
			return
		}
		format = logging.MustStringFormatter(`%{message}`)
	case Consolelog:
		var outstream *os.File
		if quiet {
			outstream = os.NewFile(uintptr(3), filepath.Join("/","dev","null"))
		} else {
			outstream = os.Stderr
		}
		backend = logging.NewLogBackend(outstream, "", 0)
		format = logging.MustStringFormatter(`%{color}%{time:15:04:05.000} | %{level:.10s} ▶%{color:reset} %{message}`)
	case Filelog:
	  var file *os.File
	  file, err = os.OpenFile(filepath.Join("/","var","log","cfdnsupdater"), os.O_RDWR|os.O_CREATE, 0600)
	  if err != nil {
	    return
	  }
	  backend = logging.NewLogBackend(file, "", 0)
		format = logging.MustStringFormatter(`%{color}%{time:15:04:05.000} | %{level:.10s} ▶%{color:reset} %{message}`)
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
