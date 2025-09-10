// Package utils provides helper functions and utilities for the bot.
// This file, logger.go, implements a simple, colorful, and formatted logger
// for displaying messages at different severity levels (Info, Error, Warn, Debug).
package utils

import (
	"fmt"
	"time"
)

// ANSI color codes for formatting log messages in the console.
const (
	Reset  = "\033[0m"
	Green  = "\033[32m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Gray   = "\033[90m"
)

// timestamp generates a formatted string of the current time in the "Asia/Jakarta" (WIB) timezone.
// It falls back to a fixed UTC+7 offset if the timezone cannot be loaded.
//
// Returns:
//   A string representing the current timestamp (e.g., "2006-01-02 15:04:05").
func timestamp() string {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.FixedZone("WIB", 7*3600)
	}
	return time.Now().In(loc).Format("2006-01-02 15:04:05")
}

// Info logs a message at the INFO level with a green color.
// Use this for general application information.
//
// Parameters:
//   msg: The message string to be logged.
func Info(msg string) {
	fmt.Printf("%s[%s] [INFO] %s%s\n", Green, timestamp(), msg, Reset)
}

// Error logs a message at the ERROR level with a red color.
// Use this for critical errors that prevent normal operation.
//
// Parameters:
//   msg: The message string to be logged.
func Error(msg string) {
	fmt.Printf("%s[%s] [ERROR] %s%s\n", Red, timestamp(), msg, Reset)
}

// Warn logs a message at the WARN level with a yellow color.
// Use this for potential issues that do not stop the application.
//
// Parameters:
//   msg: The message string to be logged.
func Warn(msg string) {
	fmt.Printf("%s[%s] [WARN] %s%s\n", Yellow, timestamp(), msg, Reset)
}

// Debug logs a message at the DEBUG level with a cyan color.
// Use this for detailed, verbose information useful for debugging.
//
// Parameters:
//   msg: The message string to be logged.
func Debug(msg string) {
	fmt.Printf("%s[%s] [DEBUG] %s%s\n", Cyan, timestamp(), msg, Reset)
}
