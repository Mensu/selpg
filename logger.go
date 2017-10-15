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
func (printer *Logger) Print(args ...interface{}) (n int, err error) {
	return fmt.Fprint(printer.Stdout, args...)
}

// Printf wraps fmt.Fprintf(stdout, ...)
func (printer *Logger) Printf(format string, args ...interface{}) (n int, err error) {
	return fmt.Fprintf(printer.Stdout, format, args...)
}

// Println wraps fmt.Fprintln(stdout, ...)
func (printer *Logger) Println(args ...interface{}) (n int, err error) {
	return fmt.Fprintln(printer.Stdout, args...)
}

func (printer *Logger) errPrintHeader() (n int, err error) {
	return fmt.Fprintf(printer.Stderr, "%s: ", printer.Progname)
}

// ErrPrint wraps fmt.FPrint(os.Stderr, ...)
func (printer *Logger) ErrPrint(args ...interface{}) (n int, err error) {
	n, err = printer.errPrintHeader()
	if err != nil {
		return
	}
	n2, err := fmt.Fprint(printer.Stderr, args...)
	n += n2
	return
}

// ErrPrintf wraps fmt.FPrintf(os.Stderr, ...)
func (printer *Logger) ErrPrintf(format string, args ...interface{}) (n int, err error) {
	n, err = printer.errPrintHeader()
	if err != nil {
		return
	}
	n2, err := fmt.Fprintf(printer.Stderr, format, args...)
	n += n2
	return
}

// ErrPrintln wraps fmt.FPrintln(os.Stderr, ...)
func (printer *Logger) ErrPrintln(args ...interface{}) (n int, err error) {
	n, err = printer.errPrintHeader()
	if err != nil {
		return
	}
	n2, err := fmt.Fprintln(printer.Stderr, args...)
	n += n2
	return
}
