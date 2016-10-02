package main

import (
	"fmt"
	"os"
)

// ErrorOut writes a formated error-message to StdErr.
func ErrorOut(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(os.Stderr, "[E] "+format, a...)
}

// PrintError prints the supplied error to StdErr.
func PrintError(err error) (int, error) {
	return ErrorOut("%v\n", err)
}
