//Log package is a logger for debugging purposes. It collects log messages of varying degrees and stores them, with the
//option to output them to disk at the user's discretion.
package log

import (
	"fmt"
	"os"
	"time"
)

var logger []Entry
var onMessageCallback func(e Entry)

type level int

const (
	DEBUG level = iota
	INFO
	WARN
	ERROR
)

func (l level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARNING"
	case ERROR:
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

func (m Entry) String() string {
	return "[" + m.Time.Format(time.Stamp) + "] " + m.Level.String() + ": " + m.Message
}

func init() {
	logger = make([]Entry, 0, 1000)
	Info("Tyumi Logger Initialized!")
}

func log(level level, text ...interface{}) {
	e := Entry{
		Level:   level,
		Time:    time.Now(),
		Message: fmt.Sprint(text...),
	}	
	logger = append(logger, e)

	if onMessageCallback != nil {
		onMessageCallback(e)
	}
}

func outputToDisk() {
	f, err := os.Create("log.txt")
	if err != nil {
		return
	}
	defer f.Close()

	for _, m := range logger {
		f.WriteString(m.String() + "\n")
	}
}

func Debug(m ...interface{}) {
	log(DEBUG, m...)
}

func Info(m ...interface{}) {
	log(INFO, m...)
}

func Warning(m ...interface{}) {
	log(WARN, m...)
}

func Error(m ...interface{}) {
	log(ERROR, m...)
}

//Use this to run code whenever a log entry is recorded, for example to print the log message to the screen.
func SetOnMessageCallback(f func(e Entry)) {
	onMessageCallback = f
}