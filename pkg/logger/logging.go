package logger

import "fmt"

// simple verbose/non-verbose logging control
const (
	VERBOSE MessageType = iota
	PROMPT
	WARN
	ERROR
)

var (
	LogLevel int = int(PROMPT)
)

type MessageType int

func getLevelColor(msgType MessageType) string {
	colorRed := "\033[31m"
	colorGreen := "\033[32m"
	colorYellow := "\033[33m"

	switch msgType {
	case VERBOSE:
		return colorGreen
	case WARN:
		return colorYellow
	case ERROR:
		return colorRed
	}

	return ""
}

func Log(msgType MessageType, format string, args ...interface{}) {
	colorReset := "\033[0m"

	if msgType >= MessageType(LogLevel) {
		color := getLevelColor(msgType)
		if color != "" {
			fmt.Print(color)
		}
		fmt.Print("mockcompose - ")
		fmt.Printf(format, args...)

		if color != "" {
			fmt.Print(colorReset)
		}
	}
}
