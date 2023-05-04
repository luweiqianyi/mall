package log

import "log"

func PrintLog(format string, v ...any) {
	log.Printf(format, v...)
}
