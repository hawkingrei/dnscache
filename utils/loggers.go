package utils

import (
	"io/ioutil"
	"os"

	"github.com/op/go-logging"
)

// InitLoggers initialize loggers
func InitLoggers(verbosity int) (err error) {
	var format logging.Formatter

	var backend logging.Backend

	switch {
	case verbosity == 0:
		backend = logging.NewLogBackend(ioutil.Discard, "", 0)
	case verbosity >= 1:
		backend = logging.NewLogBackend(os.Stdout, "", 0)
	}

	format = logging.MustStringFormatter(`%{color}%{time:15:04:05.000} | %{level:.10s} ▶%{color:reset} %{message}`)

	formatter := logging.NewBackendFormatter(backend, format)
	leveledBackend := logging.AddModuleLevel(formatter)

	switch {
	case verbosity == 1:
		leveledBackend.SetLevel(logging.INFO, "")
	case verbosity >= 2:
		leveledBackend.SetLevel(logging.DEBUG, "")
	}

	logging.SetBackend(leveledBackend)
	return
}
