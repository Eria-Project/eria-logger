package logger

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	prefixed "github.com/gz-c/logrus-prefixed-formatter"
	"github.com/sirupsen/logrus"
)

var (
	_origLogger *logrus.Logger
	_baseLogger logger
	_logLevel   *string
	_logFile    *string
	_paramsSet  bool
)

func init() {
	_origLogger = logrus.New()
	_baseLogger = logger{entry: logrus.NewEntry(_origLogger)}
	_logLevel = flag.String("log", "info", "log level [info, debug, error, warn, no]")
	_logFile = flag.String("log-file", "", "log file path (default to stderr")

	formatter := prefixed.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "02/01|15:04:05",
	}
	formatter.SetColorScheme(&prefixed.ColorScheme{
		InfoLevelStyle:   "green",
		WarnLevelStyle:   "orange",
		ErrorLevelStyle:  "red",
		FatalLevelStyle:  "red",
		PanicLevelStyle:  "red",
		DebugLevelStyle:  "blue",
		PrefixStyle:      "cyan",
		CallContextStyle: "white",
		TimestampStyle:   "white+h",
	})
	_origLogger.SetFormatter(&formatter)

}

func setParams() {
	if !_paramsSet {
		if !flag.Parsed() {
			flag.Parse()
		}

		if *_logFile != "" {
			SetFile(*_logFile)
		}
		SetLevel(*_logLevel)

		Module("logger").WithField("path", *_logFile).Info("Set log file")
		Module("logger").WithField("level", *_logLevel).Info("Set log level")
		_paramsSet = true
	}
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

	if !source {
		return e
	}

	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		e = e.WithField("file", "<???>")
		return e
	}

	slash := strings.LastIndex(file, "/")
	file = file[slash+1:]
	e = e.WithField("file", file)
	e = e.WithField("line", line)

	if f := runtime.FuncForPC(pc); f != nil {
		e = e.WithField("func", f.Name())
	}

	return e
}

// SetLevel set the logger log level
func SetLevel(level string) {
	logLevel, err := logrus.ParseLevel(level)
	if level == "no" {
		_origLogger.Out = ioutil.Discard
	} else {
		if err != nil {
			_origLogger.WithField("level", level).Warn(err)
		} else {
			_origLogger.SetLevel(logLevel)
		}
	}
}

// SetFile set the logger to output to file
func SetFile(path string) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		_origLogger.SetOutput(file)
	} else {
		Module("logger").Info("Failed to log to file, using default stderr")
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

// Module wrapper
func Module(name string) Logger {
	return _baseLogger.WithField("prefix", name)
}
