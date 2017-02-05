package logging

import (
	"fmt"
	"strings"
)

// logLevelMap maps log levels to comparable integers.
var logLevelMap = map[string]int{
	"TRACE" : 10,
	"DEBUG" : 20,
	"INFO" : 30,
	"WARN" : 40,
	"ERROR" : 50,
	"FATAL" : 60,
}

// Logger interface defines 4 different log levels : trace/debug/info/error
type Logger interface {

	// Logs at Trace level
	Trace(string)

	// Logs at Debug level
	Debug(string)

	// Logs at Info level
	Info(string)

	// Logs at Error level
	Error(string)
}

type SimpleLogger struct{
	logLevelNumeric int
}

func (s *SimpleLogger) Trace(text string) {
	if (s.logLevelNumeric <= logLevelMap["TRACE"]) {
		fmt.Println("Trace : " + strings.Replace(text, "\n", "", 10))
	}
}

func (s *SimpleLogger) Debug(text string) {
	if (s.logLevelNumeric <= logLevelMap["DEBUG"]) {
		fmt.Println("Debug :" + strings.Replace(text, "\n", "", 10))
	}
}

func (s *SimpleLogger) Info(text string) {
	if (s.logLevelNumeric <= logLevelMap["INFO"]) {
		fmt.Println("Info : " + strings.Replace(text, "\n", "", 10))
	}
}

func (s *SimpleLogger) Error(text string) {
	if (s.logLevelNumeric <= logLevelMap["ERROR"]) {
		fmt.Println("Error : " + strings.Replace(text, "\n", "", 10))
	}
}

// NewIndexLogger creates a new Logger for our indexing application.
// Log level is set to info if no level is provided.
func NewIndexLogger(logLevel *string) Logger {
	defaultLogLevel := "INFO"
	if (*logLevel != "TRACE" && *logLevel != "ERROR" && *logLevel != "DEBUG" && *logLevel != "WARN" && *logLevel != "FATAL") {
		logLevel = &defaultLogLevel
	}
	logLevelNumeric := logLevelMap[*logLevel]
	return &SimpleLogger{
		logLevelNumeric,
	}
}
