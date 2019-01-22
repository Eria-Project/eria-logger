package logger

import (
	"errors"
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var origLogger = logrus.New()
var baseLogger = logger{entry: logrus.NewEntry(origLogger)}
var disabled bool

func init() {
	formater := logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "02/01|15:04:05",
	}
	origLogger.SetFormatter(&formater)
}

// Fields is a wrapper for logrus.Fields
type Fields logrus.Fields

// Logger interface
type Logger interface {
	Trace(...interface{})
	Tracef(string, ...interface{})
	Debug(...interface{})
	Debugf(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	Warn(...interface{})
	Warnf(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})
	WithError(err error) Logger
	WithField(key string, value interface{}) Logger
	WithFields(fields Fields) Logger
	WithMessage(msg string) Logger
}

type logger struct {
	entry *logrus.Entry
}

func (l logger) sourced(source bool) *logrus.Entry {
	e := l.entry

	e = e.WithField("module", "knx")

	if !source {
		return e
	}

	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		e = e.WithField("source", "<???>")
		return e
	}

	slash := strings.LastIndex(file, "/")
	file = file[slash+1:]
	e = e.WithField("source", fmt.Sprintf("%s:%d", file, line))

	if f := runtime.FuncForPC(pc); f != nil {
		e = e.WithField("func", f.Name())
	}

	return e
}

// SetLevel set the logger log level
func SetLevel(level string) {
	logLevel, err := logrus.ParseLevel(level)
	if level == "no" {
		disabled = true
		origLogger.Out = ioutil.Discard
	} else {
		if err != nil {
			origLogger.WithField("level", level).Warn(err)
		} else {
			origLogger.SetLevel(logLevel)
		}
	}
}

func (l logger) Trace(args ...interface{}) {
	l.sourced(false).Trace(args...)
}

func (l logger) Tracef(format string, args ...interface{}) {
	l.sourced(false).Tracef(format, args...)
}

func (l logger) Debug(args ...interface{}) {
	l.sourced(false).Debug(args...)
}

func (l logger) Debugf(format string, args ...interface{}) {
	l.sourced(false).Debugf(format, args...)
}

func (l logger) Info(args ...interface{}) {
	l.sourced(false).Info(args...)
}

func (l logger) Infof(format string, args ...interface{}) {
	l.sourced(false).Infof(format, args...)
}

func (l logger) Warn(args ...interface{}) {
	l.sourced(false).Warn(args...)
}

func (l logger) Warnf(format string, args ...interface{}) {
	l.sourced(false).Warnf(format, args...)
}

func (l logger) Error(args ...interface{}) {
	l.sourced(true).Error(args...)
}

func (l logger) Errorf(format string, args ...interface{}) {
	l.sourced(true).Errorf(format, args...)
}

func (l logger) Fatal(args ...interface{}) {
	l.sourced(true).Fatal(args...)
}

func (l logger) Fatalf(format string, args ...interface{}) {
	l.sourced(true).Fatalf(format, args...)
}

func (l logger) WithField(key string, value interface{}) Logger {
	return logger{l.entry.WithField(key, value)}
}

func (l logger) WithFields(fields Fields) Logger {
	return logger{l.entry.WithFields(map[string]interface{}(fields))}
}

func (l logger) WithError(err error) Logger {
	return logger{l.entry.WithError(err)}
}

func (l logger) WithMessage(msg string) Logger {
	return logger{l.entry.WithError(errors.New(msg))}
}

func (l logger) WithMessagef(format string, args ...interface{}) Logger {
	return logger{l.entry.WithError(fmt.Errorf(format, args...))}
}

// Trace wrapper
func Trace(args ...interface{}) {
	baseLogger.sourced(false).Trace(args...)
}

// Tracef wrapper
func Tracef(format string, args ...interface{}) {
	baseLogger.sourced(false).Tracef(format, args...)
}

// Debug wrapper
func Debug(args ...interface{}) {
	baseLogger.sourced(false).Debug(args...)
}

// Debugf wrapper
func Debugf(format string, args ...interface{}) {
	baseLogger.sourced(false).Debugf(format, args...)
}

// Info wrapper
func Info(args ...interface{}) {
	baseLogger.sourced(false).Info(args...)
}

// Infof wrapper
func Infof(format string, args ...interface{}) {
	baseLogger.sourced(false).Infof(format, args...)
}

// Warn wrapper
func Warn(args ...interface{}) {
	baseLogger.sourced(false).Warn(args...)
}

// Warnf wrapper
func Warnf(format string, args ...interface{}) {
	baseLogger.sourced(false).Warnf(format, args...)
}

// Error wrapper
func Error(args ...interface{}) {
	baseLogger.sourced(true).Error(args...)
}

// Errorf wrapper
func Errorf(format string, args ...interface{}) {
	baseLogger.sourced(true).Errorf(format, args...)
}

// Fatal wrapper
func Fatal(args ...interface{}) {
	baseLogger.sourced(true).Fatal(args...)
}

// Fatalf wrapper
func Fatalf(format string, args ...interface{}) {
	baseLogger.sourced(true).Fatalf(format, args...)
}

// WithField wrapper
func WithField(key string, value interface{}) Logger {
	return baseLogger.WithField(key, value)
}

// WithFields wrapper
func WithFields(fields Fields) Logger {
	return baseLogger.WithFields(fields)
}

// WithError wrapper
func WithError(err error) Logger {
	return baseLogger.WithError(err)
}

// WithMessage wrapper
func WithMessage(msg string) Logger {
	return baseLogger.WithMessage(msg)
}

// WithMessagef wrapper
func WithMessagef(format string, args ...interface{}) Logger {
	return baseLogger.WithMessagef(format, args...)
}
