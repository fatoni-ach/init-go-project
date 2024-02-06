package services

import (
	"runtime/debug"

	log "github.com/sirupsen/logrus"
)

// WriteLog Writing Log
func WriteLog(message interface{}) {
	log.WithField("stack trace: ", string(debug.Stack())).Info(message)
}

// WriteLogWarn Writing Log
func WriteLogWarn(message interface{}) {
	log.WithField("stack trace: ", string(debug.Stack())).Warn(message)
}

// WriteLogFatal Writing Log
func WriteLogFatal(message interface{}) {
	log.WithField("stack trace: ", string(debug.Stack())).Fatal(message)
}

// RecoverPanic Writing Log
func RecoverPanic() {
	if r := recover(); r != nil {
		log.WithField("panic", r).WithField("stack trace: ", string(debug.Stack())).Error("we panicked!")
	}
}
