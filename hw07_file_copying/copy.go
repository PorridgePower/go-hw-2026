package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if err := os.MkdirAll(filepath.Dir(toPath), 0o755); err != nil {
		return err
	}
	fileDst, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer fileDst.Close()

	fileSrc, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fileSrc.Close()

	finfo, _ := fileSrc.Stat()
	if limit == 0 {
		limit = finfo.Size()
	}
	readLimit := min(limit, finfo.Size()-offset)
	section := io.NewSectionReader(fileSrc, offset, readLimit)

	buf := make([]byte, 1024)
	var totalBytes int
	for {
		readBytes, err := section.Read(buf)

		if readBytes > 0 {
			writeBytes, errW := fileDst.Write(buf[:readBytes])
			if errW != nil {
				return errW
			}
			totalBytes += writeBytes
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}
