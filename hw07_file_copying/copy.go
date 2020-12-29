package main

import (
	"errors"
	"io"
	"log"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	//fromPath validation
	fi, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	if !fi.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	//offset validation
	if offset > fi.Size() {
		return ErrOffsetExceedsFileSize
	}

	limit = adjustLimit(limit, fi.Size(), offset)

	var srcFile *os.File

	srcFile, err = os.Open(fromPath)
	if err != nil {
		return err
	}
	defer func() {
		if err = srcFile.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	dstFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer func() {
		if err = dstFile.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	_, err = srcFile.Seek(offset, 0)
	if err != nil {
		return err
	}

	bar := pb.Full.Start64(limit)
	barReader := bar.NewProxyReader(srcFile)

	_, err = io.CopyN(dstFile, barReader, limit)
	if err != nil && err != io.EOF {
		return err
	}

	bar.Finish()

	return nil
}

func adjustLimit(limit int64, fileSize int64, offset int64) int64 {
	if limit == 0 || limit > fileSize-offset {
		limit = fileSize - offset
	}

	return limit
}
