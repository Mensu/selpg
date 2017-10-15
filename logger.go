package main

import (
	"fmt"
	"io"
)

// Logger helper type to wrap fmt.PrintXX
type Logger struct {
	Progname string
	Stdout   io.Writer
	Stderr   io.Writer
}

// Print wraps fmt.Fprint(stdout, ...)
func (logger *Logger) Print(args ...interface{}) (n int, err error) {
	return fmt.Fprint(logger.Stdout, args...)
}

// Printf wraps fmt.Fprintf(stdout, ...)
func (logger *Logger) Printf(format string, args ...interface{}) (n int, err error) {
	return fmt.Fprintf(logger.Stdout, format, args...)
}

// Println wraps fmt.Fprintln(stdout, ...)
func (logger *Logger) Println(args ...interface{}) (n int, err error) {
	return fmt.Fprintln(logger.Stdout, args...)
}

func (logger *Logger) errPrintHeader() (n int, err error) {
	return fmt.Fprintf(logger.Stderr, "%s: ", logger.Progname)
}

// ErrPrint wraps fmt.FPrint(os.Stderr, ...)
func (logger *Logger) ErrPrint(args ...interface{}) (n int, err error) {
	n, err = logger.errPrintHeader()
	if err != nil {
		return
	}
	n2, err := fmt.Fprint(logger.Stderr, args...)
	n += n2
	return
}

// ErrPrintf wraps fmt.FPrintf(os.Stderr, ...)
func (logger *Logger) ErrPrintf(format string, args ...interface{}) (n int, err error) {
	n, err = logger.errPrintHeader()
	if err != nil {
		return
	}
	n2, err := fmt.Fprintf(logger.Stderr, format, args...)
	n += n2
	return
}

// ErrPrintln wraps fmt.FPrintln(os.Stderr, ...)
func (logger *Logger) ErrPrintln(args ...interface{}) (n int, err error) {
	n, err = logger.errPrintHeader()
	if err != nil {
		return
	}
	n2, err := fmt.Fprintln(logger.Stderr, args...)
	n += n2
	return
}
