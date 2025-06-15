// Log package is a logger for debugging purposes. It collects log messages of varying degrees and stores them, with the
// option to output them to disk and/or console at the user's discretion.
package log

import (
	"fmt"
	"os"
	"time"
)

var logger []Entry
var onMessageCallback func(e Entry)  // callback for when a message is logged
var printLogs bool                   // print log messages to console
var minimumLogLevel level = LVL_INFO // only log messages of this level or higher

type level uint8

const (
	LVL_DEBUG level = iota
	LVL_INFO
	LVL_WARN
	LVL_ERROR
)

func (l level) String() string {
	switch l {
	case LVL_DEBUG:
		return "DEBUG"
	case LVL_INFO:
		return "INFO"
	case LVL_WARN:
		return "WARNING"
	case LVL_ERROR:
		return "ERROR"
	default:
		return "???"
	}
}

type Entry struct {
	Level   level
	Time    time.Time
	Message string
}

func (e Entry) String() string {
	return fmt.Sprintf("[%s] %s: %s", e.Time.Format(time.TimeOnly), e.Level, e.Message)
}

func init() {
	logger = make([]Entry, 0, 1000)
}

func log(level level, messages ...any) {
	if level < minimumLogLevel || len(messages) == 0 {
		return
	}

	e := Entry{
		Level:   level,
		Time:    time.Now(),
		Message: fmt.Sprint(messages...),
	}

	if e.Message == "" {
		return
	}

	logger = append(logger, e)

	if onMessageCallback != nil {
		onMessageCallback(e)
	}

	if printLogs {
		fmt.Println(e)
	}
}

// Writes the log to log.txt.
// TODO: make this filepath customizable.
func WriteToDisk() {
	f, err := os.Create("log.txt")
	if err != nil {
		return
	}
	defer f.Close()

	for _, m := range logger {
		f.WriteString(m.String() + "\n")
	}
}

func Info(messages ...any) {
	log(LVL_INFO, messages...)
}

func Warning(messages ...any) {
	log(LVL_WARN, messages...)
}

func Error(messages ...any) {
	log(LVL_ERROR, messages...)
}

// Use this to run code whenever a log entry is recorded.
func SetOnMessageCallback(f func(e Entry)) {
	onMessageCallback = f
}

// EnableConsoleOutput will cause all log messages to be printed to the console.
func EnableConsoleOutput() {
	printLogs = true
}

func SetMinimumLogLevel(l level) {
	if l == minimumLogLevel {
		return
	}

	minimumLogLevel = l
	if l == LVL_DEBUG && !printLogs {
		Debug("Debug logging enabled. Consider running log.EnablePrinting() to print logs to the console!")
	}
}
