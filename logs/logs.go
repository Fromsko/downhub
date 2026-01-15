package logs

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Fromsko/downhub/config"
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
	cfg        *config.Config
)

// SetConfig sets the configuration for logging
func SetConfig(c *config.Config) {
	cfg = c
}

func log(level Level, msg string, args ...any) {
	// Check if we should log this level based on config
	if cfg != nil {
		configLevel := strings.ToLower(cfg.Logging.Level)
		currentLevel := ""
		switch level {
		case LevelInfo:
			currentLevel = "info"
		case LevelWarn:
			currentLevel = "warn"
		case LevelError:
			currentLevel = "error"
		}

		// Skip if current level is lower than configured level
		if shouldSkipLog(currentLevel, configLevel) {
			return
		}
	}

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

	// Determine output destination
	output := os.Stdout
	if cfg != nil && cfg.Logging.Output == "stderr" {
		output = os.Stderr
	}

	// Format based on config
	logMsg := fmt.Sprintf(msg, args...)
	if cfg != nil && cfg.Logging.Format == "json" {
		fmt.Fprintf(output, `{"timestamp":"%s","level":"%s","message":"%s"}\n`, t, strings.TrimSpace(levelStr), logMsg)
	} else {
		fmt.Fprintf(output, "%s - %s - %s\n", t, levelStr, logMsg)
	}
}

// shouldSkipLog determines if a log message should be skipped based on level
func shouldSkipLog(currentLevel, configLevel string) bool {
	levels := map[string]int{"info": 0, "warn": 1, "error": 2}
	current := levels[currentLevel]
	configured := levels[configLevel]
	return current < configured
}

func Info(msg string, args ...any)  { log(LevelInfo, msg, args...) }
func Warn(msg string, args ...any)  { log(LevelWarn, msg, args...) }
func Error(msg string, args ...any) { log(LevelError, msg, args...) }
