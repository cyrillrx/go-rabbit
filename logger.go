package rabbitmq

import "io"

type Logger interface {
	// Implements io.Writer so it can be set a output of Logger
	io.Writer

	// Debug logs a debug message. If last parameter is a map[string]string, it's content
	// is added as fields to the message.
	Debug(v ...interface{})
	// Debug logs a debug message with format. If last parameter is a map[string]string,
	// it's content is added as fields to the message.
	Debugf(format string, v ...interface{})
	// Info logs a info message. If last parameter is a map[string]string, it's content
	// is added as fields to the message.
	Info(v ...interface{})
	// Info logs a info message with format. If last parameter is a map[string]string,
	// it's content is added as fields to the message.
	Infof(format string, v ...interface{})
	// Warn logs a warning message. If last parameter is a map[string]string, it's content
	// is added as fields to the message.
	Warn(v ...interface{})
	// Warn logs a warning message with format. If last parameter is a map[string]string,
	// it's content is added as fields to the message.
	Warnf(format string, v ...interface{})
	// Error logs an error message. If last parameter is a map[string]string, it's content
	// is added as fields to the message.
	Error(v ...interface{})
	// Error logs an error message with format. If last parameter is a map[string]string,
	// it's content is added as fields to the message.
	Errorf(format string, v ...interface{})
	// Fatal logs an error message followed by a call to os.Exit(1). If last parameter is a
	// map[string]string, it's content is added as fields to the message.
	Fatal(v ...interface{})
	// Fatalf logs an error message with format followed by a call to ox.Exit(1). If last
	// parameter is a map[string]string, it's content is added as fields to the message.
	Fatalf(format string, v ...interface{})
	// Output mimics std logger interface
	Output(calldepth int, s string) error
}
