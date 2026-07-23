package main

import (
	"errors"
	"flag"
	"log"
)

var (
	ErrFileNotExist = errors.New("non existing source file")
	// ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()
	if err := validate_source(); err != nil {
		log.Fatalf("Critical error: %s", err.Error())
	}
	err := Copy(from, to, offset, limit)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Place your code here.
}
