//go:build debug

package log

func init() {
	SetMinimumLogLevel(LVL_DEBUG)
	EnableConsoleOutput()
}

func Debug(messages ...any) {
	log(LVL_DEBUG, messages...)
}
