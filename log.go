package graylog

import (
	"log"
	"os"
)

var (
	logPrefix   = "[graylog-hook] "
	errorLogger = log.New(os.Stderr, logPrefix, log.LstdFlags)
	debugLogger = log.New(os.Stdout, logPrefix, log.LstdFlags)
)
