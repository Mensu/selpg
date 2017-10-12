package main

import (
	"bufio"
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
	var cmd *exec.Cmd

	// set output to where
	if len(dest) > 0 {
		cmd = exec.Command("lp", "-d", dest)
		// cmd = exec.Command("cat", "-n")
		printer.Stdout, _ = cmd.StdinPipe()
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr
		// when subprocess is done, join it
		go func() {
			cmd.Run()
			joinSubProc <- time.Now()
		}()
	}

	// set input from where
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
	if cmd != nil {
		// await the subprocess for at most 1 sec
		go func() {
			joinSubProc <- <-time.After(time.Second)
		}()
		// join subprocess
		<-joinSubProc
	}
}
