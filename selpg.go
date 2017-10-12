package main

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

var errStartOutOfRange = errors.New("Start Out of Range")
var errEndOutOfRange = errors.New("End Out of Range")

// selpgByLine reads page ``startPagenum`` to ``endPagenum`` from ```reader``
func selpgByLine(reader *bufio.Reader, startPagenum uint, endPagenum uint, pageSize uint) (uint, error) {
	startLinenum := (startPagenum-1)*pageSize + 1
	linenum := uint(1)
	for pagenum := uint(1); ; linenum++ {
		// corrects pagenum
		if linenum != 1 && (pageSize == 1 || linenum%pageSize == 1) {
			pagenum++
		}
		if pagenum > endPagenum {
			break
		}
		line, err := reader.ReadString('\n')
		// not reached startLine
		if linenum < startLinenum {
			if err == io.EOF {
				return pagenum, errStartOutOfRange
			} else if err != nil {
				return pagenum, err
			}
			continue
		}
		line = strings.TrimSuffix(line, "\r\n")
		line = strings.TrimSuffix(line, "\n")
		if err == nil {
			printer.Println(line)
			continue
		} else if err == io.EOF {
			printer.Print(line)
			if pagenum < endPagenum {
				return pagenum, errEndOutOfRange
			}
			return pagenum, nil
		} else {
			return pagenum, err
		}
	}
	return endPagenum, nil
}

// selpgByF reads page ``startPagenum`` to ``endPagenum`` from ```reader``
func selpgByF(reader *bufio.Reader, startPagenum uint, endPagenum uint) (uint, error) {
	var pagenum uint
	for pagenum = uint(1); pagenum <= endPagenum; pagenum++ {
		line, err := reader.ReadString('\f')
		if pagenum < startPagenum {
			if err == io.EOF {
				return pagenum, errStartOutOfRange
			} else if err != nil {
				return pagenum, err
			}
			continue
		}
		if err == nil {
			printer.Print(line)
			continue
		} else if err == io.EOF {
			printer.Print(line)
			if pagenum < endPagenum {
				return pagenum, errEndOutOfRange
			}
			return pagenum, nil
		} else {
			return pagenum, err
		}
	}
	return endPagenum, nil
}
