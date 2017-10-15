package main

import (
	"bufio"
	"io"
	"os"
	"os/exec"
)

var logger *Logger
var err error

func main() {
	// construct global printer
	logger = &Logger{
		os.Args[0],
		os.Stdout,
		os.Stderr,
	}

	// process flags
	err = processFlags()
	if err != nil {
		logger.ErrPrintf("%v\n%s", err, flagset.FlagUsagesWrapped(0))
		os.Exit(2)
	}

	// parsed results
	startPagenum, _ := flagset.GetUint(flagInfo["s"])
	endPagenum, _ := flagset.GetUint(flagInfo["e"])
	pageSize, _ := flagset.GetUint(flagInfo["l"])
	useF, _ := flagset.GetBool(flagInfo["f"])
	dest, _ := flagset.GetString(flagInfo["d"])

	var subproc *exec.Cmd

	// set where to output to
	if len(dest) > 0 {
		subproc = exec.Command("lp", "-d", dest)
		// subproc = exec.Command("cat", "-n")
		logger.Stdout, err = subproc.StdinPipe()
		if err != nil {
			logger.ErrPrintf("Failed to send data to printer %s: %v\n", dest, err)
			os.Exit(1)
		}
		subproc.Stdout = os.Stdout
		subproc.Stderr = os.Stderr
		subproc.Start()
	}

	// set where to input from
	var input = os.Stdin
	var readPagenum uint
	if len(filename) > 0 {
		input, err = os.Open(filename)
		defer input.Close()
		if err != nil {
			logger.ErrPrintln(err)
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
		logger.ErrPrintf("start_page (%d) greater than total pages (%d), no output written\n", startPagenum, readPagenum)
	} else if err == errEndOutOfRange {
		logger.ErrPrintf("end_page (%d) greater than total pages (%d), less output than expected\n", endPagenum, readPagenum)
	} else if err != nil {
		logger.ErrPrintf("Unexpected error on page %d: %v\n", readPagenum, err)
		os.Exit(1)
	}

	if subproc != nil {
		logger.Stdout.(io.WriteCloser).Close()
		// await subprocess
		subproc.Wait()
	}
}
