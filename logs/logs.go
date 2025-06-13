package logs

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
)

type Level int

const (
	LevelInfo Level = iota
	LevelWarn
	LevelError
)

var (
	infoColor  = color.New(color.FgGreen).SprintFunc()
	warnColor  = color.New(color.FgYellow).SprintFunc()
	errorColor = color.New(color.FgRed).SprintFunc()
)

func log(level Level, msg string, args ...any) {
	t := time.Now().Format("2006-01-02 15:04:05")
	var levelStr string
	switch level {
	case LevelInfo:
		levelStr = infoColor("INFO ")
	case LevelWarn:
		levelStr = warnColor("WARN ")
	case LevelError:
		levelStr = errorColor("ERROR")
	}
	fmt.Fprintf(os.Stdout, "%s - %s - %s\n", t, levelStr, fmt.Sprintf(msg, args...))
}

func Info(msg string, args ...any)  { log(LevelInfo, msg, args...) }
func Warn(msg string, args ...any)  { log(LevelWarn, msg, args...) }
func Error(msg string, args ...any) { log(LevelError, msg, args...) }
