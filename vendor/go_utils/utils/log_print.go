package utils

import (
	"fmt"
	"go_utils/utils/screen"
	"time"
)

const (
	LOG_LEVEL_ERROR = iota
	LOG_LEVEL_WARNING
	LOG_LEVEL_INFO
	LOG_LEVEL_DEBUG
	LOG_LEVEL_DEBUG2
	LOG_LEVEL_DEBUG3
	LOG_LEVEL_DEBUG4
)

var LogLevel = LOG_LEVEL_INFO

func LogPrintDebug(a ...any) {
	if LogLevel >= LOG_LEVEL_DEBUG {
		t := time.Now().Format("2006-01-02 15:04:05.000")
		screen.ColoredPrint(screen.ANSI_COLOR_DEFAULT, fmt.Sprintf("[%s][%s] ", "DEBUG", t))
		fmt.Println(a...)
	}
}

func LogPrintDebug2(a ...any) {
	if LogLevel >= LOG_LEVEL_DEBUG2 {
		t := time.Now().Format("2006-01-02 15:04:05.000")
		screen.ColoredPrint(screen.ANSI_COLOR_DEFAULT, fmt.Sprintf("[%s][%s] ", "DEBUG2", t))
		fmt.Println(a...)
	}
}

func LogPrintDebug3(a ...any) {
	if LogLevel >= LOG_LEVEL_DEBUG3 {
		t := time.Now().Format("2006-01-02 15:04:05.000")
		screen.ColoredPrint(screen.ANSI_COLOR_DEFAULT, fmt.Sprintf("[%s][%s] ", "DEBUG3", t))
		fmt.Println(a...)
	}
}

func LogPrintDebug4(a ...any) {
	if LogLevel >= LOG_LEVEL_DEBUG4 {
		t := time.Now().Format("2006-01-02 15:04:05.000")
		screen.ColoredPrint(screen.ANSI_COLOR_DEFAULT, fmt.Sprintf("[%s][%s] ", "DEBUG4", t))
		fmt.Println(a...)
	}
}

func LogPrintInfo(a ...any) {
	if LogLevel >= LOG_LEVEL_INFO {
		t := time.Now().Format("2006-01-02 15:04:05.000")
		screen.ColoredPrint(screen.ANSI_COLOR_CYAN, fmt.Sprintf("[%s][%s] ", "INFO", t))
		fmt.Println(a...)
	}
}

func LogPrintWarning(a ...any) {
	if LogLevel >= LOG_LEVEL_WARNING {
		t := time.Now().Format("2006-01-02 15:04:05.000")
		screen.ColoredPrint(screen.ANSI_COLOR_YELLOW+screen.ANSI_STYLE_BOLD, fmt.Sprintf("[%s][%s] ", "WARNING", t))
		fmt.Println(a...)
	}
}

func LogPrintError(a ...any) {
	if LogLevel >= LOG_LEVEL_ERROR {
		t := time.Now().Format("2006-01-02 15:04:05.000")
		screen.ColoredPrint(screen.ANSI_COLOR_RED+screen.ANSI_STYLE_BOLD, fmt.Sprintf("[%s][%s] ", "ERROR", t))
		fmt.Println(a...)
	}
}

func SetLogLevelDebug() {
	LogLevel = LOG_LEVEL_DEBUG
}
func SetLogLevelDebug2() {
	LogLevel = LOG_LEVEL_DEBUG2
}
func SetLogLevelDebug3() {
	LogLevel = LOG_LEVEL_DEBUG3
}
func SetLogLevelDebug4() {
	LogLevel = LOG_LEVEL_DEBUG4
}
func SetLogLevelInfo() {
	LogLevel = LOG_LEVEL_INFO
}
func SetLogLevelWarning() {
	LogLevel = LOG_LEVEL_WARNING
}
func SetLogLevelError() {
	LogLevel = LOG_LEVEL_ERROR
}
