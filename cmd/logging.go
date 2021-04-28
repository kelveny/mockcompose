package cmd

import "fmt"

// simple verbose/non-verbose logging control
const (
	VERBOSE MessageType = iota
	PROMPT
)

var (
	LogLevel int = int(PROMPT)
)

type MessageType int

func Log(msgType MessageType, format string, args ...interface{}) {
	if msgType >= MessageType(LogLevel) {
		fmt.Printf(format, args...)
	}
}
