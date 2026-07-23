package main

import (
	"errors"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// os.MkdirAll(toPath)
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
	fileSrc.Seek(offset, 0)

	finfo, _ := fileSrc.Stat()
	maxBorder := min(offset+limit, finfo.Size())
	section := io.NewSectionReader(fileSrc, offset, maxBorder-offset)

	buf := make([]byte, 1024)
	var totalBytes int
	for {
		readedBytes, err := section.Read(buf)

		if readedBytes > 0 {
			// io.LimitReader()
			writeBytes, errW := fileDst.Write(buf[:readedBytes])
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
