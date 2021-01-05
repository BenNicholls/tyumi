package log

import (
	"fmt"
	"os"
	"time"
)

var logger []entry

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

type entry struct {
	level   level
	time    time.Time
	message string
}

func (m entry) String() string {
	return "[" + m.time.Format(time.Stamp) + "] " + m.level.String() + ": " + m.message
}

func init() {
	logger = make([]entry, 0, 1000)
	Info("Tyumi Logger Initialized!")
}

func log(level level, text ...interface{}) {
	logger = append(logger, entry{
		level:   level,
		time:    time.Now(),
		message: fmt.Sprint(text...),
	})

	// //if we're in debug mode, add the new message to the debugger window
	// if debug {
	// 	debugger.logList.Append(logger[len(logger)-1].String())
	// 	debugger.logList.ScrollToBottom()
	// }
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
