/**
 *s class is a simple wrapper for adding a few use things to go
 * logging, and help keep similarity between the other Ola Client libraries and
 * this one
 *
 * It provides a global logger, that is initialzed on init.
 *
 * Right now its a simple wrapper over logging but intended to get added to
 * later, and allow a consistent api
 */
package logging

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime"
	"path/filepath"
	"errors"
)

var log_color bool
var logger *Logger

func init() {
	_log := log.New(os.Stdout, "", 0)
	logger = NewLogger(_log)
}

type LogLevel uint

const (
	LOG_NONE  LogLevel = iota
	LOG_FATAL
	LOG_WARN
	LOG_INFO
	LOG_DEBUG
)

func IntToLogLevel(value int) (error, LogLevel) {
  if value == 0 {
    return nil, LOG_NONE
  } else if value == 1 {
    return nil, LOG_FATAL
  } else if value == 2 {
    return nil, LOG_WARN
  } else if value == 3 {
    return nil, LOG_INFO
  } else if value == 4 {
    return nil, LOG_DEBUG
  } else {
    return errors.New("Invalid Number for a Log Level, must be 0 to 4"),
      LOG_NONE
  }
}

func (e LogLevel) LogColor() string {
	switch e {
	case LOG_FATAL:
		return "\033[31m"
	case LOG_WARN:
		return "\033[33m"
	case LOG_INFO:
		return "\033[32m"
	}
	return ""
}

func (e LogLevel) String() string {
	switch e {
	case LOG_FATAL:
		return "Fatal"
	case LOG_WARN:
		return "Warn"
	case LOG_INFO:
		return "Info"
	case LOG_DEBUG:
		return "Debug"
	}
	return ""
}

type Logger struct {
	logger    *log.Logger
	log_level LogLevel
}

func NewLogger(log_interface *log.Logger) *Logger {
	l := new(Logger)
	l.logger = log_interface
	l.log_level = LOG_NONE
	return l
}

func (m *Logger) SetLoggingLevel(level LogLevel) {
	if level > LOG_DEBUG {
		return
	}
	m.log_level = level
}

func (m *Logger) log(log_type LogLevel, call_depth int, msg string, args ...interface{}) {
	if log_type > m.log_level {
		return
	}

	if m.logger != nil {
		_, file, line, ok := runtime.Caller(call_depth)
		if ok == false {
			m.logger.Printf(msg, args...)
			return
		}
		var buffer bytes.Buffer
		_, filename := filepath.Split(file)
		if log_color {	
			buffer.WriteString(fmt.Sprintf("%s%s:%d:%s: \033[0m", log_type.LogColor(), filename,
				line, log_type))
		} else {
			buffer.WriteString(fmt.Sprintf("%s:%d:%s: ", filename, line, log_type))
		}
		buffer.WriteString(msg)
		m.logger.Printf(buffer.String(), args...)
	}
}

// Allows us to take over the output from the log package to the global logger
// 
func (m *Logger) Write(p []byte) (int, error) {
  logger.log(LOG_INFO, 4, "External Libs: %s", p)
  return len(p), nil
}

func SetLoggingLevel(level LogLevel) {
 	logger.SetLoggingLevel(level)
}

func Debug(msg string, args ...interface{}) {
	logger.log(LOG_DEBUG, 2, msg, args...)
}

func Warn(msg string, args ...interface{}) {
	logger.log(LOG_WARN, 2, msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	logger.log(LOG_FATAL, 2, msg, args...)
}

func Info(msg string, args ...interface{}) {
	logger.log(LOG_INFO, 2, msg, args...)
}

func GetLogger() *Logger {
  return logger
}

func SetColorEnabled(enabled bool) {
	log_color = enabled
}

func ColorEnabled() bool {
 	return log_color
}