package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

// Выключил линтер превысил всего на 1 :)
//nolint:funlen
func Copy(fromPath string, toPath string, offset, limit int64) error {
	fileFrom, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("can't open file `%s`: %w", fromPath, err)
	}
	defer fileFrom.Close()

	fileFromStat, err := fileFrom.Stat()
	if err != nil {
		return fmt.Errorf("can't get file stat `%s`: %w", fromPath, err)
	}

	if fileFromStat.Size() == 0 {
		return ErrUnsupportedFile
	}

	if offset > fileFromStat.Size() {
		return ErrOffsetExceedsFileSize
	}

	_, err = fileFrom.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("can't seek: %w", err)
	}

	fileTo, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("can't create file `%s`: %w", toPath, err)
	}
	defer fileTo.Close()

	finBytesCount := limit
	if limit == 0 {
		finBytesCount = fileFromStat.Size()
	}
	var chunkSize int64 = 1024
	if chunkSize > finBytesCount {
		chunkSize = finBytesCount
	}
	var remained = finBytesCount

	progress := pb.Start64(finBytesCount)
	progress.Set(pb.Bytes, true)

	for {
		if chunkSize > remained {
			chunkSize = remained
		}
		if remained == 0 {
			break
		}
		written, err := io.CopyN(fileTo, fileFrom, chunkSize)
		progress.Add64(written)
		remained -= written
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("can't copy file %w", err)
		}
	}
	progress.Finish()

	return nil
}
