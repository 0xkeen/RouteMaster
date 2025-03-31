package utils

import "github.com/inconshreveable/log15"

type Loggable struct {
	ModuleName string
}

func (l *Loggable) log15() log15.Logger {
	return log15.New("module", l.ModuleName)
}

func (l *Loggable) Debug(msg string, ctx ...interface{}) {
	log15.Debug(msg, ctx...)
}
func (l *Loggable) Info(msg string, ctx ...interface{}) {
	log15.Info(msg, ctx...)
}
func (l *Loggable) Warn(msg string, ctx ...interface{}) {
	log15.Warn(msg, ctx...)
}
func (l *Loggable) Error(msg string, ctx ...interface{}) {
	log15.Error(msg, ctx...)
}
func (l *Loggable) Crit(msg string, ctx ...interface{}) {
	log15.Crit(msg, ctx...)
}
