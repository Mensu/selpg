package main

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"time"
)

var printer *Printer
var err error

func main() {
	// construct global printer
	printer = &Printer{
		os.Args[0],
		os.Stdout,
		os.Stderr,
	}

	// process flags
	err = processFlags()
	if err != nil {
		printer.ErrPrintf("%v\n%s", err, flagset.FlagUsagesWrapped(0))
		os.Exit(2)
	}

	// parsed results
	startPagenum, _ := flagset.GetUint(flagInfo["s"])
	endPagenum, _ := flagset.GetUint(flagInfo["e"])
	pageSize, _ := flagset.GetUint(flagInfo["l"])
	useF, _ := flagset.GetBool(flagInfo["f"])
	dest, _ := flagset.GetString(flagInfo["d"])

	joinSubProc := make(chan time.Time)
	var subproc *exec.Cmd

	// set where to output to
	if len(dest) > 0 {
		subproc = exec.Command("lp", "-d", dest)
		// subproc = exec.Command("cat", "-n")
		printer.Stdout, _ = subproc.StdinPipe()
		subproc.Stdout = os.Stdout
		subproc.Stderr = os.Stderr
		// when subprocess is done, join it
		go func() {
			subproc.Run()
			joinSubProc <- time.Now()
		}()
	}

	// set where to input from
	var input = os.Stdin
	var readPagenum uint
	if len(filename) > 0 {
		input, err = os.Open(filename)
		defer input.Close()
		if err != nil {
			printer.ErrPrintf("Unexpected error when opening %s: %v\n", filename, err)
			os.Exit(1)
		}
	}
	reader := bufio.NewReader(input)

	// select pages and print
	if useF {
		readPagenum, err = selpgByF(reader, startPagenum, endPagenum)
	} else {
		readPagenum, err = selpgByLine(reader, startPagenum, endPagenum, pageSize)
	}

	// inform errors
	if err == errStartOutOfRange {
		printer.ErrPrintf("start_page (%d) greater than total pages (%d), no output written\n", startPagenum, readPagenum)
	} else if err == errEndOutOfRange {
		printer.ErrPrintf("end_page (%d) greater than total pages (%d), less output than expected\n", endPagenum, readPagenum)
	} else if err != nil {
		printer.ErrPrintf("Unexpected error on page %d: %v\n", readPagenum, err)
		os.Exit(1)
	}

	// try to join subprocess
	if subproc != nil {
		printer.Stdout.(io.WriteCloser).Close()
		// join subprocess
		<-joinSubProc
	}
}
