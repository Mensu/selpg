package main

import (
	"fmt"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
)

var (
	flagset  = flag.NewFlagSet("selpg", flag.ContinueOnError)
	flagInfo = map[string]string{
		"s": "startpage",
		"e": "endpage",
		"l": "pagesize",
		"f": "f",
		"d": "destination",
	}
	filename = ""
)

func makeHasFlagMap() map[string]bool {
	hasFlag := make(map[string]bool)
	isFilename := false
	for _, oneArg := range os.Args {
		toContinue := false
		// check whether oneArg is a flag
		for short, full := range flagInfo {
			if strings.HasPrefix(oneArg, "--"+full) {
				hasFlag[short] = true
				// next oneArg can be a filename
				isFilename = true
				toContinue = true
				break
			} else if strings.HasPrefix(oneArg, "-"+short) {
				hasFlag[short] = true
				// next oneArg can be a filename if oneArg != "-s"
				isFilename = oneArg != short
				toContinue = true
				break
			}
		}
		if toContinue {
			continue
		}
		// oneArg is not a flag: it can be a filename
		if isFilename {
			filename = oneArg
		} // else: oneArg is like 10 in "-s 10"
		isFilename = true
	}
	return hasFlag
}

func processFlags() error {
	flagset.SortFlags = false

	// define flags
	flagset.UintP(flagInfo["s"], "s", 0, "start_page, must be specified")
	flagset.UintP(flagInfo["e"], "e", 0, "end_page, must be specified and no smaller than start_page")
	flagset.UintP(flagInfo["l"], "l", 72, "page_size, exclusive to -f")
	flagset.BoolP(flagInfo["f"], "f", false, "use \\f instead of page_size to delimit a page, exclusive to -l")
	flagset.StringP(flagInfo["d"], "d", "", "printer to output to (default stdout)")

	// parse flags
	err := flagset.Parse(os.Args[1:])
	if err != nil {
		os.Exit(2)
	}

	// parsed results
	startPagenum, _ := flagset.GetUint(flagInfo["s"])
	endPagenum, _ := flagset.GetUint(flagInfo["e"])
	pageSize, _ := flagset.GetUint(flagInfo["l"])
	dest, _ := flagset.GetString(flagInfo["d"])

	hasFlag := makeHasFlagMap()
	_, hasS := hasFlag["s"]
	_, hasE := hasFlag["e"]
	_, hasL := hasFlag["l"]
	_, hasF := hasFlag["f"]
	_, hasD := hasFlag["d"]

	if !hasS {
		return fmt.Errorf("missing argument for -s")
	}
	if startPagenum == 0 {
		return fmt.Errorf("invalid argument \"%v\" for -s: expected a positive integer", startPagenum)
	}
	if !hasE {
		return fmt.Errorf("missing argument for -e")
	}
	if endPagenum == 0 {
		return fmt.Errorf("invalid argument \"%v\" for -e: expected a positive integer", endPagenum)
	}
	if startPagenum > endPagenum {
		return fmt.Errorf("invalid argument \"%v\" for -e: end_page (%d) must be no smaller than start_page (%d)", endPagenum, endPagenum, startPagenum)
	}
	if hasL && hasF {
		return fmt.Errorf("-l and -f must not be specified simultaneously")
	}
	if hasL && pageSize == 0 {
		return fmt.Errorf("invalid argument \"%v\" for -l: expected a positive integer", pageSize)
	}
	if hasD && len(dest) == 0 {
		return fmt.Errorf("invalid argument \"%v\" for -d: should name one destination", dest)
	}
	return nil
}
