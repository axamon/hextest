package main

import (
	"os"
	"time"
)

// ExitCodes to convey status.
const (
	OK int = iota
	Warning
	Error
)

func main() {
	minute := time.Now().Minute()

	var exitCode int
	switch {
	case minute%3 == 0:
		exitCode = Error
	case minute%2 == 0:
		exitCode = Warning
	}
	os.Exit(exitCode)
}
