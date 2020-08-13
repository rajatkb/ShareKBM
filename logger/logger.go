package logger

import (
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
)

//LogType ... for logtype
type LogType int8

const (
	fatal LogType = 0
	info  LogType = 1
	warn  LogType = 2
	err   LogType = 3
)

//LogMessage ... for message type bundling
type LogMessage struct {
	logt    LogType
	message string
}

//Logger ... represent a Logger instance
type Logger struct {
	logchannel         chan LogMessage
	waitOnChannelEmpty sync.WaitGroup
}

func (l *Logger) printlogs(logt LogType, message string) {
	switch logt {
	case fatal:
		{
			log.Fatal(message)
		}
	case info:
		{
			log.Info(message)
		}
	case warn:
		{
			log.Warn(message)
		}
	case err:
		{
			log.Warn(message)
		}
	}
}

//Close ... close Logger channels
func (l *Logger) Close() {
	if l.logchannel != nil {
		close(l.logchannel) // always before for loop so that range knows about the channel being closed
		l.logchannel = nil
		l.waitOnChannelEmpty.Wait()
	}
}

func (l *Logger) readlogs() {
	go func() {
		for lm := range l.logchannel {
			l.printlogs(lm.logt, lm.message)
			l.waitOnChannelEmpty.Done()
		}
	}()
}

//Fatal ... logs fatal and shuts down app
func (l *Logger) Fatal(message string) {
	if l.logchannel != nil {
		l.logchannel <- LogMessage{logt: fatal, message: message}
		l.waitOnChannelEmpty.Add(1)
	}
}

//Info ... logs info
func (l *Logger) Info(message string) {

	if l.logchannel != nil {
		l.logchannel <- LogMessage{logt: info, message: message}
		l.waitOnChannelEmpty.Add(1)
	}

}

//Warn ... logs warn
func (l *Logger) Warn(message string) {

	if l.logchannel != nil {
		l.logchannel <- LogMessage{logt: warn, message: message}
		l.waitOnChannelEmpty.Add(1)
	}
}

//Error ... logs error
func (l *Logger) Error(message string) {
	if l.logchannel != nil {
		l.logchannel <- LogMessage{logt: err, message: message}
		l.waitOnChannelEmpty.Add(1)
	}
}

//GetLogger ... provide a Logger object
func GetLogger(bufferSize int) *Logger {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	lgr := new(Logger)
	lgr.logchannel = make(chan LogMessage, bufferSize)
	lgr.readlogs() // start go channel for reader
	return lgr
}
