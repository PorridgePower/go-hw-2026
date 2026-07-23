package main

import (
	"errors"
	"os"
)

func validate_source() error {
	srcInfo, err := os.Stat(from)
	if errors.Is(err, os.ErrNotExist) {
		return ErrFileNotExist
	}
	// TODO
	fsize := srcInfo.Size()
	if fsize <= int64(offset) {
		return ErrOffsetExceedsFileSize
	}
	return nil
}
