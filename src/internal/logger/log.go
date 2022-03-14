package logger

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func PrepareLogger(configLogLevel string) {
	log.SetFormatter(&log.TextFormatter{
		DisableColors:  false,
		DisableQuote:   true,
		FullTimestamp:  true,
		DisableSorting: false,
	})
	level := getLogLevel(configLogLevel)
	log.SetLevel(level)
	log.Debug(fmt.Sprintf("Log level was set to %s", level.String()))
}

func getLogLevel(logLevel string) log.Level {
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Fatal(err)
	}
	return level
}
